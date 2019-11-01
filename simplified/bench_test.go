package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"

	"github.com/pebbe/zmq4"
	"github.com/zeromq/goczmq"
)

var (
	zmq_port  = 8226
	czmq_port = 8227
	m_port    = 8223
)

func init() {
	// start both servers
	go StartServer(m_port)
	go StartPebbeZmqServer(zmq_port)
	go StartCZmqServer(czmq_port)
}

//StartServer for server export
func StartServer(port int) {
	arith := new(MyServer)

	server := rpc.NewServer()
	server.Register(arith)

	// server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	sport := fmt.Sprintf(":%d", port)
	l, e := net.Listen("tcp", sport)
	if e != nil {
		log.Fatal("ListenError:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("AcceptError", err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		defer conn.Close()
	}
}

func BenchmarkStandard(b *testing.B) {
	is := fmt.Sprintf("localhost:%d", m_port)
	conn, err := net.Dial("tcp", is)

	if err != nil {
		b.Fatal("CreateNewClientError:", err)
	}
	c := jsonrpc.NewClient(conn)
	if c == nil {
		b.Fatal("user not connected")
	}
	defer c.Close()
	A, B := 1, 1
	var r int
	args := &Args{A, B}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// b.RunParallel(func(pb *testing.PB) {
		// for pb.Next() {
		if err := c.Call("MyServer.Sum", args, &r); err != nil {
			b.Error(err)
		}
		if r != A+B {
			b.Fatal("Sum is wrong")
		}
		// }
		// })

	}
	b.StopTimer()
}

func BenchmarkStandardParalel(b *testing.B) {
	is := fmt.Sprintf("localhost:%d", m_port)
	conn, err := net.Dial("tcp", is)

	if err != nil {
		b.Fatal("CreateNewClientError:", err)
	}
	c := jsonrpc.NewClient(conn)
	if c == nil {
		b.Fatal("user not connected")
	}
	defer c.Close()
	A, B := 1, 1
	var r int
	args := &Args{A, B}
	b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := c.Call("MyServer.Sum", args, &r); err != nil {
				b.Error(err)
			}
			if r != A+B {
				b.Fatal("Sum is wrong")
			}
		}
	})

	// }
	b.StopTimer()
}

func StartPebbeZmqServer(port int) (err error) {
	endpoint := fmt.Sprintf("127.0.0.1:%d", port)
	ctx, err := zmq4.NewContext()
	if err != nil {
		log.Fatal(err)
		return err
	}
	sock, err := ctx.NewSocket(zmq4.ROUTER)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = sock.Bind("tcp://" + endpoint)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer sock.Close()

	for {
		// fmt.Printf("Before rcv\n")
		msg, err := sock.RecvMessageBytes(0) // this is something like accept but return only the message
		// fmt.Printf("Router id:%v\n", string(msg[2]))
		if err != nil {
			log.Fatal(err)
			return err
		}
		// go func(msg [][]byte) {
		_, err = sock.SendMessage(string(msg[0]), "", "msg[2]") // hardcoded
		// _, err = sock.SendMessage(msg) // hardcoded
		// fmt.Printf("Router id:%v\n", err)

		if err != nil {
			log.Fatal(err)
			return err
		}
		// }(msg)
	}
}

func BenchmarkPZmq(b *testing.B) {
	endpoint := fmt.Sprintf("127.0.0.1:%d", zmq_port)
	client, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		b.Fatal(err)
	}
	err = client.Connect("tcp://" + endpoint)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// b.RunParallel(func(pb *testing.PB) {
		// for pb.Next() {
		// fmt.Printf("Send id:\n")
		if _, err := client.SendMessage("1000"); err != nil {
			b.Error(err)
		}
		// fmt.Printf("Send2\n")
		_, err := client.Recv(0)
		if err != nil {
			b.Error(err)
		}
		// fmt.Printf("RCV\n")
		// }
		// })

	}
	b.StopTimer()
}

func BenchmarkPZmqParralel(b *testing.B) {
	endpoint := fmt.Sprintf("127.0.0.1:%d", zmq_port)
	client, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		b.Fatal(err)
	}
	err = client.Connect("tcp://" + endpoint)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// fmt.Printf("Send id:\n")
			if _, err := client.SendMessage("1000"); err != nil {
				b.Error(err)
			}
			if _, err := client.SendMessage("1000"); err != nil {
				b.Error(err)
			}
			// fmt.Printf("Send2\n")
			_, err := client.Recv(0)
			if err != nil {
				b.Error(err)
			}
			// fmt.Printf("RCV\n")
		}
	})

	// }
	b.StopTimer()
}

func StartCZmqServer(port int) (err error) {
	endpoint := fmt.Sprintf("127.0.0.1:%d", port)
	router := goczmq.NewRouterChanneler("tcp://" + endpoint)

	defer router.Destroy()

	for {
		msg := <-router.RecvChan // this is something like accept but return only the message
		id := msg[0]
		// fmt.Println(string(msg[1]))
		router.SendChan <- [][]byte{id, []byte{}, msg[2]}
	}
}

func BenchmarkCZmq(b *testing.B) {
	endpoint := fmt.Sprintf("127.0.0.1:%d", czmq_port)
	dealer := goczmq.NewReqChanneler("tcp://" + endpoint)
	b.ResetTimer()
	go func() {
		err := <-dealer.ErrChan
		b.Fatal(err)
	}()
	for i := 0; i < b.N; i++ {
		// b.RunParallel(func(pb *testing.PB) {
		// for pb.Next() {
		// fmt.Printf("Send id:\n")
		dealer.SendChan <- [][]byte{[]byte("100")}
		<-dealer.RecvChan
		// fmt.Printf("RCV\n")
		// }
		// })

	}
	b.StopTimer()
}

func BenchmarkCZmqParalel(b *testing.B) {
	endpoint := fmt.Sprintf("127.0.0.1:%d", czmq_port)
	dealer := goczmq.NewReqChanneler("tcp://" + endpoint)
	b.ResetTimer()
	go func() {
		err := <-dealer.ErrChan
		b.Fatal(err)
	}()
	// for i := 0; i < b.N; i++ {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// fmt.Printf("Send id:\n")
			dealer.SendChan <- [][]byte{[]byte("100")}
			<-dealer.RecvChan
			// fmt.Printf("RCV\n")
		}
	})

	// }
	b.StopTimer()
}
