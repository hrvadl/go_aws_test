package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrvadl/go_aws_test/pkg/dto"
	"github.com/hrvadl/go_aws_test/pkg/services"
)

type AuthHandler interface {
	HandleRegister(*gin.Context)
	HandleLogin(ctx *gin.Context)
	HandleConfirm(ctx *gin.Context)
}

type Auth struct {
	authSrv services.Auth
}

func NewAuthHandler(authSrv services.Auth) AuthHandler {
	return &Auth{authSrv: authSrv}
}

func (s *Auth) HandleRegister(ctx *gin.Context) {
	var body dto.SignInDTO

	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.authSrv.SignUp(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true})
}

func (s *Auth) HandleLogin(ctx *gin.Context) {
	var body dto.LoginDTO

	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := s.authSrv.Login(&body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": res})
}

func (s *Auth) HandleConfirm(ctx *gin.Context) {
	var body dto.ConfirmDTO

	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.authSrv.Confirm(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true})
}
