// +build integration

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/thomaspoignant/user-microservice/dto"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/thomaspoignant/user-microservice/entity"
	"github.com/thomaspoignant/user-microservice/testingUtils"
)

// UserAPISuite is the test suite for integration tests of API
type UserAPISuite struct {
	suite.Suite
}

var dynamoDbClient *dynamo.DB
var router *gin.Engine

func TestSuiteUserAPI(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(UserAPISuite))
}

func (suite *UserAPISuite) SetupSuite() {
	//load config files
	testingUtils.PrepareTest()

}

func (suite *UserAPISuite) SetupTest() {
	//setting up direct dynamoDB Client
	dynamoDbSession := session.Must(session.NewSession())
	awsConfig := aws.Config{
		Endpoint: aws.String(viper.GetString("DYNAMODB_ENDPOINT")),
		Region:   aws.String(viper.GetString("AWS_REGION")),
	}
	dynamoDbClient := dynamo.New(dynamoDbSession, &awsConfig)
	dynamodbTableName := viper.GetString("DYNAMODB_TABLE_NAME")
	log.Infof("Creating dynamodb table %s", dynamodbTableName)
	// delete the dynamoDbTable
	dynamoDbClient.Table(dynamodbTableName).DeleteTable()
	// create the dynamoDbTable again
	dynamoDbClient.CreateTable(dynamodbTableName, entity.User{}).Run()
	router = SetupRouter()
}

func (suite *UserAPISuite) TearDownTest() {
	// delete the dynamoDbTable
	dynamodbTableName := viper.GetString("DYNAMODB_TABLE_NAME")
	dynamoDbClient.Table(dynamodbTableName).DeleteTable()
}

// Test_user_get_invalid_uuid test error when id is not an uuid
func (suite *UserAPISuite) Test_user_get_invalid_uuid() {
	w := testingUtils.PerformHTTPRequest(router, "GET", "/v1/user/1234", nil)
	suite.Equal(http.StatusBadRequest, w.Code)

	expected := ApiErrorResponse{
		Error: "Key: 'UserGetDto.ID' Error:Field validation for 'ID' failed on the 'uuid' tag",
	}

	var response ApiErrorResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	suite.Nil(err)
	suite.Equal(expected, response)
}

// Test_user_get_valid_uuid_not_found test error when id is not an uuid
func (suite *UserAPISuite) Test_user_get_valid_uuid_not_found() {
	w := testingUtils.PerformHTTPRequest(router, "GET", "/v1/user/f831b13d-8227-47b3-a17e-f229d3b69335", nil)
	suite.Equal(http.StatusNotFound, w.Code)

	expected := ApiErrorResponse{
		Error: "No user found with id f831b13d-8227-47b3-a17e-f229d3b69335",
	}

	var response ApiErrorResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	suite.Nil(err)
	suite.Equal(expected, response)
}

