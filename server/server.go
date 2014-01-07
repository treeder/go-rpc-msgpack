package server

import (
	"net/rpc"

	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ugorji/go/codec"

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

	arith := new(Arith)

	//	listenGob(&listener, &arith)

	listenMsgpack(listener, arith)

}

func listenGob(listener net.Listener, reg *Arith) {
	//RPC Server
	rpcServer := rpc.NewServer()
	rpcServer.Register(reg)
	rpcServer.Accept(listener)
}

func listenMsgpack(listener net.Listener, reg *Arith) {
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

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
