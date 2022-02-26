package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"GoBlogClean/common"
	"GoBlogClean/domain"
	"GoBlogClean/pkg/auth"
)

type UserHandler struct {
	userUsecase auth.UserUsecase
}

func NewUserHandler(userUsecase auth.UserUsecase) UserHandler {
	return UserHandler{userUsecase: userUsecase}
}

func (uh *UserHandler) Signup(c *gin.Context) {
	var user *domain.User
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hashedBytePassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.Password = string(hashedBytePassword)

	u, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	uuidStr := u.String()
	user.ID = uuidStr

	if err := uh.userUsecase.CreateUser(user); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": " succeeded"})
}

func (uh *UserHandler) Login(c *gin.Context) {
	var user *domain.User
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userModel, err := uh.userUsecase.GetUserByUsername(user.Username)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(user.Password)); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("login success: user_id=%s", userModel.ID)

	jwtToken, err := common.CreateJWTToken(userModel)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	users, err := uh.userUsecase.GetUsers()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}
