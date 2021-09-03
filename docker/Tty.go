package docker

import (
	"bufio"
	"github.com/docker/docker/api/types"
	c "github.com/docker/docker/client"
	"io"
	"kubedebug/tcp"
)

type Container struct {
	In chan string
	Out chan string
}

type TtyResponse struct {
  Response	types.HijackedResponse
  DebugContainerId string
  Client	*c.Client
}


func ConnectionTty(tty *TtyResponse,sockect *tcp.Socket,cmd string) *Container {
	container := &Container{
		In: make(chan string),
		Out: make(chan string),
	}
	go container.stdInput(&tty.Response)
	if cmd == "sh" || cmd == "bash"{
		go container.stdOutput(&tty.Response,nil)
	}else{
		go container.stdOutput(&tty.Response,sockect)
	}

	return container
}


func (c *Container) stdInput (hijack *types.HijackedResponse)  {
	for {
		cmd := <- c.In
		_, err := hijack.Conn.Write([]byte(cmd + "\n"))
		if err != nil {
			break
		}
	}
}

func (c *Container) stdOutput (hijack *types.HijackedResponse,socket *tcp.Socket)  {
	if socket == nil {
		scanner := bufio.NewScanner(hijack.Reader)
		for scanner.Scan(){
			in := scanner.Text()
			c.Out <- in
		}

		if err := scanner.Err(); err != nil {
			//pass
		}
	}else{
		io.Copy(socket.Conn,hijack.Reader)
	}

}
