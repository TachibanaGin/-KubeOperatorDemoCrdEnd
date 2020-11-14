package test

import (
	client_go "Crd-End/client-go"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNodePort(t *testing.T) {

	svcNodePortList, err := client_go.Clientset.CoreV1().Services("").List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		fmt.Println("error:" + err.Error())
	}
	for _,value := range svcNodePortList.Items{
		npresult := fmt.Sprint(value.Spec.Ports[0].NodePort)
		fmt.Println(npresult)
	}

}
