package main

import (
	"context"
	"fmt"

	"github.com/go-zeromq/zmq4"
)

func NewZmqServer() Server {
	return &ZmqServer{
		end:    make(chan struct{}, 1),
		socket: zmq4.NewRouter(context.Background(), zmq4.WithID(zmq4.SocketIdentity("router"))),
	}
}

type ZmqServer struct {
	socket zmq4.Socket
	end    chan struct{}
}

func (s *ZmqServer) LisenAndServe(endpoint string) (err error) {
	err = s.socket.Listen("tcp://" + endpoint)
	if err != nil {
		return err
	}
	defer s.socket.Close()

	for {
		select {
		case <-s.end: // stop the server
			return nil
		default:
		}
		msg, err := s.socket.Recv() // this is something like accept but return only the message
		if err != nil {
			return err
		}
		go func(msg zmq4.Msg) {
			if len(msg.Frames) >= 2 {
				// id := msg.Frames[0]
				// fmt.Println(string(msg.String()))
				err = s.socket.Send(zmq4.NewMsgFrom(msg.Frames...))
				if err != nil {
					return
				}
			}
		}(msg)
	}
}

func (s *ZmqServer) Close() (err error) {
	s.end <- struct{}{}
	close(s.end)
	return s.socket.Close()
}

func NewZmqClient(endpoint string, i int) (Client, error) {
	id := zmq4.SocketIdentity(fmt.Sprintf("dealer-%d", i))
	dealer := zmq4.NewDealer(context.Background(), zmq4.WithID(id))
	err := dealer.Dial("tcp://" + endpoint)
	if err != nil {
		return nil, err
	}
	return &ZmqClient{
		socket: dealer,
	}, nil

}

type ZmqClient struct {
	socket zmq4.Socket
}

func (c *ZmqClient) Send(b []byte) error {
	return c.socket.Send(zmq4.NewMsgFrom(b))
}

func (c *ZmqClient) Receive() ([]byte, error) {
	msg, err := c.socket.Recv()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Router id:%s\n", string(msg.String()))
	return msg.Frames[0], nil
}

func (c *ZmqClient) Close() error {
	return c.socket.Close()
}

func NewZmqClient2(endpoint string, i int) (Client, error) {
	dealer := zmq4.NewReq(context.Background())
	err := dealer.Dial("tcp://" + endpoint)
	if err != nil {
		return nil, err
	}
	return &ZmqClient2{
		socket: dealer,
	}, nil

}

type ZmqClient2 struct {
	socket zmq4.Socket
}

func (c *ZmqClient2) Send(b []byte) error {
	return c.socket.Send(zmq4.NewMsgFrom(b))
}

func (c *ZmqClient2) Receive() ([]byte, error) {
	msg, err := c.socket.Recv()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Router id:%s\n", string(msg.String()))
	return msg.Frames[0], nil
}

func (c *ZmqClient2) Close() error {
	return c.socket.Close()
}
