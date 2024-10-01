// package main

// import (
// 	"example.com/go-project/auth"
// 	"example.com/go-project/authrequired"
// 	"example.com/go-project/config"
// 	"example.com/go-project/controller"
// 	"example.com/go-project/helper"
// 	"example.com/go-project/model"
// 	"example.com/go-project/model/repository"
// 	"example.com/go-project/services"
// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator"
// )

// func main() {
// 	print("Server Started.\n\n\n")

// 	// Setup the database and validation
// 	db := config.DatabaseConnection()
// 	validate := validator.New()

// 	// AutoMigrate tables
// 	db.AutoMigrate(&model.Tags{}, &model.Neche{}, &model.Users{})

// 	// Tags setup
// 	tagsRepository := repository.NewTagsRepositoryImpl(db)
// 	tagsService := services.NewTagsServiceImpl(tagsRepository, validate)
// 	tagsController := controller.NewTagsController(tagsService)

// 	// Neches setup
// 	nechesRepository := repository.NewNecheRepositoryImpl(db)
// 	nechesService := services.NewNecheServiceImpl(nechesRepository, validate, tagsRepository)
// 	nechesController := controller.NewNecheController(nechesService)

// 	// User setup
// 	userRepo := repository.NewUsersRepository(db)
// 	userService := services.NewUsersService(userRepo)
// 	userController := controller.NewUsersController(userService)

// 	// Create the base router
// 	router := gin.Default()

// 	// Public routes (no authentication required)
// 	publicRouter := router.Group("/user")
// 	publicRouter.POST("/register", userController.RegisterUser)
// 	publicRouter.POST("/login", userController.Login)

// 	// Admin routes (requires Admin role)
// 	adminRouter := router.Group("/admin")
// 	adminRouter.Use(authrequired.RoleBasedAuth("Admin"))
// 	{
// 		adminRouter.DELETE("/neches/:necheId", nechesController.Delete)
// 		adminRouter.DELETE("/tags/:tagId", tagsController.Delete)
// 		// Uncomment if you have additional admin routes
// 		// adminRouter.GET("/sorted", tagsController.FindAllSorted)
// 	}

// 	// User routes (requires User or Admin role)
// 	userRouter := router.Group("/user")
// 	userRouter.Use(authrequired.RoleBasedAuth("User"))
// 	{
// 		userRouter.GET("/tags", tagsController.FindAll)
// 		userRouter.GET("/neches", nechesController.FindAll)
// 		userRouter.PATCH("/tags/:tagId", tagsController.Update)
// 		userRouter.POST("/tags", tagsController.Create)
// 		userRouter.GET("/tags/:tagId", tagsController.FindById)
// 	}

// 	// Initialize Google OAuth routes
// 	auth.NewAuth()

// 	// Start the server
// 	err := router.Run(":8888")
// 	helper.ErrorPanic(err)
// }

package main

import (
	"example.com/go-project/auth"
	"example.com/go-project/authrequired"
	"example.com/go-project/config"
	"example.com/go-project/controller"
	"example.com/go-project/helper"
	"example.com/go-project/model"
	"example.com/go-project/model/repository"
	"example.com/go-project/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func main() {
	print("Server Started.\n\n\n")

	// Setup the database and validation
	db := config.DatabaseConnection()
	validate := validator.New()

	// AutoMigrate tables
	db.AutoMigrate(&model.Tags{}, &model.Neche{}, &model.Users{})

	// Tags setup
	tagsRepository := repository.NewTagsRepositoryImpl(db)
	tagsService := services.NewTagsServiceImpl(tagsRepository, validate)
	tagsController := controller.NewTagsController(tagsService)

	// Neches setup
	nechesRepository := repository.NewNecheRepositoryImpl(db)
	nechesService := services.NewNecheServiceImpl(nechesRepository, validate, tagsRepository)
	nechesController := controller.NewNecheController(nechesService)

	// User setup
	userRepo := repository.NewUsersRepository(db)
	userService := services.NewUsersService(userRepo)
	userController := controller.NewUsersController(userService)

	// Create the base router
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Initialize Google OAuth

	// Public routes (no authentication required)
	publicRouter := router.Group("/user")
	publicRouter.POST("/register", userController.RegisterUser)
	publicRouter.POST("/login", userController.Login)

	// Admin routes (requires Admin role)
	adminRouter := router.Group("/admin")
	adminRouter.Use(authrequired.RoleBasedAuth("Admin"))
	{
		adminRouter.DELETE("/neches/:necheId", nechesController.Delete)
		adminRouter.DELETE("/tags/:tagId", tagsController.Delete)
	}

	// User routes (requires User or Admin role)
	userRouter := router.Group("/user")
	userRouter.Use(authrequired.RoleBasedAuth("User"))
	{
		userRouter.GET("/tags", tagsController.FindAll)
		userRouter.GET("/neches", nechesController.FindAll)
		userRouter.PATCH("/tags/:tagId", tagsController.Update)
		userRouter.POST("/tags", tagsController.Create)
		userRouter.GET("/tags/:tagId", tagsController.FindById)
	}

	auth.NewAuth(router, userService)
	// Start the server
	err := router.Run(":8888")
	helper.ErrorPanic(err)
}
