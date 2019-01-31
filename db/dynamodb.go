package db

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var dynamoDbClient *dynamo.DB

// GetDynamodbTable return a dynamoDB table to manipulate data
func GetDynamodbTable(tableName string) (*dynamo.Table, error) {
	if tableName == "" {
		return nil, fmt.Errorf("you must supply a table name")
	}
	setupDynamoDBClient()
	table := dynamoDbClient.Table(tableName)
	return &table, nil

}

func setupDynamoDBClient() {
	if dynamoDbClient == nil {
		// we setup the dynamoDb connection
		dynamoDbSession := session.Must(session.NewSession())
		awsConfig := aws.Config{
			Endpoint: aws.String(viper.GetString("DYNAMODB_ENDPOINT")),
			Region:   aws.String(viper.GetString("AWS_REGION")),
		}
		dynamoDbClient = dynamo.New(dynamoDbSession, &awsConfig)
	}
}
