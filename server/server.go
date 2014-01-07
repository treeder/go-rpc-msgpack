package main

import (
	"net/rpc"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ugorji/go/codec"
	"github.com/iron-io/rpctest/common"

	"reflect"
)

func main() {

	listenPort := 8080
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(listenPort))
	if err != nil {
		log.Println("cannot listen:", err)
		panic(err)
	}
	// todo: Add listener.Close() thing like beanstalkq.go?

	arith := new(common.Arith)

	//	listenGob(&listener, &arith)

	listenMsgpack(listener, arith)

}

func listenGob(listener net.Listener, reg *common.Arith) {
	//RPC Server
	rpcServer := rpc.NewServer()
	rpcServer.Register(reg)
	rpcServer.Accept(listener)
}

func listenMsgpack(listener net.Listener, reg *common.Arith) {
	// create and configure Handle
	var (
		mh codec.MsgpackHandle
	)

	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	rpcServer := rpc.NewServer()
	rpcServer.Register(reg)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err accepting:", err)
		}
		rpcCodec := codec.GoRpc.ServerCodec(conn, &mh)
		rpcServer.ServeCodec(rpcCodec)

	}
}
