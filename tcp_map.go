package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type EndpointInfo struct {
	Ipstr string `json:"ipstr"`
}

type Config struct {
	Host    EndpointInfo   `json:"host"`
	MapList []EndpointInfo `json:"maplist"`
}

var config Config

func init() {
	f, err := os.Open("tcp_map.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
	//设置随机数种子
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("welcome to tony tcp map!")
	fromaddr := config.Host.Ipstr
	fromlistener, err := net.Listen("tcp", fromaddr)

	if err != nil {
		log.Fatalf("Unable to listen on: %s, error: %s\n", fromaddr, err.Error())
	}
	defer fromlistener.Close()

	tcpMap(fromlistener)

}

func tcpMap(listener net.Listener) {
	for {
		con, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		//使用算法随机选取一个
		toIpStr := RandomSelect(config.MapList).Ipstr

		go func() {
			toCon, err := net.Dial("tcp", toIpStr)
			if err != nil {
				fmt.Printf("can not connect to %s", toIpStr)
				return
			}
			go handleConnection(con, toCon)
			go handleConnection(toCon, con)
		}()
	}
}

func RandomSelect(endPoints []EndpointInfo) EndpointInfo {
	index := rand.Intn(len(endPoints))
	return endPoints[index]
}

func handleConnection(r, w net.Conn) {
	defer r.Close()
	defer w.Close()

	var buffer = make([]byte, 100000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			break
		}

		n, err = w.Write(buffer[:n])
		if err != nil {
			break
		}
	}

}
