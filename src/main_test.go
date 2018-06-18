package main

import (
	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
  "github.com/spf13/viper"
	"net/http"
	"testing"
)

// Tests GIVE API 
func TestGIVEAPI(t *testing.T) {
	var err error

	// Read configuration
	err = readconfig()
	if err != nil {
		log.Error("configuration file error: %s\n", err)
		return
	}
  authToken := viper.Get("auth_token").(string)

	r := gofight.New()

  // mock data
  mock_new := "{\"id\":\"0x9AbAB02EcBe8A917C266681B37d1f45f56191bDb\", \"name\":\"Danny Wood\", \"date_of_birth\":\"2010-11-10\", \"parents_emails\":[\"pa@gmail.com\",\"ma@gmail.com\"], \"school_name\":\"1st School of Hawaii\", \"id_tag_name\":\"Big Danny\"}"
  t.Log(mock_new)
  mock_upd := "{\"id\":\"0x9AbAB02EcBe8A917C266681B37d1f45f56191bDb\", \"id_tag_name\":\"Medium Danny\"}"

	t.Log("Testing Kid POST ")
	r.POST("/api/v1/give/kids").
		SetDebug(true).
    SetHeader(gofight.H{"GIVEAPIToken": authToken}).
    SetBody(string(mock_new)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, res.Code)
		})

	t.Log("Testing Kid PUT ")
	r.PUT("/api/v1/give/kids").
		SetDebug(true).
    SetHeader(gofight.H{"GIVEAPIToken": authToken}).
    SetBody(string(mock_upd)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, res.Code)
		})
}
