package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

type hasuraCustomClaim struct {
	UserId       string   `json:"x-hasura-user-id"`
	DefaultRole  string   `json:"x-hasura-default-role"`
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	Role         string   `json:"x-hasura-role"`
}

type customJWTClaim struct {
	HasuraNameSpace hasuraCustomClaim `json:"urlshortening/jwt/claims"`
	jwt.StandardClaims
}

var (
	jwtSecret    []byte
	jwtExpiresAt int
	jwtIssuer    string
)

func init() {
	jwtSecretString := GetEnvVar("JWT_SECRET")
	jwtSecret = []byte(jwtSecretString)

	expiry, err := strconv.Atoi(GetEnvVar("JWT_EXPIRE"))
	if err != nil {
		log.Debug().Err(err).Msg("Invalid value of JWT_EXPIRE env variable, falling back to default 3600")
		jwtExpiresAt = 3600
	} else {
		jwtExpiresAt = expiry
	}

	jwtIssuer = GetEnvVar("JWT_ISSUER")
}

func CreateJWTToken(userEmail string) (string, error) {

	// This application only uses static roles
	myHasuraClaims := hasuraCustomClaim{
		userEmail,
		"user",
		[]string{"user"},
		"user",
	}

	// Initialize custom claims
	myClaims := customJWTClaim{
		myHasuraClaims,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + int64(jwtExpiresAt),
			Issuer:    jwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	ss, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Debug().Err(err).Msgf("Error occurred while signing JWT token")
		return "", errors.New("failed to generate JWT token")
	}
	return ss, nil
}
