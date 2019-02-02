package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	uuid "github.com/satori/go.uuid"
)

func Test_ConvertUserGetDtoToUser(t *testing.T) {
	randomID := uuid.NewV4().String()
	dto := UserGetDto{
		ID: randomID,
	}
	user := dto.ConvertToUser()
	assert.NotNil(t, user.ID)
	assert.Equal(t, randomID, user.ID)
}

func Test_ConvertUserDtoToUser(t *testing.T) {
	dto := UserDto{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "name@example.com",
	}
	user := dto.ConvertToUser()

	assert.Empty(t, user.ID)
	assert.Equal(t, dto.FirstName, user.FirstName)
	assert.Equal(t, dto.LastName, user.LastName)
	assert.Equal(t, dto.Email, user.Email)
}

func Test_ConvertUserPatchDtoToUser(t *testing.T) {
	dto := UserPatchDto{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "name@example.com",
	}
	user := dto.ConvertToUser()

	assert.Empty(t, user.ID)
	assert.Equal(t, dto.FirstName, user.FirstName)
	assert.Equal(t, dto.LastName, user.LastName)
	assert.Equal(t, dto.Email, user.Email)
}
