package main

import (
	"fmt"
	"github.com/spf13/cobra"
	back "kubedebug/kubernetes/backup"
	d "kubedebug/kubernetes/debug"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{Use: "kubectl"}
	//debug
	rootCmd.AddCommand(debugTools())
	//backup etcd cluster
	rootCmd.AddCommand(backupEtcdTools())
	//todo 升级集群
	//todo 修改证书过期时间
	rootCmd.Execute()
}

func backupEtcdTools() *cobra.Command{
	var outPutDir,kubeConfig,confiMapName string
	var backupCobra = &cobra.Command{
		Use:   "backup",
		Short: "backup etcd snapshot data",
		Long: "backup etcd snapshot data , you know for safe data , but use only in kubedm environment",
		Example: "kubectl utils backup -o /var/lib/etcd",
		Run: func(cmd *cobra.Command, args []string) {
			back.Backup(outPutDir,kubeConfig,confiMapName)
		},
	}
	backupCobra.Flags().StringVarP(&outPutDir, "outPutDir", "o", "/var/lib/etcd", "etcd snapshot data output dir,generate names based on time")
	backupCobra.Flags().StringVarP(&confiMapName, "kubeadmConfig", "c", "kubeadm-config", "k8s cluster init kubeadm-config")
	backupCobra.Flags().StringVarP(&kubeConfig, "kubeConfigDir", "f", "/root/.kube/config", "k8s kube config dir")
	return backupCobra
}

func debugTools() *cobra.Command {
	var podName,image,config,rcmd string
	var debugCobra = &cobra.Command{
		Use:   "debug",
		Short: "debug target pod",
		Long: "debug target pod , default debug container is nicolaka/netshoot:latest, default kubeConfigDir /root/.kube/config",
		Example: "kubectl utils debug -c sh -m busybox -f /root/.kube/config -p nginx-xxx",
		Run: func(cmd *cobra.Command, args []string) {
			if podName == ""{
				fmt.Println("requires a pod name argument")
				os.Exit(0)
			}
			d.Debug(podName,image,config,rcmd)
		},
	}
	debugCobra.Flags().StringVarP(&podName, "podName", "p", "", "target pod name")
	debugCobra.Flags().StringVarP(&image, "image", "m", "nicolaka/netshoot:latest", "debug container image name")
	debugCobra.Flags().StringVarP(&config, "kubeConfigDir", "f", "/root/.kube/config", "k8s kube config dir")
	debugCobra.Flags().StringVarP(&rcmd, "cmd", "c", "bash", "container attach command")
	return debugCobra
}
