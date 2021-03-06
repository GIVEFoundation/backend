package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
  "github.com/gin-contrib/static"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
  "github.com/rs/xid"
	"net/http"
	"os"
	"fmt"
	"time"
)

var log = logrus.New()

const (
	chainError        = 100 << iota
	dbConnectionError = 100 << iota
	kidsCreationError = 100 << iota
	uploadError       = 100 << iota
	mainGroup         = "api/v1/give"
  mediaURL          = "/media"
)

// APIError JSONAPI compatible error
type APIError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Files holds uploaded files urls
type UploadedFile struct {
  Filename      string          `json:"filename"`
  URL           string          `json:"url"`
}

// UploadedFiles set of uploaded files
type UploadedFiles struct {
  Files         []UploadedFile  `json:"files"`
}

// Kid data
type Kid struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	DateOfBirth   string         `json:"date_of_birth"`
	ParentsEmail  pq.StringArray `json:"parents_emails" gorm:"type:string[]"`
	StudentsPhoto string         `json:"students_photo"`
	SchoolName    string         `json:"school_name"`
	IDTagName     string         `json:"id_tag_name"`
}

// Kids is a Collection of kids in JSON API standard
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
    v1.POST("/upload", UploadFiles) // Upload photos, and other media files
	}
  // Serve uploaded media 
  r.Use(static.Serve(mediaURL, static.LocalFile("./upload", true)))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

// UploadFiles upload files interface
func UploadFiles(c *gin.Context) {
	outputDir := viper.Get("output_dir").(string)
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		  ae := APIError{Code: uploadError, Title: "Internal API error (upload) ", Detail: err.Error()}
		  c.JSON(http.StatusBadRequest, ae)
		  return
		return
	}
	files := form.File["files"]
	// Should return an array to be JSON API compatible
	var uploaded []UploadedFile
	for _, file := range files {
    guid := xid.New()
		outputFile := fmt.Sprintf("%s/%s_%s", outputDir,guid,file.Filename)
	  uploaded = append(uploaded, UploadedFile{Filename: file.Filename ,
                URL: fmt.Sprintf("%s/%s_%s",mediaURL,guid,file.Filename) })

		if err := c.SaveUploadedFile(file, outputFile); err != nil {
		  ae := APIError{Code: uploadError, Title: "Internal API error (upload) ", Detail: err.Error()}
		  c.JSON(http.StatusBadRequest, ae)
		  return
		}
	}

	c.JSON(http.StatusOK, UploadedFiles{Files: uploaded})
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

// UpdateKids updates data of a kid
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
