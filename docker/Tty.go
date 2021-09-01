package docker

import (
	"bufio"
	"github.com/docker/docker/api/types"
	c "github.com/docker/docker/client"
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


func ConnectionTty(tty *TtyResponse) *Container {
	container := &Container{
		In: make(chan string),
		Out: make(chan string),
	}
	go container.stdInput(&tty.Response)
	go container.stdOutput(&tty.Response)
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

func (c *Container) stdOutput (hijack *types.HijackedResponse)  {
	scanner := bufio.NewScanner(hijack.Reader)

	for scanner.Scan(){
		in := scanner.Text()
		c.Out <- in
	}

	if err := scanner.Err(); err != nil {
		//pass
	}
}
