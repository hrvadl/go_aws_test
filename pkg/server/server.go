package server

import (
	"log"
	"net/http"

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
	jwt := services.NewJWTValidator(cfg)
	session := aws.NewSession(cfg)
	cognito := aws.NewCognito(session)
	auth := services.NewAuthService(cognito, cfg, jwt)

	authH := handlers.NewAuthHandler(auth)

	protected := s.srv.Group("", auth.CheckIdentityMiddleware)
	public := s.srv.Group("")
	static := s.srv.Group("/public")

	static.StaticFS("", http.Dir("./pkg/public"))

	public.POST("/login", authH.HandleLogin)
	public.POST("/sign-up", authH.HandleRegister)
	public.POST("/confirm", authH.HandleConfirm)

	protected.GET("/home", func(ctx *gin.Context) {})
	protected.GET("/log-out", authH.HandleLogout)
	protected.GET("/me", authH.HandleGetMe)

	return nil
}

func (s *Server) Run() error {
	return s.srv.Run()
}
