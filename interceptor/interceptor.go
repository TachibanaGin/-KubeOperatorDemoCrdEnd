package interceptor

import (
	"Crd-End/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Interceptor() gin.HandlerFunc{
	return func(c *gin.Context){

		if c.Request.Header.Get("token") != "" && c.Request.Header.Get("token") != "null" {
			token := c.Request.Header.Get("token")
			fmt.Println("token:",token)
			value, err := redis.Client.Get(token).Result()
			if err != nil {
				c.Abort()
				c.JSON(http.StatusUnauthorized,gin.H{"message":"身份验证失败"})
				return
			}else if value == "redis: nil" {
				c.Abort()
				c.JSON(http.StatusUnauthorized,gin.H{"message":"身份验证失败"})
				return
			}else {
				redis.Client.Expire(token,30 * time.Minute)
				//c.JSON(http.StatusOK,gin.H{"message":"身份验证成功"})
				c.Next()
			}
		}else {
			c.Abort()
			c.JSON(401,gin.H{
				"message":"身份验证失败",
				"code": "401",
			})
			return
		}
	}
}

//		//这一部分可以替换成从session/cookie中获取，
//			username:=c.Query("username")
//			password:=c.Query("password")
//
//		if username=="ft" && password =="123"{
//			c.JSON(http.StatusOK,gin.H{"message":"身份验证成功"})
//			c.Next()  //该句可以省略，写出来只是表明可以进行验证下一步中间件，不写，也是内置会继续访问下一个中间件的
//		}else {
//			c.Abort()
//			c.JSON(http.StatusUnauthorized,gin.H{"message":"身份验证失败"})
//			return// return也是可以省略的，执行了abort操作，会内置在中间件defer前，return，写出来也只是解答为什么Abort()之后，还能执行返回JSON数据
//		}