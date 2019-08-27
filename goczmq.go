package main

import (
	"fmt"
	"time"

	"github.com/zeromq/goczmq"
)

func NewCZmqServer() Server {
	return &CZmqServer{
		end: make(chan struct{}, 1),
	}
}

type CZmqServer struct {
	socket *goczmq.Channeler
	end    chan struct{}
}

func (s *CZmqServer) LisenAndServe(endpoint string) (err error) {
	router := goczmq.NewRouterChanneler("tcp://" + endpoint)

	s.socket = router
	defer s.socket.Destroy()

	for {
		select {
		case <-s.end: // stop the server
			return nil
		case err := <-s.socket.ErrChan:
			return err
		default:
		}
		msg := <-s.socket.RecvChan // this is something like accept but return only the message
		go func(msg [][]byte) {
			if len(msg) >= 2 {
				id := msg[0]
				// fmt.Println(string(msg[1]))
				s.socket.SendChan <- [][]byte{id, msg[1]}
			}
		}(msg)
	}
}

func (s *CZmqServer) Close() (err error) {
	s.end <- struct{}{}
	close(s.end)
	s.socket.Destroy()
	return nil

}

func NewCZmqClient(endpoint string, i int) (Client, error) {
	dealer, err := goczmq.NewDealer("tcp://" + endpoint)
	if err != nil {
		return nil, err
	}
	dealer.SetOption(goczmq.SockSetIdentity(fmt.Sprintf("dealer-%v", i)))
	return &CZmqClient{
		socket: dealer,
	}, nil

}

type CZmqClient struct {
	socket *goczmq.Sock
}

func (c *CZmqClient) Send(b []byte) error {
	return c.socket.SendFrame(b, goczmq.FlagNone)
}

func (c *CZmqClient) Receive() ([]byte, error) {
	reply, err := c.socket.RecvMessage()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Router id:%s\n", string(reply[0]))
	return reply[0], nil
}

func (c *CZmqClient) Close() error {
	c.socket.Destroy()
	return nil
}

func NewCZmqClient2(endpoint string, i int) (Client, error) {

	dealer := goczmq.NewDealerChanneler("tcp://" + endpoint)
	select {
	case err := <-dealer.ErrChan:
		return nil, err
	case <-time.After(10 * time.Millisecond):
	}
	return &CZmqClient2{
		socket: dealer,
	}, nil

}

type CZmqClient2 struct {
	socket *goczmq.Channeler
}

func (c *CZmqClient2) Send(b []byte) error {
	select {
	case err := <-c.socket.ErrChan:
		return err
	default:
	}
	c.socket.SendChan <- [][]byte{b}
	// fmt.Printf("Send id:%s\n", string(b))

	return nil
}

func (c *CZmqClient2) Receive() ([]byte, error) {
	select {
	case err := <-c.socket.ErrChan:
		return nil, err
	default:
	}
	reply := <-c.socket.RecvChan
	// fmt.Printf("Reply id:%s\n", string(reply[0]))
	return reply[0], nil
}

func (c *CZmqClient2) Close() error {
	c.socket.Destroy()
	return nil
}
