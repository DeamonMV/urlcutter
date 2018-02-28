package main

import (
	json2 "encoding/json"
	"fmt"
	"github.com/nats-io/go-nats-streaming"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"urlcutter/generator/sub/checker"

	"os"
	"os/signal"
)

const azbuka string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const azbukaLength = len(azbuka)
//объявляем перменную  на которую навешиваем WaitGoup, которая позволит завершить все начатые ГОрутины для создания кодов
var wg sync.WaitGroup
//var codeLenght int = 4

//////////////////////////
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

/////////////////////////////////////////////////////////////////
func powerOf(length int) int {
	if length == 0 {
		return length
	}
	var l int = length
	var result int = azbukaLength

	for i := 0; i < l-1; i++ {
		result = result * azbukaLength
	}
	return result
}

////////////////////////////////////////////////////////////////
func codeGenerator(codeLenght int, powerOf int, c_codes chan<- string) {

	var indexes = [azbukaLength]int{}
	var result string

	type Message struct {
		Index string
		Code  string
	}

	for i := 0; i < powerOf; i++ {

		for elOfAzbuka := 0; elOfAzbuka < codeLenght; elOfAzbuka++ {
			result += string(azbuka[indexes[elOfAzbuka]])
		}

		data := &Message{
			Index: strconv.Itoa(i),
			Code:  result}

		json, err := json2.Marshal(data)
		if err != nil {
			fmt.Printf("Error marshal struc: &v", err)
		}
		log.Info(string(json))
		c_codes <- string(json)

		result = ""

		for index := 0; index < azbukaLength; index++ {
			indexes[index]++
			if indexes[index] < azbukaLength {
				break
			}
			indexes[index] = 0
		}
	}
	//закрываем канал что бы не расходовать ресурсы
	close(c_codes)
}

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
/////////////////////////////////////////////////////////////////
func main() {

	c_startgen := make(chan int)
	c_shutdown := make(chan struct{})
	//c_codeGenIsWork := make(chan bool)
	var wg_gors sync.WaitGroup

	log.Info("Main func Started")

	conn := connect(clientID(), "SUB")
	defer conn.Close()
	defer log.Info("SUB: Connection Closed")
	log.Info("SUB: Gorutine started")

	msgConn := func(m *stan.Msg) {
		m.Ack()
		log.Info("Recive data from NATS: ", m)

		mapa := checker.Unmarshal(m.Data)
		codeint, err := checker.Checker(mapa)
		if err != "" {
			log.Warn(err)
		}
		log.Info("####### CODELEN #########: ", codeint)
		c_startgen <- codeint
	}

	aw, _ := time.ParseDuration("60s")
	_, err := conn.QueueSubscribe("dbproc", "generator", msgConn, stan.DeliverAllAvailable(), stan.SetManualAckMode(), stan.AckWait(aw), stan.DurableName("dur1"))

	if err != nil {
		fmt.Printf("Error %v\n", err)
	}

	log.Info("SUB: works")
////////////////////////////////////////////////////////////////////////////////////////
	go func()  {

		conn := connect(clientID(), "CodeGen")
		wg_gors.Add(1)
		defer wg_gors.Done()
		defer conn.Close()
		defer log.Info("GodeGen: Connection Closed")
		log.Info("CodeGen: Gorutine with CodeGenerator started")
		for {
			select {
			case startgen :=  <-c_startgen:
				{
					c_codes := make(chan string)
					//c_codeGenIsWork <- true
					//defer close(c_codeGenIsWork)
					log.Info("Work")
					log.Info("Codelen: ", startgen)
					//if startgen != 0 {
					powerof := powerOf(startgen)
					i := 0
					go codeGenerator(startgen, powerof, c_codes)

					for msg := range c_codes {
						i++
						log.Info("index of gorutines: ", i)
						// увеличиваем количество ГОрутин, для которых надо будет дожлаться закрытие
						wg.Add(1)
						go publishMsg(msg, conn)
						log.Info("Message was published: ", msg)
						time.Sleep(300 * time.Millisecond)
					}
					log.Info("CodeGen: done code generation")
				}

				case  <- c_shutdown:
				{
					log.Info("CodeGen: Recive Shutdown signal")
					return
				}
			}

		}
	}()
////////////////////////////////////////////////////////////////////////////////////////////
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {

				<- signalChan
					log.Info("\nReceived an interrupt and closing connection...\n\n")

					//conn.Close()
					close(c_shutdown)
					wg_gors.Wait()
					cleanupDone <- true



	}()
	<- cleanupDone

}
