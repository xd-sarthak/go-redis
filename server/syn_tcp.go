package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/xd-sarthak/go-redis/config"
	"github.com/xd-sarthak/go-redis/core"
)

func readCommand(c net.Conn) (*core.RedisCmd, error) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	
	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondWithError(err error, c net.Conn) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c net.Conn) {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondWithError(err, c)
	}
}

func RunSyncTCPServer(){
	log.Println("starting a synchronous server on ",config.Host, config.Port)

	var connection_client int = 0

	address := config.Host+":"+strconv.Itoa(config.Port)

	listner, err := net.Listen("tcp",address)
	if err != nil {
		panic(err)
	}

	// first infinite loop, this loops for tcp connections which are ready to send data
	for{

		c,err := listner.Accept() // waiting for new client to connect -> blocking call
		if err != nil {
			panic(err)
		}

		connection_client += 1
		log.Println("client connected with address ", c.RemoteAddr(), "concurrent clients ", connection_client)

		// second infinite loop, we wait on the socket till it keeps sending data
		for {
			// over the socket we read the data sent by the client
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				connection_client -= 1
				log.Println("client disconnected with address ", c.RemoteAddr(), "concurrent clients ", connection_client)
				if err == io.EOF {
					break
				}
				log.Println("error reading command from client ", err)
			}

			// we process the command and send the response back to the client
			respond(cmd, c)
		}
	}
}