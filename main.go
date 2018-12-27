package main

import (
	"os"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var initialized = false
var ginLambda *ginadapter.GinLambda

// Handler to wrap gin to lambda
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		router := SetupRouter()
		ginLambda = ginadapter.New(router)
		initialized = true
	}

	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.Proxy(req)
}

func main() {
	switch runAs := os.Getenv("RUN_AS"); runAs {
	case "lambda":
		log.Info("Run as lambda")
		lambdaRun()
	default:
		log.Info("Run locally")
		localRun()
	}
}

func localRun() {
	router := SetupRouter()
	router.Run() // listen and serve on 0.0.0.0:8080
}

func lambdaRun() {
	lambda.Start(Handler)
}

// SetupRouter determine what to do for each api calls
func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}
