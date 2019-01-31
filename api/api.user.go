package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/guregu/dynamo"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thomaspoignant/user-microservice/dto"
	"github.com/thomaspoignant/user-microservice/entity"
)

// UserController definition of a controller for user
type userController struct {
	userService *entity.UserService
}

// NewUserController is the constructor of the object
func NewUserController() *userController {
	userController := new(userController)
	service, err := entity.NewUserService(viper.GetString("DYNAMODB_TABLE_NAME"))
	if err != nil {
		panic("Impossible to init User service : " + err.Error())
	}
	userController.userService = service
	return userController
}

// GetByID return a user by id
// @Summary return user with given id
// @Description return user with given id
// @Tags user
// @Consume json
// @Produce  json
// @Param id path string true "ID of the user"
// @Success 200 {object} entity.User
// @Failure 400 {object} api.ApiErrorResponse "Parameters error"
// @Failure 404 {object} api.ApiErrorResponse "No user found with this id"
// @Failure 500 {object} api.ApiErrorResponse "System error"
// @Router /v1/user/{id} [get]
func (controller userController) GetByID(c *gin.Context) {
	var dto dto.UserGetDto
	if err := c.ShouldBindUri(&dto); err != nil {
		log.Warnf("Invalid input (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusBadRequest, ApiErrorResponse{Error: err.Error()})
		return
	}
	user := dto.ConvertToUser()
	if err := controller.userService.FindByID(&user); err != nil {
		if err == dynamo.ErrNotFound {
			log.Warnf("No user found (%s) : %s", dto.ID, err.Error())
			c.JSON(http.StatusNotFound, ApiErrorResponse{Error: "No user found with id " + dto.ID})
			return
		}
		log.Errorf("Impossible to find user with id %s (%s) : %s", dto.ID, c.Request.RequestURI, err.Error())
		c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Create a new User
// @Summary Create a new user
// @Description create and save user in database
// @Tags user
// @Param user body dto.UserDto false "user object in json"
// @Consume json
// @Produce  json
// @Success 201 {object} entity.User
// @Failure 400 {object} api.ApiErrorResponse "Parameters error"
// @Failure 500 {object} api.ApiErrorResponse "System error"
// @Router /v1/user/ [post]
func (controller userController) Create(c *gin.Context) {
	var dto dto.UserDto
	if err := c.ShouldBindBodyWith(&dto, binding.JSON); err != nil {
		log.Warnf("Invalid input (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusBadRequest, ApiErrorResponse{Error: err.Error()})
		return
	}

	userToInsert := dto.ConvertToUser()
	userToInsert.ID = "" // to be sure to create new element
	if err := controller.userService.Save(&userToInsert); err != nil {
		log.Errorf("Impossible to create user (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, userToInsert)
}

// CompleteUpdate a User
// @Summary Full update for a user object
// @Description update and save user in database
// @Tags user
// @Param id path string true "ID of the user"
// @Param user body dto.UserDto false "user object in json"
// @Consume json
// @Produce  json
// @Success 200 {object} entity.User "User is updated"
// @Success 201 {object} entity.User "User is created"
// @Failure 400 {object} api.ApiErrorResponse "Parameters error"
// @Failure 404 {object} api.ApiErrorResponse "No user found with this id"
// @Failure 500 {object} api.ApiErrorResponse "System error"
// @Router /v1/user/{id} [put]
func (controller userController) CompleteUpdate(c *gin.Context) {
	// read id in url path
	var idDto dto.UserGetDto
	if err := c.ShouldBindUri(&idDto); err != nil {
		log.Warnf("Invalid input (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusBadRequest, ApiErrorResponse{Error: err.Error()})
		return
	}
	user := idDto.ConvertToUser()
	user.ID = idDto.ID
	createUser := false
	if err := controller.userService.FindByID(&user); err != nil {
		if err != dynamo.ErrNotFound {
			log.Errorf("Impossible to find user with id %s (%s) : %s", idDto.ID, c.Request.RequestURI, err.Error())
			c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
			return
		}
		createUser = true
	}

	var dto dto.UserDto
	if err := c.ShouldBindBodyWith(&dto, binding.JSON); err != nil {
		log.Warnf("Invalid input (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusBadRequest, ApiErrorResponse{Error: err.Error()})
		return
	}

	userToUpdate := dto.ConvertToUser()
	userToUpdate.ID = idDto.ID
	userToUpdate.CreatedAt = user.CreatedAt
	if err := controller.userService.Save(&userToUpdate); err != nil {
		log.Errorf("Impossible to create/update user with id %s (%s) : %s", idDto.ID, c.Request.RequestURI, err.Error())
		c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
		return
	}

	var httpStatus int
	if createUser {
		httpStatus = http.StatusCreated
	} else {
		httpStatus = http.StatusOK
	}
	c.JSON(httpStatus, userToUpdate)
}

// PartialUpdate update only some fields
// @Summary Update only some fields
// @Description Update only some fields
// @Tags user
// @Param id path string true "ID of the user"
// @Param user body dto.UserPatchDto false "user with only modified fields in json"
// @Consume json
// @Produce  json
// @Success 200 {object} entity.User "User is updated"
// @Failure 400 {object} api.ApiErrorResponse "Parameters error"
// @Failure 404 {object} api.ApiErrorResponse "No user found with this id"
// @Failure 500 {object} api.ApiErrorResponse "System error"
// @Router /v1/user/{id} [patch]
func (controller userController) PartialUpdate(c *gin.Context) {
	var idDto dto.UserGetDto
	if err := c.ShouldBindUri(&idDto); err != nil {
		log.Warnf("Invalid input (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusBadRequest, ApiErrorResponse{Error: err.Error()})
		return
	}
	user := idDto.ConvertToUser()
	user.ID = idDto.ID
	if err := controller.userService.FindByID(&user); err != nil {
		if err == dynamo.ErrNotFound {
			log.Warnf("No user found (%s) : %s", idDto.ID, err.Error())
			c.JSON(http.StatusNotFound, ApiErrorResponse{Error: "No user found with id " + idDto.ID})
			return
		}
		log.Errorf("Impossible to find user with id %s (%s) : %s", idDto.ID, c.Request.RequestURI, err.Error())
		c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
		return
	}

	var dto dto.UserPatchDto
	if err := c.ShouldBindBodyWith(&dto, binding.JSON); err != nil {
		log.Warnf("Invalid input (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusBadRequest, ApiErrorResponse{Error: err.Error()})
		return
	}
	userToUpdate := dto.ConvertToUser()

	// merging old user with modified fields
	if err := mergo.Merge(&userToUpdate, user); err != nil {
		log.Errorf("Impossible to merge old data with new one %s / %s : %s", user, userToUpdate, err.Error())
		return
	}

	if err := controller.userService.Save(&userToUpdate); err != nil {
		log.Errorf("Impossible to create/update user with id %s (%s) : %s", idDto.ID, c.Request.RequestURI, err.Error())
		c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, userToUpdate)
}

// Delete delete a user by id
// @Summary delete user with given id
// @Description delete user with given id
// @Tags user
// @Consume json
// @Produce  json
// @Param id path string true "ID of the user"
// @Success 204 "User is correctly deleted"
// @Failure 400 {object} api.ApiErrorResponse "Parameters error"
// @Failure 404 {object} api.ApiErrorResponse "No user found with this id"
// @Failure 500 {object} api.ApiErrorResponse "System error"
// @Router /v1/user/{id} [delete]
func (controller userController) Delete(c *gin.Context) {
	var dto dto.UserGetDto
	if err := c.ShouldBindUri(&dto); err != nil {
		log.Warnf("Invalid input (%s) : %s", c.Request.RequestURI, err.Error())
		c.JSON(http.StatusBadRequest, ApiErrorResponse{Error: err.Error()})
		return
	}
	user := dto.ConvertToUser()
	if err := controller.userService.FindByID(&user); err != nil {
		if err == dynamo.ErrNotFound {
			log.Warnf("No user found (%s) : %s", dto.ID, err.Error())
			c.JSON(http.StatusNotFound, ApiErrorResponse{Error: "No user found with id " + dto.ID})
			return
		}
		log.Errorf("Impossible to find user with id %s (%s) : %s", dto.ID, c.Request.RequestURI, err.Error())
		c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
		return
	}

	if err := controller.userService.Delete(&user); err != nil {
		log.Errorf("Impossible to delete user with id %s (%s) : %s", dto.ID, c.Request.RequestURI, err.Error())
		c.JSON(http.StatusInternalServerError, ApiErrorResponse{Error: err.Error()})
	}

	c.JSON(http.StatusNoContent, nil) // return 204
}
