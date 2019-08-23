package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

type Server interface {
	LisenAndServe(endpoint string) error
	Close() error
}

type Client interface {
	Send([]byte) error
	Receive() ([]byte, error)
	Close() error
}

const (
	MaxByteBuffer = 400 // do not send over this
)

func NewStandardServer() Server {
	return &StandardServer{
		end: make(chan struct{}, 1),
	}
}

type StandardServer struct {
	end  chan struct{}
	conn net.Listener
}

func (s *StandardServer) LisenAndServe(endpoint string) (err error) {
	s.conn, err = net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}
	for {
		select {
		case <-s.end: // stop the server
			return nil
		default:
		}
		conn, err := s.conn.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			b := make([]byte, MaxByteBuffer)
			defer c.Close()

			for {
				n, err := c.Read(b)
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("<%s> Read error: %s", c.RemoteAddr(), err)
					return
				}
				if n == 0 {
					return
				}
				n, err = c.Write(b)
				if err != nil {
					log.Printf("<%s> Write error: %s", c.RemoteAddr(), err)
					return
				}
			}
		}(conn)
	}
}

func (s *StandardServer) Close() (err error) {
	s.end <- struct{}{}
	close(s.end)
	return s.conn.Close()
}

func NewStandardClient(endpoint string) (Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	return &StandardClient{
		conn: conn,
	}, nil

}

type StandardClient struct {
	conn net.Conn
}

func (c *StandardClient) Send(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *StandardClient) Receive() ([]byte, error) {
	b := make([]byte, MaxByteBuffer)
	_, err := c.conn.Read(b)
	if err != nil {
		return nil, err
	}
	return bytes.Trim(b, "\x00"), nil
	// return b, nil
}

func (c *StandardClient) Close() error {
	return c.conn.Close()
}
