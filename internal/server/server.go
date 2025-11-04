package server

import (
	"log"
	"os"

	"github.com/irmadev7/tripmate-backend/internal/itinerary"
	"github.com/irmadev7/tripmate-backend/internal/model"
	"github.com/irmadev7/tripmate-backend/internal/pkg/config"
	"github.com/irmadev7/tripmate-backend/internal/user"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	r  *gin.Engine
	db *gorm.DB
}

func New() *Server {
	_ = godotenv.Load()
	r := gin.Default()

	db := config.ConnectDB()
	db.AutoMigrate(&model.User{})

	// register routes
	user.RegisterRoutes(r, db)
	itinerary.RegisterRoutes(r, db)

	return &Server{r: r, db: db}
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running at http://localhost:" + port)
	return s.r.Run(":" + port)
}
