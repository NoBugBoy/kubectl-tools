package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
	"strings"
	"time"
)

func CreateDebugContainer(targetContainerId , image string, client *dockerclient.Client) *TtyResponse {
	ctx := context.Background()
	opts := types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
		Stdin:  true,
	}
	cf := &container.Config{
		Tty: true,
		OpenStdin: true,
		StdinOnce: true,
		Cmd: []string{"bash"},
		AttachStdin: true,
		AttachStderr: true,
		AttachStdout: true,
		Image: image,
	}
	fmt.Printf("target container id = %s",targetContainerId)
	targetId := fmt.Sprintf("container:%s",strings.Replace(targetContainerId,"docker://","",1))
	hostConf := &container.HostConfig{
		NetworkMode: container.NetworkMode(targetId),
		UsernsMode:  container.UsernsMode(targetId),
		IpcMode:     container.IpcMode(targetId),
		PidMode:     container.PidMode(targetId),
		CapAdd:      []string{"SYS_PTRACE", "SYS_ADMIN"},
	}
	cc ,e := client.ContainerCreate(ctx,cf,hostConf,nil,nil,"")
	if e != nil{
		fmt.Printf("create container error %s", e)
		panic(e)
	}
	err := client.ContainerStart(ctx, cc.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	attached, err := client.ContainerAttach(ctx, cc.ID, opts)
	if err != nil {
		fmt.Printf("attach container error %s",err)
		panic(err)
	}
	response := &TtyResponse{
		Response: attached,
		DebugContainerId: cc.ID,
		Client: client,
	}
	return response
}

func ClearAndClose(containerId string,client *dockerclient.Client)  {
	fmt.Println("debug over clear debug container start .")
	ctx := context.Background()
	t := time.Duration(10) * time.Millisecond
	fmt.Println(">>> waiting for stop debug container . <<<")
	err := client.ContainerStop(ctx, containerId, &t)
	if err != nil {
		fmt.Println(err)
	}
	ops := types.ContainerRemoveOptions{
		Force: true,
	}
	fmt.Println(">>> waiting for remove debugger container . <<<")
	err = client.ContainerRemove(ctx, containerId, ops)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(">>> clear and close successful . <<<")
}

func PullDebugImg(imageName string,client *dockerclient.Client){
	ctx := context.Background()
	fmt.Println("pull image" + imageName)
	c, err := client.ImagePull(ctx,imageName,types.ImagePullOptions{})
	defer func() {
		if c != nil {
			c.Close()
		}
	}()
	if err != nil {
		fmt.Printf("pull images erro %s \n",err)
	}


}