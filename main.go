package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// api represents data about a record api.
type api struct {
	ReqID     string `json:"req_ID"`
	RestID    string `json:"rest_ID"`
	Timestamp string `json:"time"`
}

// api slice to seed record api data.

var restID = 0

func main() {
	rand.Seed(time.Now().UnixNano())
	restID = rand.Intn(100)
	router := gin.Default()
	router.GET("/api", getRequest)

	router.Run("0.0.0.0:8080")
}

// getapi responds with the list of all api as JSON.
func getRequest(c *gin.Context) {
	var apiO = api{ReqID: fmt.Sprint(rand.Intn(1000000)), RestID: fmt.Sprint(restID), Timestamp: time.Now().Format(time.RFC3339)}
	c.IndentedJSON(http.StatusOK, apiO)
}
