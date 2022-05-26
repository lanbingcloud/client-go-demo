package main

import (
	"log"

	"github.com/ChrisLo0751/client-go-demo/01/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1 获取配置文件
	// 集群外获取k8s配置文件
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// 集群内获取k8s配置文件
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("can't get config")
		}
		config = inClusterConfig
	}
	// 2 通过配置获取客户端对象
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("can't crate client")
	}
	// 3 获取informer factory
	factory := informers.NewSharedInformerFactory(clientset, 0)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()
	// 4 创建控制器
	controller := pkg.NewController(clientset, serviceInformer, ingressInformer)
	// 5 启动informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	// 6 启动控制器
	controller.Run(stopCh)
}
