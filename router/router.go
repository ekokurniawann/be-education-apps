package router

import (
	"be-education/config"
	"be-education/handler"
	"be-education/middleware"
	"be-education/repository"
	"be-education/service"
	"be-education/utils"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitRouter(db *sqlx.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Tambahkan middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Ubah sesuai kebutuhan
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	jwtUtil := utils.NewJWTUtil(cfg.SecretKey)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtUtil)
	userHandler := handler.NewUserHandler(userService, cfg.Server.BaseURL)

	userChapterRepo := repository.NewUserChapterRepository(db)
	userChapterService := service.NewUserChapterService(userChapterRepo)
	userChapterHandler := handler.NewUserChapterHandler(userChapterService)

	authMiddleware := middleware.NewAuthMiddleware(cfg)

	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateMahasiswa)
			users.GET("/profile", authMiddleware.Auth(), userHandler.GetProfile)
			users.POST("/profile/image", authMiddleware.Auth(), userHandler.UpdateProfileImage)
			users.GET("/summary/students", authMiddleware.Auth(), authMiddleware.RequireRole("admin"), userHandler.GetStudentSummary)
			users.GET("/summary/admins", authMiddleware.Auth(), authMiddleware.RequireRole("admin"), userHandler.GetAdminSummary)
			users.POST("/admin", userHandler.CreateAdmin)
			users.GET("/mahasiswa", authMiddleware.Auth(), authMiddleware.RequireRole("admin"), userHandler.GetMahasiswaUsers)
			users.DELETE("/:id", authMiddleware.Auth(), authMiddleware.RequireRole("admin"), userHandler.DeleteUser)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}

		userChapters := api.Group("/user-chapters")
		{
			userChapters.Use(authMiddleware.Auth())
			userChapters.POST("", userChapterHandler.CreateUserChapter)
			userChapters.GET("", userChapterHandler.GetUserQuizScores)
			userChapters.POST("/check-completion", userChapterHandler.CheckUserChapterCompletion)
			userChapters.GET("/summary/all-scores", authMiddleware.RequireRole("admin"), userChapterHandler.GetAllUsersChapterScores)
		}
	}

	return r
}
