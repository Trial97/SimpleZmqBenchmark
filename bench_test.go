package main

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

const (
	endpoint  = "127.0.0.1:2567"
	msg       = "Hello"
	query     = `{"method":"ApierV2.SetBalance","params":[{"Tenant":"cgrates.org","Account":"1003","BalanceType":"*monetary","BalanceUUID":null,"BalanceID":"Bal2","Directions":null,"Value":10,"ExpiryTime":null,"RatingSubject":null,"Categories":null,"DestinationIds":null,"TimingIds":null,"Weight":10,"SharedGroups":null,"Blocker":null,"Disabled":null}],"id":0}`
	noClients = 100
)

func helpBenchmark(s Server, c []Client, b *testing.B) {
	clen := len(c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(clen)
		for _, cl := range c {
			go func(clnt Client) {
				defer wg.Done()
				if err := clnt.Send([]byte(msg)); err != nil {
					b.Error(err)
				}
				if out, err := clnt.Receive(); err != nil {
					b.Error(err)
				} else if !reflect.DeepEqual(string(out), msg) {
					b.Errorf("Expected: %q received: %q", msg, out) // not posible
				}
			}(cl)
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkStandard(b *testing.B) {
	s := NewStandardServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	var c []Client
	for i := 0; i < noClients; i++ {
		cl, err := NewStandardClient(endpoint)
		if err != nil {
			b.Fatal(err)
		}
		c = append(c, cl)
	}
	// fmt.Println("Started Standard")
	helpBenchmark(s, c, b)
	// fmt.Println("Done Standard")
	s.Close()
	for _, cl := range c {
		cl.Close()
	}
}

func BenchmarkZmq(b *testing.B) {
	s := NewZmqServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	var c []Client
	for i := 0; i < noClients; i++ {
		cl, err := NewZmqClient(endpoint, i)
		if err != nil {
			b.Fatal(err)
		}
		c = append(c, cl)
	}
	// fmt.Println("Started ZMQ")
	helpBenchmark(s, c, b)
	// fmt.Println("Done ZMQ")
	s.Close()
	for _, cl := range c {
		cl.Close()
	}
}

func BenchmarkZmq2(b *testing.B) {
	s := NewZmqServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	var c []Client
	for i := 0; i < noClients; i++ {
		cl, err := NewZmqClient2(endpoint, i)
		if err != nil {
			b.Fatal(err)
		}
		c = append(c, cl)
	}
	// fmt.Println("Started ZMQ")
	helpBenchmark(s, c, b)
	// fmt.Println("Done ZMQ")
	s.Close()
	for _, cl := range c {
		cl.Close()
	}
}

// func BenchmarkPZmq(b *testing.B) {
// 	s, err := NewPZmqServer()
// 	if err != nil {
// 		b.Fatal(err)
// 	}
// 	go s.LisenAndServe(endpoint)
// 	time.Sleep(time.Second)
// 	var c []Client
// 	for i := 0; i < noClients; i++ {
// 		cl, err := NewPZmqClient(endpoint)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		c = append(c, cl)
// 	}
// 	// fmt.Println("Started ZMQ")
// 	helpBenchmark(s, c, b)
// 	// fmt.Println("Done ZMQ")
// 	s.Close()
// 	for _, cl := range c {
// 		cl.Close()
// 	}
// }

// func BenchmarkPZmq2(b *testing.B) {
// 	s, err := NewPZmqServer()
// 	if err != nil {
// 		b.Fatal(err)
// 	}
// 	go s.LisenAndServe(endpoint)
// 	time.Sleep(time.Second)
// 	var c []Client
// 	for i := 0; i < noClients; i++ {
// 		cl, err := NewPZmqClient2(endpoint)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		c = append(c, cl)
// 	}
// 	// fmt.Println("Started ZMQ")
// 	helpBenchmark(s, c, b)
// 	// fmt.Println("Done ZMQ")
// 	s.Close()
// 	for _, cl := range c {
// 		cl.Close()
// 	}
// }

func BenchmarkCZmq(b *testing.B) {
	s := NewCZmqServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	var c []Client
	for i := 0; i < noClients; i++ {
		cl, err := NewCZmqClient(endpoint, i)
		if err != nil {
			b.Fatal(err)
		}
		c = append(c, cl)
	}
	// fmt.Println("Started ZMQ")
	helpBenchmark(s, c, b)
	// fmt.Println("Done ZMQ")
	s.Close()
	for _, cl := range c {
		cl.Close()
	}
}

func BenchmarkCZmq2(b *testing.B) {
	s := NewCZmqServer()
	go s.LisenAndServe(endpoint)
	time.Sleep(time.Second)
	var c []Client
	for i := 0; i < noClients; i++ {
		cl, err := NewCZmqClient2(endpoint, i)
		if err != nil {
			b.Fatal(err)
		}
		c = append(c, cl)
	}
	// fmt.Println("Started ZMQ")
	helpBenchmark(s, c, b)
	// fmt.Println("Done ZMQ")
	s.Close()
	for _, cl := range c {
		cl.Close()
	}
}
