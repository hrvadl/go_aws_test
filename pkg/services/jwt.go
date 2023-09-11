package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hrvadl/go_aws_test/pkg/config"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type JWTValidator interface {
	Validate(string) error
}

type CognitoJWTValidator struct {
	cfg *config.Env
}

func NewJWTValidator(cfg *config.Env) JWTValidator {
	return &CognitoJWTValidator{
		cfg: cfg,
	}
}

func (v *CognitoJWTValidator) Validate(jwtToken string) error {
	pKey, err := v.getPublicKeys(v.cfg.AWSRegion, v.cfg.UserPoolID)

	if err != nil {
		log.Fatal("Error trying to get Cognito public keys, check your config")
	}

	keySet, _ := jwk.Parse(pKey)

	parsedToken, err := jwt.Parse([]byte(jwtToken), jwt.WithKeySet(keySet))

	if parsedToken == nil {
		return errors.New("invalid token")
	}

	tokenUse, _ := parsedToken.Get("token_use")

	if err != nil {
		return errors.New("invalid token")
	}

	if parsedToken.Issuer() != fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", v.cfg.AWSRegion, v.cfg.UserPoolID) {
		return errors.New("token is from a different pool_id")
	}

	if tokenUse != "id" && tokenUse != "access" {
		return errors.New("token is from a different source")
	}

	if time.Now().After(parsedToken.Expiration()) {
		return errors.New("token expired")
	}

	return nil
}

func (v *CognitoJWTValidator) getPublicKeys(region string, cognitoPoolId string) ([]byte, error) {
	var url = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, cognitoPoolId)

	resp, err :=
		http.Get(url)

	if err != nil {
		fmt.Println("Error fetching public keys")
		return nil, errors.New("Error")
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return body, nil
}
