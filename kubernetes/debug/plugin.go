package debug

import (
	"context"
	"fmt"
	v "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Plugin struct {
	ClientSet *kubernetes.Clientset
	NodeName string
	ContainerId string
	Namespace string
}

func runPlugin(podName , kubeConfigDir string) *Plugin {
	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigDir)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	allNs,err := clientset.CoreV1().Namespaces().List(ctx,v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _,item := range allNs.Items{
		pod, err := clientset.CoreV1().Pods(item.Name).Get(ctx,podName, v1.GetOptions{})
		if err != nil{
			continue
		}
		if pod != nil {
			return &Plugin{
				ClientSet:   clientset,
				NodeName:    pod.Spec.NodeName,
				ContainerId: getContainerIdByPod(pod.Spec.Containers[0].Name,pod),
				Namespace:   item.Name,
			}
		}
	}
	return nil
}


func getContainerIdByPod(containerName string , pod *v.Pod) string {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name != containerName {
			continue
		}
		if containerStatus.State.Running == nil {
			fmt.Println("container " + containerName + "not running")
		}
		return containerStatus.ContainerID
    }
    return ""
}
