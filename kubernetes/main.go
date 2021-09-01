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
	var debug = &cobra.Command{
		Use:   "debug",
		Short: "debug pod name",
		Example: "kubectl utils debug {podName} {kubeConfigDir}",
		Run: func(cmd *cobra.Command, args []string) {
			debug(args)
		},
	}
	var rootCmd = &cobra.Command{Use: "kubectl"}
	rootCmd.AddCommand(debug)
	rootCmd.Execute()
}

func debug(args []string)  {
	// 0 podName 1 kubeConfig
	var plugin *Plugin
	if len(args) == 2 {
		plugin = runPlugin(args[0],args[1])
	}else{
		plugin = runPlugin(args[0],"/root/.kube/config")
	}
	if plugin == nil {
		fmt.Println("can not found pod" + args[0])
		os.Exit(1)
	}
	pod := GetPod("debug-k8s",plugin.Namespace,plugin.ContainerId,plugin.NodeName)

	startPod(plugin.Namespace,plugin.ClientSet,pod,plugin.NodeName)

	defer deletePod(plugin.Namespace,plugin.ClientSet)

	conn := connectionTcp(plugin.NodeName)

	fmt.Println("------------------------------------------")
	fmt.Println("- plugin connected ~ please input cmd >> -")
	fmt.Println("------------------------------------------")
	exit := make(chan string)
	go func() {
		for{
			b := make([]byte,10240)
			read, err := conn.Read(b)
			if string(b[:read]) == "closed"{
				fmt.Println("\ndebug container closed ...")
				exit <- "closed"
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
	stdin := exec.Command("stty", "erase","^H")
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

			time.Sleep(2 * time.Second)
			continue
		}
		// 尝试连接
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		conn.Write([]byte("\n"))
		return conn
	}
}
