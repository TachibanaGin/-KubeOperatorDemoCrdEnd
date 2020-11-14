package controller

import (
	client_go "Crd-End/client-go"
	front "Crd-End/pkg/apis/front/v1"
	"Crd-End/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func GetFrontList(c *gin.Context) {

	selectFront := c.Request.FormValue("frontname")
	namespaces := c.Request.Header.Get("user")
	fmt.Println(selectFront,namespaces)
	if selectFront != ""{
		front , err := client_go.FrontClientset.FrontV1().Fronts(namespaces).Get(context.TODO(),selectFront,metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err){
				c.JSON(200, gin.H{
					"frontList": "",
				})
			}else{
				panic(err.Error())
			}
		}
		c.JSON(200, gin.H{
			"frontList": front,
		})
	}else {
		frontList , err := client_go.FrontClientset.FrontV1().Fronts(namespaces).List(context.TODO(),metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		c.JSON(200, gin.H{
			"frontList": frontList.Items,
		})
	}

}

func AddFront(c *gin.Context) {

	replicas := c.Request.FormValue("replicas")
	replicasInt ,err := strconv.Atoi(replicas)
	if err !=nil {
		fmt.Println("erro为:", err)
	}
	replicasInt32 := int32(replicasInt)

	port := c.Request.FormValue("port")
	portInt ,err := strconv.Atoi(port)
	if err !=nil {
		fmt.Println("erro为:", err)
	}
	portInt32 := int32(portInt)

	nodeport := c.Request.FormValue("nodeport")
	nodeportInt ,err := strconv.Atoi(nodeport)
	if err !=nil {
		fmt.Println("erro为:", err)
	}
	nodeportInt32 := int32(nodeportInt)

	frontname := c.Request.FormValue("frontname")
	namespaces := c.Request.Header.Get("user")

	front := &front.Front{
		ObjectMeta: metav1.ObjectMeta{
			Name: frontname,
			Namespace: namespaces,
		},
		Spec: front.FrontSpec{
			&replicasInt32,
			c.Request.FormValue("image"),
			[]corev1.ServicePort{
				{
					Port: 		portInt32,
					TargetPort:	intstr.FromInt(portInt) ,
					NodePort:	nodeportInt32,
				},
			},
		},
	}
	//frontname := c.Request.FormValue("frontname")
	//image := c.Request.FormValue("image")
	//port := c.Request.FormValue("port")
	//nodeport := c.Request.FormValue("nodeport")

	//判断是否存在dep
	depList, err := client_go.Clientset.AppsV1().Deployments(namespaces).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
	}
	for _,value := range depList.Items{
		result := fmt.Sprintf("%+v",value.ObjectMeta.Name)
		if frontname == result {
			fmt.Println("error:" + "Deployment-"+ "["+result+"]" + " already exists,请修改Frontname")
			c.JSON(200, gin.H{
				"message" : "Deployment-"+ "["+result+"]" + " already exists,请修改Frontname",
			})
			return
		}
	}
	//判断是否存在svc和nodeport
	svcList, err := client_go.Clientset.CoreV1().Services(namespaces).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
	}
	for _,value := range svcList.Items{
		result := fmt.Sprintf("%+v",value.ObjectMeta.Name)
		if frontname == result {
			fmt.Println("error:" + "Service-"+ "["+result+"]" + " already exists,请修改Frontname")
			c.JSON(200, gin.H{
				"message" : "Service-"+ "["+result+"]" + " already exists,请修改Frontname",
			})
			return
		}
	}
	svcNodePortList, err := client_go.Clientset.CoreV1().Services("").List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
	}
	for _,value := range svcNodePortList.Items{
		for _,nplist := range value.Spec.Ports{
			npresult := fmt.Sprint(nplist.NodePort)
			if nodeport == npresult {
				fmt.Println("error:" + "NodePort-"+ "["+nodeport+"]" + " already exists,请修改对外端口")
				c.JSON(200, gin.H{
					"message" : "NodePort-"+ "["+nodeport+"]" + " already exists,请修改对外端口",
				})
				return
			}
		}
		//npresult := fmt.Sprint(value.Spec.Ports[0].NodePort)
		//if nodeport == npresult {
		//	fmt.Println("error:" + "NodePort-"+ "["+nodeport+"]" + " already exists,请修改对外端口")
		//	c.JSON(200, gin.H{
		//		"message" : "NodePort-"+ "["+nodeport+"]" + " already exists,请修改对外端口",
		//	})
		//	return
		//}
	}

	result , err := client_go.FrontClientset.FrontV1().Fronts(namespaces).Create(context.TODO(),front,metav1.CreateOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}
	fmt.Printf("Created Front %q.\n", result.GetObjectMeta().GetName())

	c.JSON(200, gin.H{
		"success" : true,
	})

}