func (suite *UserAPISuite) Test_insert_user() {
	userDto := dto.UserDto{
		Email:     "name@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	requestByte, _ := json.Marshal(userDto)
	requestReader := bytes.NewReader(requestByte)

	c := testingUtils.PerformHTTPRequest(router, "POST", "/v1/user/", requestReader)
	suite.Equal(http.StatusCreated, c.Code)

	var got entity.User
	err := json.Unmarshal([]byte(c.Body.String()), &got)
	suite.Nil(err)
	suite.NotNil(got.ID)
	suite.Equal(userDto.Email, got.Email)
	suite.Equal(userDto.FirstName, got.FirstName)
	suite.Equal(userDto.LastName, got.LastName)
}

func (suite *UserAPISuite) Test_insert_user_and_read_it() {
	userDto := dto.UserDto{
		Email:     "name@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	requestByte, _ := json.Marshal(userDto)
	requestReader := bytes.NewReader(requestByte)

	c := testingUtils.PerformHTTPRequest(router, "POST", "/v1/user/", requestReader)
	suite.Equal(http.StatusCreated, c.Code)

	var got entity.User
	err := json.Unmarshal([]byte(c.Body.String()), &got)
	suite.Nil(err)

	c = testingUtils.PerformHTTPRequest(router, "GET", "/v1/user/"+got.ID, nil)
	suite.Equal(http.StatusOK, c.Code)

	var read entity.User
	err = json.Unmarshal([]byte(c.Body.String()), &read)
	suite.Nil(err)

	suite.Equal(got, read)
}

func (suite *UserAPISuite) Test_insert_user_and_delete_it() {
	userDto := dto.UserDto{
		Email:     "name@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	requestByte, _ := json.Marshal(userDto)
	requestReader := bytes.NewReader(requestByte)

	c := testingUtils.PerformHTTPRequest(router, "POST", "/v1/user/", requestReader)
	suite.Equal(http.StatusCreated, c.Code)

	var got entity.User
	err := json.Unmarshal([]byte(c.Body.String()), &got)
	suite.Nil(err)

	c = testingUtils.PerformHTTPRequest(router, "DELETE", "/v1/user/"+got.ID, nil)
	suite.Equal(http.StatusNoContent, c.Code)

	c = testingUtils.PerformHTTPRequest(router, "GET", "/v1/user/"+got.ID, nil)
	suite.Equal(http.StatusNotFound, c.Code)
}

func (suite *UserAPISuite) Test_insert_with_put_and_read() {
	userDto := dto.UserDto{
		Email:     "name@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	requestByte, _ := json.Marshal(userDto)
	requestReader := bytes.NewReader(requestByte)

	id := uuid.NewV4().String()
	c := testingUtils.PerformHTTPRequest(router, "PUT", "/v1/user/"+id, requestReader)
	suite.Equal(http.StatusCreated, c.Code)

	var got entity.User
	err := json.Unmarshal([]byte(c.Body.String()), &got)
	suite.Nil(err)
	suite.Equal(id, got.ID)

	c = testingUtils.PerformHTTPRequest(router, "GET", "/v1/user/"+id, nil)
	suite.Equal(http.StatusOK, c.Code)
	var read entity.User
	err = json.Unmarshal([]byte(c.Body.String()), &read)
	suite.Nil(err)

	suite.Equal(got, read)
}

func (suite *UserAPISuite) Test_insert_user_and_update_it() {
	userDto := dto.UserDto{
		Email:     "name@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	requestByte, _ := json.Marshal(userDto)
	requestReader := bytes.NewReader(requestByte)

	c := testingUtils.PerformHTTPRequest(router, "POST", "/v1/user/", requestReader)
	suite.Equal(http.StatusCreated, c.Code)

	var inserted entity.User
	err := json.Unmarshal([]byte(c.Body.String()), &inserted)
	suite.Nil(err)

	updatedUser := dto.UserDto{
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	requestByte, _ = json.Marshal(updatedUser)
	requestReader = bytes.NewReader(requestByte)

	c = testingUtils.PerformHTTPRequest(router, "PUT", "/v1/user/"+inserted.ID, requestReader)
	suite.Equal(http.StatusOK, c.Code)

	var updated entity.User
	err = json.Unmarshal([]byte(c.Body.String()), &updated)
	suite.Nil(err)

	suite.NotEqual(inserted, updated)
	suite.Equal("john.doe@example.com", updated.Email)
	suite.Equal("John", updated.FirstName)
	suite.Equal("Doe", updated.LastName)
	suite.Equal(inserted.ID, updated.ID)
}

func (suite *UserAPISuite) Test_delete_invalid_id() {
	c := testingUtils.PerformHTTPRequest(router, "DELETE", "/v1/user/f831b13d-8227-b3-a17e-f229d3b69335", nil)
	suite.Equal(http.StatusBadRequest, c.Code)
}

func (suite *UserAPISuite) Test_delete_not_found_id() {
	c := testingUtils.PerformHTTPRequest(router, "DELETE", "/v1/user/f831b13d-8227-47b3-a17e-f229d3b69335", nil)
	suite.Equal(http.StatusNotFound, c.Code)
}

func (suite *UserAPISuite) Test_insert_user_and_patch_it() {
	userDto := dto.UserDto{
		Email:     "name@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	requestByte, _ := json.Marshal(userDto)
	requestReader := bytes.NewReader(requestByte)

	c := testingUtils.PerformHTTPRequest(router, "POST", "/v1/user/", requestReader)
	suite.Equal(http.StatusCreated, c.Code)

	var inserted entity.User
	err := json.Unmarshal([]byte(c.Body.String()), &inserted)
	suite.Nil(err)

	updatedUser := dto.UserDto{
		Email: "john.doe@example.com",
	}
	requestByte, _ = json.Marshal(updatedUser)
	requestReader = bytes.NewReader(requestByte)

	c = testingUtils.PerformHTTPRequest(router, "PATCH", "/v1/user/"+inserted.ID, requestReader)
	suite.Equal(http.StatusOK, c.Code)

	var updated entity.User
	err = json.Unmarshal([]byte(c.Body.String()), &updated)
	suite.Nil(err)

	suite.NotEqual(inserted, updated)
	suite.Equal("john.doe@example.com", updated.Email)
	suite.Equal("John", updated.FirstName)
	suite.Equal("Doe", updated.LastName)
	suite.Equal(inserted.ID, updated.ID)
}

func (suite *UserAPISuite) Test_patch_not_found_id() {
	updatedUser := dto.UserDto{
		Email: "john.doe@example.com",
	}
	requestByte, _ := json.Marshal(updatedUser)
	requestReader := bytes.NewReader(requestByte)

	c := testingUtils.PerformHTTPRequest(router, "PATCH", "/v1/user/f831b13d-8227-47b3-a17e-f229d3b69335", requestReader)
	suite.Equal(http.StatusNotFound, c.Code)
}

func (suite *UserAPISuite) Test_patch_invalid_id() {
	updatedUser := dto.UserDto{
		Email: "john.doe@example.com",
	}
	requestByte, _ := json.Marshal(updatedUser)
	requestReader := bytes.NewReader(requestByte)

	c := testingUtils.PerformHTTPRequest(router, "PATCH", "/v1/user/f831b13d-8227-b3-a17e-f229d3b69335", requestReader)
	suite.Equal(http.StatusBadRequest, c.Code)
}
