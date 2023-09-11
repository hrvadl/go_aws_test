package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hrvadl/go_aws_test/pkg/aws"
	"github.com/hrvadl/go_aws_test/pkg/config"
	"github.com/hrvadl/go_aws_test/pkg/handlers"
	"github.com/hrvadl/go_aws_test/pkg/services"
	"github.com/joho/godotenv"
)

type Server struct {
	srv *gin.Engine
}

func New() *Server {
	return &Server{
		srv: gin.Default(),
	}
}

func (s *Server) Setup() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("cannot load .env file: %v", err)
	}

	cfg := config.NewEnv()
	session := aws.NewSession()
	cognito := aws.NewCognito(session)
	auth := services.NewAuthService(cognito, cfg)

	authH := handlers.NewAuthHandler(auth)

	s.srv.POST("/login", authH.HandleLogin)
	s.srv.POST("/sign-up", authH.HandleRegister)
	s.srv.POST("/confirm", authH.HandleConfirm)

	return nil
}

func (s *Server) Run() error {
	return s.srv.Run()
}
