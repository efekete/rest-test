package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// api represents data about a record api.
type api struct {
	ResponseText string            `json:"Text,omitempty"`
	Version      string            `json:"Version,omitempty"`
	Timestamp    string            `json:"Timestamp,omitempty"`
	Environment  map[string]string `json:"Environment,omitempty"`
}

// api slice to seed record api data.

var (
	listenFlag      = flag.String("listen", "0.0.0.0:8080", "address and port to listen")
	textFlag        = flag.String("text", "", "text to put in response")
	versionFlag     = flag.String("version", "", "display version information")
	successCodeFlag = flag.Int("status-code", 200, "http response code, e.g.: 200")
	erroCodeFlag    = flag.Int("error-code", 503, "http response code for errors, e.g.: 503")
	errorRateFlag   = flag.Int("error-rate", 2, "error rate between 0 and 100, e.g.: 20 for 20%")

	includeEnvFlag       = flag.String("include-env-vars", "", "set the name of environment variables to include in response separated by comma")
	includeTimestampFlag = flag.Bool("include-timestamp", false, "set to true to include timestamp in response")

	readinessDelayFlag       = flag.Int("readiness-delay", 0, "delay in seconds before readiness endpoint returns 200")
	readinessDelayGitterFlag = flag.Int("readiness-gitter", 0, "max readiness gitter in seconds (random delay between 0 and this value)")
)

func main() {
	flag.Parse()

	router := gin.Default()
	router.GET("/", getRequest)
	router.GET("/ready", readiness)

	gitter := 0
	if *readinessDelayGitterFlag > 0 {
		gitter = rand.Intn(*readinessDelayGitterFlag)
	}
	sleepTime := *readinessDelayFlag + gitter
	fmt.Printf("Sleeping for %v seconds", sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)
	router.Run(*listenFlag)
}

func readiness(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// getapi responds with the list of all api as JSON.
func getRequest(c *gin.Context) {
	reponseCode := *successCodeFlag
	chance := rand.Intn(99)
	apiResponse := api{ResponseText: *textFlag}
	if chance < *errorRateFlag {
		reponseCode = *erroCodeFlag
	}

	if *versionFlag != "" {
		c.Header("X-App-Version", *versionFlag)
		apiResponse.Version = *versionFlag
	}

	if *includeEnvFlag != "" {
		envVars := strings.Split(*includeEnvFlag, ",")
		envMap := make(map[string]string)
		for i := range envVars {
			envMap[envVars[i]] = os.Getenv(envVars[i])
		}
		apiResponse.Environment = envMap //os.Getenv(*includeEnvFlag)
	}

	if *includeTimestampFlag {
		apiResponse.Timestamp = time.Now().Format(time.RFC3339)
	}

	c.JSON(reponseCode, apiResponse)
}
