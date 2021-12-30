package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/models"
	"github.com/lakhinsu/urlshortening-authservice/utils"
)

func Login(c *gin.Context) {
	var requestData models.User
	binderr := c.ShouldBindJSON(&requestData)

	if binderr != nil {
		fmt.Println(binderr.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	results, err := utils.GetLdapUser(requestData.Username, requestData.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	jwtToken, err := utils.CreateJWTToken(results["email"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong while authenticating user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}
