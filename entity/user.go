package entity

import (
	"time"

	"github.com/guregu/dynamo"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/thomaspoignant/user-microservice/db"
)

// User is the entity in the DB
type User struct {
	// id of the User (format UUID)
	ID        string `dynamo:"id,hash" json:"id" example:"8da8adc3-0ae9-47b2-884c-ee41e691ff57"`
	FirstName string `dynamo:"first_name" json:"first_name" example:"John"`
	LastName  string `dynamo:"last_name" json:"last_name,omitempty" example:"Doe"`
	Email     string `dynamo:"email" json:"email" example:"name@example.com"`
	// creation date of the entry (format example 2019-01-17T21:03:08.373394+01:00), can't be update throught the API
	CreatedAt time.Time `dynamo:"created_at" json:"created_at" example:"2019-01-17T21:03:08.373394+01:00"`
	// last update date of the entry (format example 2019-01-17T21:03:08.373394+01:00)
	UpdatedAt time.Time `dynamo:"updated_at" json:"updated_at" example:"2019-01-17T21:03:08.373394+01:00"`
}

// UserService is the struct of the service
type UserService struct {
	UserTable *dynamo.Table
}

// NewUserService init the connection to dynamoDB database
func NewUserService(tableName string) (*UserService, error) {
	table, err := db.GetDynamodbTable(tableName)
	if err != nil {
		return nil, err
	}
	return &UserService{
		UserTable: table,
	}, nil
}

// Save a user in DB
func (service *UserService) Save(user *User) error {
	now := time.Now()
	if user.ID == "" {
		user.ID = uuid.NewV4().String()
	}
	if user.CreatedAt == *new(time.Time) {
		user.CreatedAt = now
	}
	user.UpdatedAt = now
	return service.UserTable.Put(user).Run()
}

// FindByID retrieve a user from DB
func (service *UserService) FindByID(user *User) error {
	log.WithField("user", user).Info("Retrieve user")
	return service.UserTable.Get("id", user.ID).One(user)
}

// Delete a user from DB
func (service *UserService) Delete(user *User) error {
	return service.UserTable.Delete("id", user.ID).Run()
}
