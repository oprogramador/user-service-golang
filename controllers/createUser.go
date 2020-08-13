package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/models"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func CreateUser(userManager UserManager) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		reqBody, _ := ioutil.ReadAll(ginContext.Request.Body)
		var user models.User
		err := json.Unmarshal(reqBody, &user)
		if err != nil {
			re := regexp.MustCompile(`[A-Za-z.]* of type [A-Za-z]*`)
			ginContext.String(400, strings.ReplaceAll(string(re.Find([]byte(err.Error()))), "of type", "should be of type"))
			return
		}

		err = userManager.Save(&user)
		if err != nil {
			log.Println(err)
			ginContext.String(500, "")
		}
		ginContext.JSON(201, user)
	}
}
