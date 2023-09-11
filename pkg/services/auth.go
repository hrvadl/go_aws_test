package services

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/hrvadl/go_aws_test/pkg/config"
	"github.com/hrvadl/go_aws_test/pkg/dto"
)

var UserPasswordAuthFlow = "USER_PASSWORD_AUTH"
var TokenAuthFlow = "REFRESH_TOKEN_AUTH"

type Auth interface {
	SignUp(in *dto.SignInDTO) error
	Login(in *dto.LoginDTO) (*cognitoidentityprovider.AuthenticationResultType, error)
	Confirm(in *dto.ConfirmDTO) error
	CheckIdentityMiddleware(*gin.Context)
}

type AuthService struct {
	jwt      JWTValidator
	config   *config.Env
	cognito  *cognitoidentityprovider.CognitoIdentityProvider
	clientID string
}

func NewAuthService(cognito *cognitoidentityprovider.CognitoIdentityProvider, cfg *config.Env, jwt JWTValidator) *AuthService {
	return &AuthService{
		clientID: cfg.CognitoID,
		cognito:  cognito,
		config:   cfg,
		jwt:      jwt,
	}
}

func (a *AuthService) SignUp(in *dto.SignInDTO) error {
	_, err := a.cognito.SignUp(&cognitoidentityprovider.SignUpInput{
		Username: &in.Username,
		Password: &in.Password,
		ClientId: &a.clientID,
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			in.GetCognitoEmailAttribute(),
		},
	})

	if err != nil {
		return err
	}

	return err
}

func (a *AuthService) Confirm(in *dto.ConfirmDTO) error {
	_, err := a.cognito.ConfirmSignUp(&cognitoidentityprovider.ConfirmSignUpInput{
		Username:         &in.Username,
		ClientId:         &a.clientID,
		ConfirmationCode: &in.Code,
	})

	return err
}

func (a *AuthService) Login(in *dto.LoginDTO) (*cognitoidentityprovider.AuthenticationResultType, error) {
	params := map[string]*string{
		"USERNAME": aws.String(in.Username),
		"PASSWORD": aws.String(in.Password),
	}

	res, err := a.cognito.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       &UserPasswordAuthFlow,
		ClientId:       &a.clientID,
		AuthParameters: params,
	})

	if err != nil {
		return nil, err
	}

	return res.AuthenticationResult, nil
}

func (a *AuthService) CheckIdentityMiddleware(ctx *gin.Context) {
	tokenH := ctx.GetHeader("Authorization")
	parts := strings.Split(tokenH, " ")

	if len(parts) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid auth method",
		})
		return
	}

	token := parts[1]

	if err := a.jwt.Validate(token); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Next()
}
