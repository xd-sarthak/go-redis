package core

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")
var RESP_OK []byte = []byte("+OK\r\n")
var RESP_ZERO []byte = []byte(":0\r\n")
var RESP_ONE []byte = []byte(":1\r\n")
var RESP_MINUS_1 []byte = []byte(":-1\r\n")
var RESP_MINUS_2 []byte = []byte(":-2\r\n")

// evalPING is the implementation of PING
func evalPING(args []string) []byte {
	var b []byte

	if len(args) >= 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		b = Encode("PONG",true)
	} else {
		b = Encode(args[0],true)
	}

	return b
}

//evalSET is the implementation of SET command
// SET k v
func evalSET(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'set' command"), false)
	}

	var key,value string
	var expiration int64 = -1

	key = args[0]
	value = args[1]
	objType, objEncoding := deduceTypeEncoding(value)

	for i := 2; i < len(args); i++ {
		switch args[i] {
			case "EX","ex":
				i++;
				if i == len(args) {
					return Encode(errors.New("ERR syntax error"), false)
				}

				exDurationSec, err := strconv.ParseInt(args[i], 10, 64)
				if err != nil {
					return Encode(errors.New("ERR value is not an integer or out of range"), false)
				}
				expiration = exDurationSec * 1000
			default:
				return Encode(errors.New("ERR syntax error"), false)
		}
	}

	 // putting the key and value in the hashmap
	 Put(key, NewObj(value, expiration, objType, objEncoding))
	 return RESP_OK
}

// evalGET is the implementation of GET command
// GET k
// if expired or not found, return nil
func evalGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}

	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		return RESP_NIL
	}

	// if the object has expired, return nil
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli(){
		return RESP_NIL
	}

	return Encode(obj.Value, true)
}

// evalTTL is the implementation of TTL command
// TTL k
// if expired or not found, return -2
// if found but no expiration, return -1
func evalTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ttl' command"), false)
	}

	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		return RESP_MINUS_2
	}

	if obj.ExpiresAt == -1 {
		return RESP_MINUS_1
	}

	remainingMs := obj.ExpiresAt - time.Now().UnixMilli()
	if remainingMs <= 0 {
		return RESP_MINUS_2
	}

	return Encode(remainingMs/1000, false)
}

// if exist evalDEL returns how many keys were deleted, otherwise returns 0
func evalDEL(args []string) []byte {
	var countDeleted int = 0

	for _, key := range args {
		if ok := Del(key); ok {
			countDeleted++
		}
	}
	return Encode(countDeleted, false)
}

// evalEXPIRE is the implementation of EXPIRE command
// EXPIRE k seconds
// if the key does not exist, return 0
// if the key exists, set the expiration and return 1 
func evalEXPIRE(args []string) []byte {
	if len(args) <= 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'expire' command"), false)
	}

	var key string = args[0]
	expirationSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}

	// if the key does not exist, return 0
	obj := Get(key)
	if obj == nil {
		return RESP_ZERO
	}

	obj.ExpiresAt = time.Now().UnixMilli() + expirationSec*1000
	return RESP_ONE
}

func evalBGWRITEAOF(args []string) []byte {
	//TODO: make it async
	DumpAllAOF()
	return RESP_OK
}

func evalINCR(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'incr' command"), false)
	}

	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		Put(key, NewObj("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT))
	}

	if err := assertType(obj.TypeEncoding, OBJ_TYPE_STRING); err != nil {
		return Encode(err, false)
	}

	if err := assertEncoding(obj.TypeEncoding, OBJ_ENCODING_INT); err != nil {
		return Encode(err, false)
	}

	i,_ := strconv.ParseInt(obj.Value.(string), 10, 64)
	i++
	obj.Value = strconv.FormatInt(i, 10)

	return Encode(i, false)
}


func EvalAndRespond(cmds RedisCmds, c io.ReadWriter) {

	var response []byte
	buf := bytes.NewBuffer(response)

	for _,cmd := range cmds {
	switch cmd.Cmd {
	case "PING":
		 buf.Write(evalPING(cmd.Args))
	case "SET":
		buf.Write(evalSET(cmd.Args))
	case "GET":
		buf.Write(evalGET(cmd.Args))
	case "TTL":
		buf.Write(evalTTL(cmd.Args))
	case "DEL":
		buf.Write(evalDEL(cmd.Args))
	case "EXPIRE":
		buf.Write(evalEXPIRE(cmd.Args))
	case "BGWRITEAOF":
		buf.Write(evalBGWRITEAOF(cmd.Args))
	case "INCR":
		buf.Write(evalINCR(cmd.Args))
	default:
		buf.Write(evalPING(cmd.Args))
	}
}

	c.Write(buf.Bytes())
}