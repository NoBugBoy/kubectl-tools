package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/watch"
	"k8s.io/kubernetes/pkg/client/conditions"
)

func GetPod(podName,ns,containerId,node string) *corev1.Pod {
	privileged := true
	pod := &corev1.Pod{
		TypeMeta: v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      podName,
			Namespace: ns,
		},
		Spec: corev1.PodSpec{
			HostPID:  true,
			NodeName: node,
			Containers: []corev1.Container{
				{
					Name:            "debug-k8s",
					Image:           "yujian1996/debug-k8s:v1.2",
					ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
					SecurityContext: &corev1.SecurityContext{
						Privileged: &privileged,
					},
					Env: []corev1.EnvVar{
						{
							Name: "CONTAINERID",
							Value: containerId,
						},
					},
					Ports: []corev1.ContainerPort{
						{
							Name:          "tcp",
							HostPort:      19675,
							ContainerPort: 19675,
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "docker",
							MountPath: "/var/run/docker.sock",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "docker",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/var/run/docker.sock",
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	//y, err := yaml.Marshal(pod)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("pod print in yaml: ", string(y))
	return pod
}

func startPod(ns string,clientset *kubernetes.Clientset,pod *corev1.Pod,nodeName string) (string,string) {
	ctx := context.Background()
	debugPod := clientset.CoreV1().Pods(ns)
	_, err := debugPod.Create(ctx, pod, v1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	watcher, err := clientset.CoreV1().Pods(pod.Namespace).Watch(ctx,v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	_, err = watch.UntilWithoutRetry(ctx, watcher, conditions.PodRunning)
	if err != nil {
		fmt.Printf("Error occurred while waiting for pod to run:  %v\n", err)
	}
	return nodeName,"debug-k8s"

}
var locks = 1
func deletePod(ns string,clientset *kubernetes.Clientset) {
	if locks == 1 {
		fmt.Println("clear debug-k8s pod ...")
		ctx := context.Background()
		err := clientset.CoreV1().Pods(ns).Delete(ctx,"debug-k8s",v1.DeleteOptions{})
		if err != nil {
			panic(err)
		}
		locks = -1
	}
}
