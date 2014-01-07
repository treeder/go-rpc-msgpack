package main

import (
	"github.com/iron-io/rpctest/common"
	"fmt"
	"net/rpc"
	"log"
	"github.com/ugorji/go/codec"

	"net"
	"reflect"
)

func main() {

//	client := gobClient()
	client := msgpackClient()
	// Synchronous call
	args := &common.Args{7, 10}
	var reply int
	err := client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	quotient := new(common.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done	// will be equal to divCall
	// check errors, print, etc.
	fmt.Println("reply:" , replyCall.Reply, " / ", replyCall.Reply)

}

func gobClient() *rpc.Client {
	//RPC Communication (client side)
	s := "localhost:8080"
	client, err := rpc.Dial("tcp", s)
	if err != nil {
		log.Println("cannot dial:", err)
		panic(err)
	}
return client
}

func msgpackClient() *rpc.Client {
	var (
		mh codec.MsgpackHandle
	)

	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	//RPC Communication (client side)
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("dieal error:", err)
	}
	rpcCodec := codec.GoRpc.ClientCodec(conn, &mh)
	client := rpc.NewClientWithCodec(rpcCodec)
	return client
}
