package stan_operation

import (

	"github.com/nats-io/go-nats-streaming"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"

)
////////////////////////////////////////////////////////////////
func connect(name string, clientId string) stan.Conn {
	clientId += name
	conn, err := stan.Connect("test-cluster", clientId)
	if err != nil {
		log.Fatalf("Cann't connect : %v, make shure NATS Streaming server is running at %s", err, clientID())
	}
	log.Info("Connection performed for ClientID: ", clientId)
	return conn

}

////////////////////////////////////////////////////////////////
var acb = func(lguid string, err error) {
	if err != nil {
		log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
	}
	// Эта ф-ция будет выполнять действия последней. так как на каждое сообщение будет взят свой uuid
	// тут мы делаем декримент от общего количества "открытых" ГОрутин
	log.Info("Gorutine inside WG of CodeGenerator Done work")
	wg.Done()
}

/////////////////////////////////////////////////////////////////
func publishMsg(msg string, conn stan.Conn) {
	guid, err := conn.PublishAsync("foo", []byte(msg), acb)
	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
	if guid == "" {
		log.Fatal("Expected non-empty guid to be returned.")
	}
}

////////////////////////////////////////////////////////////////
func clientID() string {
	number := 1 + rand.Intn(9999-1)
	clientId := strconv.Itoa(number)
	//log.Info("ClientID is: ", clientId)
	return clientId
}
