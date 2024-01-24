package main

import (
	"github.com/nats-io/stan.go"
	"os"
)

func publish(path string, sc stan.Conn) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = sc.Publish("addNewOrder", b)
	if err != nil {
		return
	}
}

func main() {
	sc, err := stan.Connect("test-cluster", "publish", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		panic(err)
	}
	defer func(sc stan.Conn) {
		err := sc.Close()
		if err != nil {

		}
	}(sc)

	publish("model.json", sc)
	publish("m2.json", sc)
}
