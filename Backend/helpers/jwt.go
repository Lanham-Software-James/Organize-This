package helpers

import (
	"errors"
	"fmt"
	"organize-this/config"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenHelper interface {
	VerifyToken(tokenString string, performValidation bool) (*jwt.Token, error)
	ExtractClaims(token *jwt.Token) (jwt.MapClaims, error)
}

type DefaultTokenHelper struct{}

func (h *DefaultTokenHelper) VerifyToken(tokenString string, performValidation bool) (*jwt.Token, error) {
	region := config.AWSRegion()
	userPool := config.CognitoUserPoolID()
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPool)
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPool)

	// Get the JWKS from Cognito
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, errors.New("Failed to get JWKS")
	}

	// Parse and validate the token
	var token *jwt.Token
	if performValidation {
		token, err = jwt.Parse(tokenString, jwks.Keyfunc,
			jwt.WithValidMethods([]string{"RS256"}),
			jwt.WithIssuer(issuer),
		)
	} else {
		token, err = jwt.Parse(tokenString, jwks.Keyfunc,
			jwt.WithValidMethods([]string{"RS256"}),
			jwt.WithIssuer(issuer),
			jwt.WithoutClaimsValidation(),
		)
	}

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}

	return token, nil
}

func (h *DefaultTokenHelper) ExtractClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid Token Claims")
	}

	return claims, nil
}
