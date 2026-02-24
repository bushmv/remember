package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/bushmv/remember/docs"
	"github.com/bushmv/remember/src/middlewares"
	"github.com/bushmv/remember/src/services/auth"
	authRepo "github.com/bushmv/remember/src/services/auth/data/repository"
	"github.com/bushmv/remember/src/services/records"
	recordsRepo "github.com/bushmv/remember/src/services/records/data/repository"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer {token}
func main() {
	r := gin.Default()

	jwtSecretKey := os.Getenv("jwtSecretKey")
	jwtManager := auth.NewJWTManager([]byte(jwtSecretKey))
	authMiddleware := middlewares.NewAuthMiddleware(jwtManager)
	CreateAuthService(r, *jwtManager)
	CreateRecordsService(r, authMiddleware.RequireAuth())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func CreateAuthService(r *gin.Engine, jwtManager auth.JWTManager) {
	authRepository := authRepo.NewInMemoryRepository()
	authInteractor := auth.NewAuthInteractor(authRepository, jwtManager)
	userService := auth.NewAuthService(authInteractor)
	userService.BindRoutes(r)
}

func CreateRecordsService(r *gin.Engine, authMiddleware gin.HandlerFunc) {
	recordsRepository := recordsRepo.NewInMemoryRepository()
	recordsInteractor := records.NewRecordsInteractor(recordsRepository)
	recordsService := records.NewRecordsService(recordsInteractor)
	recordsService.BindRoutes(r, authMiddleware)
}

// HealthCheck godoc
// @Summary      Проверка здоровья сервиса
// @Description  Возвращает статус сервиса
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
