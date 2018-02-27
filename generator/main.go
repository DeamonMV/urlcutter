package main

import (
	"log"
	"sync"
	"math/rand"
	"github.com/nats-io/go-nats-streaming"
	"time"
	json2 "encoding/json"
	"strconv"
	"fmt"
)

const azbuka string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const azbukaLength = len(azbuka)

var codeLenght int = 4

//объявляем перменную  на которую навешиваем WaitGoup, которая позволит завершить все начатые ГОрутины
var wg sync.WaitGroup
//////////////////////////
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
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
		Code string
	}

	for i := 0; i < powerOf; i++ {

		for elOfAzbuka := 0; elOfAzbuka < codeLenght; elOfAzbuka++ {
			result += string(azbuka[indexes[elOfAzbuka]])
		}

		data := &Message{
			Index: strconv.Itoa(i),
			Code: result}

		json, err:= json2.Marshal(data)
		if err != nil{
			fmt.Printf("Error marshal struc: &v", err)
		}

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
func connect() stan.Conn {

	conn, err := stan.Connect("test-cluster", "stan2")
	if err != nil {
		log.Fatalf("Cann't connect : %v, make shure NATS Streaming server is running at %s", err)
	}
	return conn

}

////////////////////////////////////////////////////////////////
var acb = func(lguid string, err error) {
	if err != nil {
		log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
	}
	// Эта ф-ция будет выполнять действия последней. так как на каждое сообщение будет взят свой uuid
	// тут мы делаем декримент от общего количества "открытых" ГОрутин
	wg.Done()
	//else {
	//	log.Printf("Received ack for msg id %s\n", lguid)
	//}

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

func main() {




	c_codes := make(chan string)
	powerof := powerOf(codeLenght)
	println(powerof)

	conn := connect()
	go codeGenerator(codeLenght, powerof, c_codes)

	for msg := range c_codes {
		// увеличиваем количество ГОрутин, для которых надо будет дожлаться закрытие
		wg.Add(1)
		go publishMsg(msg, conn)
		time.Sleep(300 * time.Millisecond)
	}
	// ждем когда отработают все ГОрутины
	wg.Wait()
	conn.Close()
}
