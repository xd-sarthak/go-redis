package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")

// evalPING is the implementation of PING
func evalPING(args []string, c io.ReadWriter) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		b = Encode("PONG",true)
	} else {
		b = Encode(args[0],true)
	}

	_, err := c.Write(b)
	return err
}

//evalSET is the implementation of SET command
// SET k v
func evalSET(args []string, c io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("ERR wrong number of arguments for 'set' command")
	}

	var key,value string
	var expiration int64 = -1

	key = args[0]
	value = args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
			case "EX","ex":
				i++;
				if i == len(args) {
					return errors.New("ERR syntax error")
				}

				exDurationSec, err := strconv.ParseInt(args[i], 10, 64)
				if err != nil {
					return errors.New("ERR value is not an integer or out of range")
				}
				expiration = exDurationSec * 1000
			default:
				return errors.New("ERR syntax error")
		}
	}

	 // putting the key and value in the hashmap
	 Put(key, NewObj(value, expiration))
	 c.Write([]byte("+OK\r\n"))
	 return nil
}

// evalGET is the implementation of GET command
// GET k
// if expired or not found, return nil
func evalGET(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'get' command")
	}

	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		c.Write(RESP_NIL)
		return nil
	}

	// if the object has expired, return nil
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli(){
		c.Write(RESP_NIL)
		return nil
	}

	c.Write(Encode(obj.Value, true))
	return nil
}

// evalTTL is the implementation of TTL command
// TTL k
// if expired or not found, return -2
// if found but no expiration, return -1
func evalTTL(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'ttl' command")
	}

	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		c.Write([]byte(":-2\r\n"))
		return nil
	}

	if obj.ExpiresAt == -1 {
		c.Write([]byte(":-1\r\n"))
		return nil
	}

	remainingMs := obj.ExpiresAt - time.Now().UnixMilli()
	if remainingMs <= 0 {
		c.Write([]byte(":-2\r\n"))
		return nil
	}

	c.Write(Encode(remainingMs/1000, false))
	return nil
}


func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	case "SET":
		return evalSET(cmd.Args,c)
	case "GET":
		return evalGET(cmd.Args,c)
	case "TTL":
		return evalTTL(cmd.Args,c)
	default:
		return evalPING(cmd.Args, c)
	}
}