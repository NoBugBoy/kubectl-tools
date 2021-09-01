/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	dockerclient "github.com/docker/docker/client"
	"io"
	kd "kubedebug/docker"
	tcp "kubedebug/tcp"
	"os"
	"os/signal"
	"strconv"
)
func main() {
	containerId := os.Getenv("CONTAINERID")
	if containerId == strconv.Itoa(1) {
		os.Exit(1)
	}

	client,err  := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("connect docker error ")
		panic(err)
	}
	kd.PullDebugImg("nicolaka/netshoot:latest",client)

	ttyRes := kd.CreateDebugContainer(containerId,"nicolaka/netshoot:latest",client)

	container := kd.ConnectionTty(ttyRes)

	fmt.Println("waiting for plugin connect ..")
	conn,create := tcp.OpenSocket()
	<- create
	fmt.Println("plugin connect successful! create debug container..")

	go func() {
		for{
			str := <- container.Out
			_, err := conn.Conn.Write([]byte("\n" + str ))
			if err == io.EOF {
				killMe(ttyRes)
			}
		}
	}()
	go killMe(ttyRes)
	defer conn.Conn.Close()
	func() {
		for{
			b := make([]byte,10240)
			read, err := conn.Conn.Read(b)
			if err == io.EOF {
				killMe(ttyRes)
			}
			context := string(b[:read])
			if "exit\n" == context{
				kd.ClearAndClose(ttyRes.DebugContainerId,ttyRes.Client)
				conn.Conn.Write([]byte("closed"))
			}
			fmt.Println(string(b[:read]))
			container.In <- string(b[:read])
		}
	}()
}


func killMe(ttyRes *kd.TtyResponse){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	 <-c
	kd.ClearAndClose(ttyRes.DebugContainerId,ttyRes.Client)
	os.Exit(0)
}

