package server

import (
	"fmt"
	"log"
	"os"

	"github.com/irmadev7/tripmate-backend/internal/auth"
	"github.com/irmadev7/tripmate-backend/internal/itinerary"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/cache"
	"github.com/irmadev7/tripmate-backend/internal/pkg/config"
	"github.com/irmadev7/tripmate-backend/internal/repository"
	"github.com/irmadev7/tripmate-backend/internal/user"
	itineraryV1 "github.com/irmadev7/tripmate-backend/internal/v1/itinerary"
	userV1 "github.com/irmadev7/tripmate-backend/internal/v1/user"
	userV2 "github.com/irmadev7/tripmate-backend/internal/v2/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "github.com/irmadev7/tripmate-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	r  *gin.Engine
	db *gorm.DB
}

func New() (*Server, error) {
	_ = godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	db := config.ConnectDB()
	db.AutoMigrate(&model.User{})

	// register routes
	v1 := r.Group("/api/v1")
	v2 := r.Group("/api/v2")

	// repo
	placeRepo := repository.NewPlaceRepository(db)
	userRepo := repository.NewUserRepository(db)
	itineraryRepo := repository.NewItineraryRepository(db)

	// redis
	redisClient := cache.NewRedis()

	// service
	tokenService := auth.NewTokenService(secret)
	userService := user.NewService(&userRepo, tokenService)
	userServiceV2 := user.NewServiceV2(&userRepo, tokenService, redisClient)
	itineraryService := itinerary.NewService(&itineraryRepo, &userRepo, &placeRepo, redisClient)

	// routes
	userV1.RegisterRoutes(v1, userService, tokenService)
	userV2.RegisterRoutes(v2, userServiceV2, tokenService)
	itineraryV1.RegisterRoutes(v1, itineraryService, tokenService)

	return &Server{r: r, db: db}, nil
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running at http://localhost:" + port)
	return s.r.Run(":" + port)
}
