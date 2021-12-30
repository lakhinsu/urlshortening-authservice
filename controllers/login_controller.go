package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/models"
	"github.com/lakhinsu/urlshortening-authservice/utils"
	"github.com/rs/zerolog/log"
)

func Login(c *gin.Context) {
	var requestData models.User
	binderr := c.ShouldBindJSON(&requestData)

	request_id := c.GetString("x-request-id")

	if binderr != nil {
		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	results, status, err := utils.GetLdapUser(requestData.Username, requestData.Password)
	if err != nil {
		if status == http.StatusUnauthorized {
			log.Error().Err(err).Str("request_id", request_id).
				Msg("Error occurred while authenticating user with LDAP server")
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid username or password",
			})
			return
		}
		log.Error().Err(err).Str("request_id", request_id).
			Msg("Error occurred while authenticating user with LDAP server")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	jwtToken, err := utils.CreateJWTToken(results["email"])
	if err != nil {
		log.Error().Err(err).Str("request_id", request_id).
			Msg("Error occurred while creating a JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}
