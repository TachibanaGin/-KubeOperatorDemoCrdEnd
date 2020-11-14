package controller

import (
	client_go "Crd-End/client-go"
	"Crd-End/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	Z "github.com/go-redis/redis/v7"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math"
	"strconv"
	"time"
)

func GetFrontLis_Bak(c *gin.Context) {

	redisKey := "frontList"

	pagesize := c.Request.FormValue("pagesize")
	pagesizeInt ,err := strconv.Atoi(pagesize)
	if err !=nil {
		fmt.Println("erro为:", err)
	}
	pagesizeFloat := float64(pagesizeInt)
	page := c.Request.Form.Get("page")
	pageInt ,err := strconv.Atoi(page)
	if err !=nil {
		fmt.Println("erro为:", err)
	}
	keys , err := redis.Client.Exists(redisKey).Result()
	if err != nil {
		panic(err.Error())
	}
	if keys > 0 {
		redis.Client.Expire(redisKey,30 * time.Minute)
		count , err := redis.Client.ZCard(redisKey).Result()
		if err != nil {
			panic(err.Error())
		}
		countInt := int(count)
		countFloat := float64(countInt)
		pageCount := math.Ceil(countFloat / pagesizeFloat)
		if pageInt > 1 {
			frontList , _ := redis.Client.ZRange(redisKey,(int64(pageInt)-1)*int64(pagesizeInt),int64(pagesizeInt)*int64(pageInt)-1).Result()
			c.JSON(200, gin.H{
				"pagecount": pageCount,
				"count": countInt,
				"frontlist": frontList,
			})
		}

		frontList := redis.Client.ZRange(redisKey,0,int64(pagesizeInt)-1)
		c.JSON(200, gin.H{
			"pagecount": pageCount,
			"count": countInt,
			"frontlist": frontList,
		})
	}else{
		frontList , err := client_go.FrontClientset.FrontV1().Fronts("").List(context.TODO(),metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		for key,value := range frontList.Items{
			tokenValue,error := json.Marshal(value)
			if error != nil {
				println("JSON ERR:", error)
			}
			result := Z.Z{Score:float64(key),Member:tokenValue}
			_ , err := redis.Client.ZAdd(redisKey,&result).Result()
			if err != nil {
				panic(err.Error())
			}
		}
		redis.Client.Expire(redisKey,30 * time.Minute)

		count , err := redis.Client.ZCard(redisKey).Result()
		if err != nil {
			panic(err.Error())
		}
		countInt := int(count)
		countFloat := float64(countInt)
		pageCount := math.Ceil(countFloat / pagesizeFloat)

		if pageInt > 1 {
			frontListRedis , _ := redis.Client.ZRange(redisKey,(int64(pageInt)-1)*int64(pagesizeInt),int64(pagesizeInt)*int64(pageInt)-1).Result()
			c.JSON(200, gin.H{
				"pagecount": pageCount,
				"count": countInt,
				"frontlist": frontListRedis,
			})
		}

		frontListRedis := redis.Client.ZRange(redisKey,0,int64(pagesizeInt)-1)
		c.JSON(200, gin.H{
			"pagecount": pageCount,
			"count": countInt,
			"frontlist": frontListRedis,
		})
	}
}

func GetFront(c *gin.Context) {

	deployments, err := client_go.Clientset.AppsV1().Deployments("default").List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fronts, err := client_go.FrontClientset.FrontV1().Fronts("default").List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var resultlist_A []string
	for _,value := range deployments.Items{
		result := fmt.Sprintf("%+v",value.ObjectMeta.Name)
		fmt.Print(result + "\n")
		resultlist_A = append(resultlist_A, result)
	}
	var resultlist_B []string
	for _,value := range fronts.Items{
		result := fmt.Sprintf("%+v",value.ObjectMeta.Name)
		fmt.Print(result + "\n")
		resultlist_B = append(resultlist_B, result)
	}
	c.JSON(200, gin.H{
		"messageA": resultlist_A,
		//"messageB": resultlist_B,
		"messageB": fronts,
	})

}