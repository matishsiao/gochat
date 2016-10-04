package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"syscall"
	"time"
)

func SetUlimit(number uint64) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Getting Rlimit ", err)
	}
	rLimit.Max = number
	rLimit.Cur = number
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Getting Rlimit ", err)
	}
	log.Println("set file limit done:", rLimit)
}

func decodeString(data string) string {
	keybyte, err := hex.DecodeString(string(data[:2]))
	if err != nil {
		return ""
	}
	key := uint(keybyte[0])
	decodeStr := ""
	if len(data)%2 != 0 {
		return ""
	}
	for i := 2; i < len(data)-1; i += 2 {
		value := data[i : i+2]
		byte, err := hex.DecodeString(value)
		if err != nil {
			continue
		}
		decodeStr += string(uint(byte[0]) ^ key)
	}
	return decodeStr

}

func encodeString(data string, key int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	if key == 0 || key > 99 {
		key = r.Intn(99)
	}
	encodeStr := fmt.Sprintf("%02s", fmt.Sprintf("%x", key))

	for _, v := range data {
		encodeStr += fmt.Sprintf("%02s", fmt.Sprintf("%x", int(v)^key))
	}
	return encodeStr
}
