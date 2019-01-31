package entity

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	"github.com/spf13/viper"
	"github.com/thomaspoignant/user-microservice/testingUtils"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// UserEntitySuite is a test suite for the entity
type UserEntitySuite struct {
	suite.Suite
}

var uniqDbName string
var userServiceTest *UserService
var dynamoDbClient *dynamo.DB

func TestSuiteUserEntity(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(UserEntitySuite))
}

func (suite *UserEntitySuite) SetupSuite() {
	//load config files
	testingUtils.PrepareTest()

	//setting up direct dynamoDB Client
	dynamoDbSession := session.Must(session.NewSession())
	awsConfig := aws.Config{
		Endpoint: aws.String(viper.GetString("DYNAMODB_ENDPOINT")),
		Region:   aws.String(viper.GetString("AWS_REGION")),
	}
	dynamoDbClient = dynamo.New(dynamoDbSession, &awsConfig)
}

func (suite *UserEntitySuite) SetupTest() {
	// create a table especially for the test
	uniqDbName = xid.New().String()
	log.Infof("Database name for Test : %s", uniqDbName)
	dynamoDbClient.CreateTable(uniqDbName, User{}).Run()

	// create a service user
	var err error
	userServiceTest, err = NewUserService(uniqDbName)
	suite.Nil(err)
}

func (suite *UserEntitySuite) TearDownTest() {
	// delete the dynamoDbTable
	dynamoDbClient.Table(uniqDbName).DeleteTable()
	emptyUserService := UserService{}
	userServiceTest = &emptyUserService
	uniqDbName = ""
}

// TestEmptyTableName newUserService with an empty table name
func (suite *UserEntitySuite) TestEmptyTableName() {
	_, err := NewUserService("")
	suite.NotNil(err)
}

// TestInsertUser insert a user and test that object is updated
func (suite *UserEntitySuite) TestInsertUser() {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}

	err := userServiceTest.Save(&user)
	suite.Nil(err)
	suite.NotNil(user.ID)
	suite.NotNil(user.UpdatedAt)
	suite.NotNil(user.CreatedAt)

	value, err := dynamoDbClient.Table(uniqDbName).Get("id", user.ID).Count()
	suite.Nil(err)
	suite.Equal(int64(1), value)
}

//TestInsertAndReadUser insert a user and try to read what is in database
func (suite *UserEntitySuite) TestInsertAndReadUser() {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}

	err := userServiceTest.Save(&user)
	suite.Nil(err)

	value, err := dynamoDbClient.Table(uniqDbName).Get("id", user.ID).Count()
	suite.Nil(err)
	suite.Equal(int64(1), value)

	result := User{
		ID: user.ID,
	}

	err = userServiceTest.FindByID(&result)
	suite.Nil(err)

	expectedCreatedAt, _ := user.CreatedAt.MarshalJSON()
	gotCreatedAt, _ := result.CreatedAt.MarshalJSON()
	suite.Equal(expectedCreatedAt, gotCreatedAt)

	expectedUpdatedAt, _ := user.UpdatedAt.MarshalJSON()
	gotUpdatedAt, _ := result.UpdatedAt.MarshalJSON()
	suite.Equal(expectedUpdatedAt, gotUpdatedAt)

	suite.Equal(user.FirstName, result.FirstName)
	suite.Equal(user.LastName, result.LastName)
	suite.Equal(user.ID, result.ID)
}

//TestInsertAndDeleteUser insert a user and delete it
func (suite *UserEntitySuite) TestInsertAndDeleteUser() {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}
	err := userServiceTest.Save(&user)
	suite.Nil(err)

	numberOfUser, err := dynamoDbClient.Table(uniqDbName).Get("id", user.ID).Count()
	suite.Nil(err)
	suite.Equal(int64(1), numberOfUser)

	err = userServiceTest.Delete(&user)
	suite.Nil(err)

	numberOfUserAfterDelete, err := dynamoDbClient.Table(uniqDbName).Get("id", user.ID).Count()
	suite.Nil(err)
	suite.Equal(int64(0), numberOfUserAfterDelete)
}

// TestUpdateUser is creating and updating a user
func (suite *UserEntitySuite) TestUpdateUser() {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
	}

	err := userServiceTest.Save(&user)
	suite.Nil(err)

	value, err := dynamoDbClient.Table(uniqDbName).Get("id", user.ID).Count()
	suite.Nil(err)
	suite.Equal(int64(1), value)

	result := user
	err = userServiceTest.FindByID(&user)
	suite.Nil(err)
	expectedCreatedAt, _ := user.CreatedAt.MarshalJSON()
	gotCreatedAt, _ := result.CreatedAt.MarshalJSON()
	suite.Equal(expectedCreatedAt, gotCreatedAt)
	expectedUpdatedAt, _ := user.UpdatedAt.MarshalJSON()
	gotUpdatedAt, _ := result.UpdatedAt.MarshalJSON()
	suite.Equal(expectedUpdatedAt, gotUpdatedAt)
	suite.Equal(user.FirstName, result.FirstName)
	suite.Equal(user.LastName, result.LastName)
	suite.Equal(user.ID, result.ID)

	updatedUser := user
	updatedUser.FirstName = "John2"
	updatedUser.LastName = "Doe2"

	log.Infof("updatedUser : %s", updatedUser)
	err = userServiceTest.Save(&updatedUser)
	suite.Nil(err)
	gotUpdatedCreatedAt, _ := updatedUser.CreatedAt.MarshalJSON()
	suite.Equal(expectedCreatedAt, gotUpdatedCreatedAt)
	gotUpdatedUpdatedAt, _ := updatedUser.UpdatedAt.MarshalJSON()
	suite.NotEqual(expectedUpdatedAt, gotUpdatedUpdatedAt)
	suite.NotEqual(user.FirstName, updatedUser.FirstName)
	suite.NotEqual(user.LastName, updatedUser.LastName)
	suite.Equal(user.ID, result.ID)

	//we verify that we only have one entry
	value, err = dynamoDbClient.Table(uniqDbName).Get("id", updatedUser.ID).Count()
	suite.Nil(err)
	suite.Equal(int64(1), value)
}
