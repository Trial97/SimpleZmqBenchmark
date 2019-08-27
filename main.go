package main

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

const (
	NO_ITER = 10 ^ 5
)

const (
	endpoint  = "127.0.0.1:2567"
	msg       = "Hello"
	query     = `{"method":"ApierV2.SetBalance","params":[{"Tenant":"cgrates.org","Account":"1003","BalanceType":"*monetary","BalanceUUID":null,"BalanceID":"Bal2","Directions":null,"Value":10,"ExpiryTime":null,"RatingSubject":null,"Categories":null,"DestinationIds":null,"TimingIds":null,"Weight":10,"SharedGroups":null,"Blocker":null,"Disabled":null}],"id":0}`
	noClients = 100
)

func helpBenchmark(s Server, c Client, n int) {
	tStart := time.Now()
	for i := 0; i < n; i++ {
		if err := c.Send([]byte(msg)); err != nil {
			log.Fatal(err)
		}
		if out, err := c.Receive(); err != nil {
			log.Fatal(err)
		} else if !reflect.DeepEqual(string(out), msg) {
			log.Fatalf("Expected: %q received: %q", msg, out) // not posible
		}
	}
	fmt.Println(time.Now().Sub(tStart))
}

func Standard() {
	s := NewStandardServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	c, err := NewStandardClient(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Started Standard")
	helpBenchmark(s, c, NO_ITER)
	// fmt.Println("Done Standard")
	s.Close()
	c.Close()
}

func Zmq() {
	s := NewZmqServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	c, err := NewZmqClient(endpoint, 0)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Started ZMQ")
	helpBenchmark(s, c, NO_ITER)
	// fmt.Println("Done ZMQ")
	s.Close()
	c.Close()
}

func CZmq() {
	s := NewCZmqServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	c, err := NewCZmqClient(endpoint, 0)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Started ZMQ")
	helpBenchmark(s, c, NO_ITER)
	// fmt.Println("Done ZMQ")
	s.Close()
	c.Close()
}

func main() {
	fmt.Println("Standard")
	Standard()
	fmt.Println("ZMQ")
	Zmq()
	fmt.Println("CZMQ")
	CZmq()
}