func DelFront(c *gin.Context) {

	namespaces := c.Request.Header.Get("user")
	frontname := c.Request.FormValue("frontname")
	err := client_go.FrontClientset.FrontV1().Fronts(namespaces).Delete(context.TODO(),frontname,metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}
	fmt.Printf("Dekete Front %q.\n", frontname)

	c.JSON(200, gin.H{
		"success" : true,
	})

}

func UpdateFront(c *gin.Context){

	replicas := c.Request.FormValue("replicas")
	namespaces := c.Request.Header.Get("user")
	replicasInt ,err := strconv.Atoi(replicas)
	if err !=nil {
		fmt.Println("erro为:", err)
	}
	replicasInt32 := int32(replicasInt)

	port := c.Request.FormValue("port")
	portInt ,err := strconv.Atoi(port)
	if err !=nil {
		fmt.Println("erro为:", err)
	}
	portInt32 := int32(portInt)

	frontname := c.Request.FormValue("frontname")


	oldFront , err := client_go.FrontClientset.FrontV1().Fronts(namespaces).Get(context.TODO(),frontname,metav1.GetOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}

	oldFront.Spec.Ports[0].Port = portInt32
	oldFront.Spec.Ports[0].TargetPort = intstr.FromInt(portInt)
	oldFront.Spec.Image = c.Request.FormValue("image")
	oldFront.Spec.Replicas = &replicasInt32

	result , err := client_go.FrontClientset.FrontV1().Fronts(namespaces).Update(context.TODO(),oldFront,metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}
	fmt.Printf("Update Front %q.\n", result.GetObjectMeta().GetName())

	c.JSON(200, gin.H{
		"success" : true,
	})


}

func GetFrontOne(c *gin.Context){

	frontname := c.Request.FormValue("frontname")
	namespaces := c.Request.Header.Get("user")
	front , err := client_go.FrontClientset.FrontV1().Fronts(namespaces).Get(context.TODO(),frontname,metav1.GetOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}
	svc , err := client_go.Clientset.CoreV1().Services(namespaces).Get(context.TODO(),frontname,metav1.GetOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}
	deploy , err := client_go.Clientset.AppsV1().Deployments(namespaces).Get(context.TODO(),frontname,metav1.GetOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}
	replica , err := client_go.Clientset.AppsV1().ReplicaSets(namespaces).List(context.TODO(),metav1.ListOptions{LabelSelector:utils.LabelsForRs(frontname)})
	if err != nil {
		fmt.Println("error:" + err.Error())
		c.JSON(200, gin.H{
			"message" : err.Error(),
		})
		return
	}
	pods , err := client_go.Clientset.CoreV1().Pods(namespaces).List(context.TODO(),metav1.ListOptions{LabelSelector:utils.LabelsForRs(frontname)})
	fmt.Println(frontname)
	c.JSON(200, gin.H{
		"success" : true,
		"front": front,
		"deploy": deploy,
		"svc": svc,
		"replica": replica,
		"pods": pods.Items,
	})
}

func GetPodDetail(c *gin.Context){

	podName := c.Request.FormValue("podname")
	containerName := c.Request.FormValue("containername")
	namespaces := c.Request.Header.Get("user")

	logs := client_go.Clientset.CoreV1().Pods(namespaces).GetLogs(podName,&corev1.PodLogOptions{Container:containerName})
	resp, err := logs.DoRaw(context.TODO())
	if err != nil {
		fmt.Println("error:" + err.Error())
	}
	//fmt.Println(string(resp))
	c.JSON(200, gin.H{
		"success" : true,
		"logs": string(resp),
	})

}

func GetDescribePod(c *gin.Context){

	podName := c.Request.FormValue("podname")
	//containerName := c.Request.FormValue("containername")
	namespaces := c.Request.Header.Get("user")

	describe ,err := client_go.Clientset.CoreV1().Pods(namespaces).Get(context.TODO(),podName,metav1.GetOptions{})
	//logs := client_go.Clientset.CoreV1().Pods("default").GetLogs(podName,&corev1.PodLogOptions{Container:containerName})
	if err != nil {
		fmt.Println("error:" + err.Error())
	}
	//fmt.Println(string(resp))
	c.JSON(200, gin.H{
		"success" : true,
		"logs": describe,
	})

}
