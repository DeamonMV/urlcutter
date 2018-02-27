package main

import (
	"github.com/nats-io/go-nats-streaming"
	"fmt"
	"os"
	"os/signal"
	"time"
	"math/rand"
	"strconv"
	log "github.com/sirupsen/logrus"
	"urlcutter/generator/sub/checker"
)



func connect(clientId string) stan.Conn {

	conn, err := stan.Connect("test-cluster", clientId)
	if err != nil {
		log.Fatalf("Cann't connect : %v, make shure NATS Streaming server is running at %s", err)
	}
	return conn

}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}
//////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {

	number := random(0, 9)
clientId := strconv.Itoa(number)
log.Info("ClientID is: ", clientId)

msgConn := func(m *stan.Msg) {
	m.Ack()
	log.Info("Recive data from NATS: ", m)

	mapa := checker.Unmarshal(m.Data)
	codeint, err := checker.Checker(mapa)
	if err != ""{
		log.Warn(err)
	}
	fmt.Printf("####### CODE #########", codeint)

	//log.Info("Code is ", mapa)

}

conn := connect(string(clientId))

	aw, _ := time.ParseDuration("60s")
	_, err := conn.QueueSubscribe("foo", "bar", msgConn, stan.DeliverAllAvailable(), stan.SetManualAckMode(), stan.AckWait(aw), stan.DurableName("dur1"))
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
//	time.Sleep(10 * time.Millisecond)



	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Info("\nReceived an interrupt, unsubscribing and closing connection...\n\n")

			conn.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone



}