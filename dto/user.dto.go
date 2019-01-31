package dto

import (
	"github.com/thomaspoignant/user-microservice/entity"
)

// UserGetDto used when calling get API
type UserGetDto struct {
	// id of the User (format UUID)
	ID string `binding:"uuid,required" uri:"id" example:"8da8adc3-0ae9-47b2-884c-ee41e691ff57"`
}

// ConvertToUser convert UserGetDto to User
func (dto *UserGetDto) ConvertToUser() entity.User {
	return entity.User{
		ID: dto.ID,
	}
}

// UserDto describe a create/update request
type UserDto struct {
	FirstName string `json:"first_name" example:"John" binding:"required"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"name@example.com" binding:"emailvalidator,exists"`
}

// ConvertToUser convert UserCreateDto to User
func (dto *UserDto) ConvertToUser() entity.User {
	return entity.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
	}
}

// UserPatchDto describe a partial update request
type UserPatchDto struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"name@example.com"`
}

// ConvertToUser convert UserPatchDto to User
func (dto *UserPatchDto) ConvertToUser() entity.User {
	return entity.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
	}
}
