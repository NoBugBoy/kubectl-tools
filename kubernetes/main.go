package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	var podName string
	var image string
	var config string
	var rcmd string
	var debug = &cobra.Command{
		Use:   "debug",
		Short: "debug target pod",
		Long: "debug target pod , default debug container is nicolaka/netshoot:latest, default kubeConfigDir /root/.kube/config",
		Example: "kubectl utils debug -c sh -m busybox -f /root/.kube/config -p nginx-xxx",
		Run: func(cmd *cobra.Command, args []string) {
			if podName == ""{
				fmt.Println("requires a pod name argument")
				os.Exit(0)
			}
			debug(podName,image,config,rcmd)
		},
	}
	var rootCmd = &cobra.Command{Use: "kubectl"}
	rootCmd.AddCommand(debug)
	debug.Flags().StringVarP(&podName, "podName", "p", "", "target pod name")
	debug.Flags().StringVarP(&image, "image", "m", "nicolaka/netshoot:latest", "debug container image name")
	debug.Flags().StringVarP(&config, "kubeConfigDir", "f", "/root/.kube/config", "k8s kube config dir")
	debug.Flags().StringVarP(&rcmd, "cmd", "c", "bash", "container attach command")
	rootCmd.Execute()
}

func debug(podName,image,config,cmd string)  {
	// 0 podName 1 kubeConfig
	var plugin *Plugin
	plugin = runPlugin(podName,config)

	if plugin == nil {
		fmt.Println("can not found pod" + podName)
		os.Exit(1)
	}
	pod := GetPod("debug-k8s",plugin.Namespace,plugin.ContainerId,image,plugin.NodeName,cmd)

	startPod(plugin.Namespace,plugin.ClientSet,pod,plugin.NodeName)

	defer deletePod(plugin.Namespace,plugin.ClientSet)

	conn := connectionTcp(plugin.NodeName)
	fmt.Println("connected waiting for pull debug image...")
	exit := make(chan string)
	pullOk := make(chan string)
	go func() {
		for{
			b := make([]byte,10240)
			read, err := conn.Read(b)
			if string(b[:read]) == "closed"{
				fmt.Println("\ndebug container closed ...")
				exit <- "closed"
			}
			if string(b[:read]) == "pulled"{
				pullOk <- "pulled"
				continue
			}
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}else{
				fmt.Print("\n" + string(b[:read]))
			}
		}
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c
		conn.Write([]byte("exit\n"))
		<-exit
		deletePod(plugin.Namespace,plugin.ClientSet)
		os.Exit(1)
	}()
	<- pullOk
	fmt.Println("")
	fmt.Println("------------------------------------------")
	fmt.Println("- plugin connected ~ please input cmd >> -")
	fmt.Println("------------------------------------------")
	stdin := exec.Command("stty", "erase ^H")
	stdin.Stdin = os.Stdin

	for{
		reader := bufio.NewReader(stdin.Stdin)
		text ,err := reader.ReadString('\n')
		if err != nil{
			fmt.Println(err)
		}
		conn.Write([]byte(text))

		if text == "exit\n" {
			<-exit
			deletePod(plugin.Namespace,plugin.ClientSet)
			break
		}

	}
}

func connectionTcp(nodeName string) net.Conn {
	fmt.Println("waiting for debug container ready ...")
	for  {
		tcpAddr, err := net.ResolveTCPAddr("tcp",nodeName + ":19675")
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		conn.Write([]byte("\n"))
		return conn
	}
}
