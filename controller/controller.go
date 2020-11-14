package controller

import (
	"Crd-End/model"
	"Crd-End/mysql"
	"Crd-End/redis"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)



func Login(c *gin.Context) {

	var user model.User
	var count int
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	fmt.Println(username,password)

	mysql.Db.Where("username = ? AND password = ?", username, password).Find(&user).Count(&count)

	if count > 0 {
		//生成token
		tokenInit := md5.New()
		io.WriteString(tokenInit,strconv.FormatInt(time.Now().Unix(), 10))
		token := fmt.Sprintf("%x", tokenInit.Sum(nil))  //w.Sum(nil)将w的hash转成[]byte格式
		fmt.Println(token)

		//将user转为json
		tokenValue,error := json.Marshal(user)
		if error != nil {
			println("JSON ERR:", error)
		}
		//存redis
		err := redis.Client.Set(token, tokenValue, 30 * time.Minute).Err()
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{
			"message" : user,
			"token" : token,
		})
	}else {
		c.JSON(http.StatusUnauthorized,gin.H{"message":"身份验证失败"})
	}

}

