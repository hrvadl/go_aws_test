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

var (
	UserPasswordAuthFlow = "USER_PASSWORD_AUTH"
	TokenAuthFlow        = "REFRESH_TOKEN_AUTH"
	UserCtxKey           = "USER"
)

type Auth interface {
	SignUp(in *dto.SignUpDTO) error
	Login(in *dto.LoginDTO) (*string, error)
	Confirm(in *dto.ConfirmDTO) error
	CheckIdentityMiddleware(*gin.Context)
	GetUserByName(username string) (*dto.UserDTO, error)
	Logout(token string) error
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

func (a *AuthService) SignUp(in *dto.SignUpDTO) error {
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

func (a *AuthService) Login(in *dto.LoginDTO) (*string, error) {
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

	return res.AuthenticationResult.AccessToken, nil
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
	username, err := a.jwt.Validate(token)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Set(UserCtxKey, username)
	ctx.Next()
}

func (a *AuthService) GetUserByName(username string) (*dto.UserDTO, error) {
	res, err := a.cognito.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &a.config.UserPoolID,
		Username:   &username,
	})

	if err != nil {
		return nil, err
	}

	var userEmail string

	for _, attr := range res.UserAttributes {
		if *attr.Name == "email" {
			userEmail = *attr.Value
		}
	}

	return &dto.UserDTO{Email: userEmail, Username: *res.Username}, nil
}

func (a *AuthService) Logout(username string) error {
	_, err := a.cognito.AdminUserGlobalSignOut(&cognitoidentityprovider.AdminUserGlobalSignOutInput{
		UserPoolId: &a.config.UserPoolID,
		Username:   &username,
	})

	return err
}
