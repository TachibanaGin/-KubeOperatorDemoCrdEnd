package client_go

import (
	"Crd-End/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	home = homeDir()
	//dev
	//kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//dev
	//kubeconfig = flag.String("kubeconfig", filepath.Join("config", "admin.conf"), "(optional) absolute path to the kubeconfig file")
	//config, _ = clientcmd.BuildConfigFromFlags("", *kubeconfig)

	//proc
	config, _ = rest.InClusterConfig()

	Clientset, _ = kubernetes.NewForConfig(config)
	FrontClientset, _ = versioned.NewForConfig(config)
)
