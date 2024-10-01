package controller

import (
	"net/http"

	"example.com/go-project/auth"
	"example.com/go-project/model"
	"example.com/go-project/services"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	usersService *services.UsersService
}

func NewUsersController(service *services.UsersService) *UsersController {
	return &UsersController{usersService: service}
}

func (controller *UsersController) RegisterUser(ctx *gin.Context) {
	var user model.Users
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the Create method
	err = controller.usersService.Create(user)
	if err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// User created successfully
	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"status": "ok",
		"data":   user,
		"msg":    "User added successfully.",
	})
}

func (controller *UsersController) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON data from request body to loginData struct
	err := ctx.ShouldBindJSON(&loginData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user
	user, err := controller.usersService.Authenticate(loginData.Email, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token including user role
	token, err := auth.GenerateJWT(user.Id, user.Email, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	// Add Bearer prefix to token
	bearerToken := "Bearer " + token

	// Set Authorization header in the response
	ctx.Header("Authorization", bearerToken)

	// Return token and user data
	ctx.JSON(http.StatusOK, gin.H{
		"token": bearerToken,
		"user":  user,
	})
}
