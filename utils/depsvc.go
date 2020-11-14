package utils

import (
	front "Crd-End/pkg/apis/front/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func NewDepForCR(cr,result *front.Front) *appsv1.Deployment {
	replicas := cr.Spec.Replicas
	containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range cr.Spec.Ports {
		cport := corev1.ContainerPort{}
		cport.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}
	z := true
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				//*metav1.NewControllerRef(cr, schema.GroupVersionKind{
				//	Group: appsv1.SchemeGroupVersion.Group,
				//	Version: appsv1.SchemeGroupVersion.Version,
				//	Kind: "Front",
				//}),
				{
					"apps/v1",
					"Front",
					result.ObjectMeta.Name,
					result.ObjectMeta.UID,
					&z,
					&z,
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": cr.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": cr.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  cr.Name + "-" + "pod",
							Image: cr.Spec.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: containerPorts,
						},
					},
				},
			},
		},
	}
	return deployment
}

func NewSvcForCR(cr,result *front.Front) *corev1.Service {
	containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range cr.Spec.Ports {
		cport := corev1.ContainerPort{}
		cport.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}
	z := true
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				//*metav1.NewControllerRef(cr, schema.GroupVersionKind{
				//	Group: appsv1.SchemeGroupVersion.Group,
				//	Version: appsv1.SchemeGroupVersion.Version,
				//	Kind: "Front",
				//}),
				{
					"apps/v1",
					"Front",
					result.ObjectMeta.Name,
					result.ObjectMeta.UID,
					&z,
					&z,
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app" : cr.Name,
			},
			Ports: cr.Spec.Ports,
			Type: corev1.ServiceTypeNodePort,
		},
	}
	return service
}

//dep := depsvc.NewDepForCR(front,result)
//resultDep , err := client_go.Clientset.AppsV1().Deployments("default").Create(context.TODO(),dep,metav1.CreateOptions{})
//if err != nil {
//	//panic(err)
//	c.JSON(200, gin.H{
//		"message" : err.Error(),
//	})
//	return
//}
//fmt.Printf("Created Deployment %q.\n", resultDep.GetObjectMeta().GetName())
//
//svc := depsvc.NewSvcForCR(front,result)
//resultSvc , err := client_go.Clientset.CoreV1().Services("default").Create(context.TODO(),svc,metav1.CreateOptions{})
//if err != nil {
//	//panic(err)
//	c.JSON(200, gin.H{
//		"message" : err.Error(),
//	})
//	return
//}
//fmt.Printf("Created Service %q.\n", resultSvc.GetObjectMeta().GetName())
