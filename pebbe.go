package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/pebbe/zmq4"
)

// for the momment desn't work because
// Assertion failed: !_current_out (src/router.cpp:191)

func NewPZmqServer() (Server, error) {
	ctx, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	sock, err := ctx.NewSocket(zmq4.ROUTER)
	if err != nil {
		return nil, err
	}
	return &PZmqServer{
		end:    make(chan struct{}, 1),
		ctx:    ctx,
		socket: sock,
	}, nil
}

type PZmqServer struct {
	ctx    *zmq4.Context
	socket *zmq4.Socket
	end    chan struct{}
}

func (s *PZmqServer) LisenAndServe(endpoint string) (err error) {
	err = s.socket.Bind("tcp://" + endpoint)
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
		msg, err := s.socket.RecvMessageBytes(0) // this is something like accept but return only the message
		// fmt.Printf("Router id:%v\n", msg[0])
		if err != nil {
			return err
		}
		go func(msg [][]byte) {
			_, err = s.socket.SendMessage(string(msg[0]), "10") // hardcoded
			fmt.Printf("Router id:%v\n", err)

			if err != nil {
				return
			}
		}(msg)
	}
}

func (s *PZmqServer) Close() (err error) {
	s.end <- struct{}{}
	close(s.end)
	return s.socket.Close()
}

func NewPZmqClient(endpoint string) (Client, error) {

	client, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		return nil, err
	}
	err = client.Connect("tcp://" + endpoint)
	if err != nil {
		return nil, err
	}
	return &PZmqClient{
		socket: client,
	}, nil

}

type PZmqClient struct {
	socket *zmq4.Socket
}

func (c *PZmqClient) Send(b []byte) error {
	_, err := c.socket.SendBytes(b, 0)
	// fmt.Printf("Send id:%s\n", string(b))
	return err
}

func (c *PZmqClient) Receive() ([]byte, error) {
	poller := zmq4.NewPoller()
	poller.Add(c.socket, zmq4.POLLIN)
	polled, err := poller.Poll(10 * time.Second)
	if err != nil {
		return nil, err
	}
	if len(polled) == 1 {

		reply, err := c.socket.RecvBytes(0)
		if err != nil {
			return nil, err
		}
		// fmt.Printf("Router id:%s\n", string(reply))
		return reply, nil
	}
	return nil, errors.New("Time out")
}

func (c *PZmqClient) Close() error {
	return c.socket.Close()
}

func NewPZmqClient2(endpoint string) (Client, error) {

	client, err := zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		return nil, err
	}
	err = client.Connect("tcp://" + endpoint)
	if err != nil {
		return nil, err
	}
	return &PZmqClient{
		socket: client,
	}, nil

}

type PZmqClient2 struct {
	socket *zmq4.Socket
}

func (c *PZmqClient2) Send(b []byte) error {
	_, err := c.socket.SendBytes(b, 0)
	// fmt.Printf("Send id:%s\n", string(b))
	return err
}

func (c *PZmqClient2) Receive() ([]byte, error) {
	// poller := zmq4.NewPoller()
	// poller.Add(c.socket, zmq4.POLLIN)
	// polled, err := poller.Poll(10 * time.Second)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(polled) == 1 {

	reply, err := c.socket.RecvBytes(0)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Router id:%s\n", string(reply))
	return reply, nil
	// }
	// return nil, errors.New("Time out")
}

func (c *PZmqClient2) Close() error {
	return c.socket.Close()
}
