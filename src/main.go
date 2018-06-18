package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
  "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

var log = logrus.New()

const (
	chainError        = 100 << iota
	dbConnectionError = 100 << iota
	kidsCreationError = 100 << iota
	mainGroup         = "api/v1/give"
)

// APIError JSONAPI compatible error
type APIError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// APIError JSONAPI compatible error
type Kid struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	DateOfBirth   string   `json:"date_of_birth"`
	ParentsEmail  pq.StringArray `json:"parents_emails" gorm:"type:string[]"`
	StudentsPhoto []byte   `json:"students_photo"`
	SchoolName    string   `json:"school_name"`
	IDTagName     string   `json:"id_tag_name"`
}

type Kids struct {
	Kids []Kid `json:"kids"`
}

func main() {
	// Read configuration
	err := readconfig()
	if err != nil {
		log.Error("configuration file error: %s\n", err)
		return
	}

	// Setup logging
	logf, err := os.OpenFile("give_api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("log file error: %s\n", err)
		return
	}
	defer logf.Close()

	log.Formatter = new(logrus.JSONFormatter)
	log.Out = logf

	port := viper.Get("port").(string)
	restEngine().Run(port)
}

// Unauthorize generates Authentication error response
func Unauthorize(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
	c.Abort()
}

// TokenAuthMiddleware is a Middleware handler function
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := viper.Get("auth_token").(string)

		token := c.GetHeader("GIVEAPIToken")
		if token == "" {
			Unauthorize(c, 401, "GIVE Key Required")
			return
		}

		if token != authToken {
			Unauthorize(c, 401, "Invalid Key")
			return
		}
		c.Next()
	}
}

//restEngine returns a new gin engine
func restEngine() *gin.Engine {

	r := gin.Default()

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, GIVEAPIToken",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	v1 := r.Group(mainGroup)
	v1.Use(TokenAuthMiddleware())
	{
		v1.POST("/kids", CreateKids) // New Kid on the block
		v1.PUT("/kids", UpdateKids)  // Update Kid data
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

// CreateKids register a kid using wallet address as ID
func CreateKids(c *gin.Context) {
	dbURL := viper.Get("db_url").(string)

	// Opens database connection
	db, err := gorm.Open("postgres", dbURL)
	defer db.Close()

	if err != nil {
		ae := APIError{Code: dbConnectionError, Title: "Internal API error (db) ", Detail: err.Error()}
		c.JSON(http.StatusUnprocessableEntity, ae)
		return
	}

	var kid Kid
	c.BindJSON(&kid)

	if err := db.Create(&kid).Error; err != nil {
		ae := APIError{Code: kidsCreationError, Title: "Kid creation error", Detail: err.Error()}
		c.JSON(http.StatusUnprocessableEntity, ae)
		return
	}

	// Should return an array to be JSON API compatible
	var kids []Kid
	kids = append(kids, kid)

	// Return JSON result
	c.JSON(http.StatusOK, Kids{Kids: kids})
}

// UpdateTrade log an update transaction into blockchain platform
func UpdateKids(c *gin.Context) {
	var kid Kid
	c.BindJSON(&kid)

	// Should return an array to be JSON API compatible
	var kids []Kid
	kids = append(kids, kid)
	// Return JSON result
	c.JSON(http.StatusOK, Kids{Kids: kids})
}

// RegisterTrade log a transaction into blockchain platform
func RegisterTrade(c *gin.Context) {

	//	c.BindJSON(obj)
	//	log.Info(obj)

	id := time.Now().UnixNano()
	c.JSON(http.StatusOK, gin.H{"transactionID": id})
}

//readconfig read and parse configuration file
func readconfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}

	return nil
}
