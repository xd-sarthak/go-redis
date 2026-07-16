package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xd-sarthak/go-redis/config"
)

// TODO: support non kv data structures
// TODO: support sync aof writes
func dumpKey(fp *os.File, key string, obj *Obj) {
	cmd := fmt.Sprintf("SET %s %s",key, obj.Value)
	tokens := strings.Split(cmd," ")
	fp.Write(Encode(tokens,false))
}

// DumpAllAOF dumps all the data in the store to the aof file
func DumpAllAOF() {
	fp, err := os.OpenFile(config.AOFFilePath,os.O_CREATE | os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Print("error",err)
		return
	}
	log.Println("dumping all data to aof file")
	for k,obj := range store {
		dumpKey(fp,k,obj)
	}
	log.Println("dumping all data to aof file completed")
	
}