package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"github.com/paradoxxl/gomastermind/msg"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/echo"
var guess = make(chan interface{})

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan interface{})

	guess = make(chan interface{})

	go receive(ws)

	sendNewGame(ws)
	<-guess
	sendGuess1(ws)
	<-guess
	sendGuess2(ws)
	<-c

}

func sendNewGame(ws *websocket.Conn) {
	data := msg.NewGameMsg{NbrColors: 2, MaxTries: 100, CodeLength: 2}
	sendData := EncodeNewGameMsg(&data)
	send(&sendData, ws)
}
func sendGuess1(ws *websocket.Conn) {
	g := make(msg.Code, 2)
	g[0] = 0
	g[1] = 1
	data := msg.GuessMsg{g}
	sendData := EncodeGuessMsg(&data)
	send(&sendData, ws)
}
func sendGuess2(ws *websocket.Conn) {
	g := make(msg.Code, 3)
	g[0] = 1
	g[1] = 0
	g[2] = 0
	data := msg.GuessMsg{g}
	sendData := EncodeGuessMsg(&data)
	send(&sendData, ws)
}

func send(msg *string, ws *websocket.Conn) {
	if err := websocket.Message.Send(ws, *msg); err != nil {
		log.Printf("cannot send: %v", err)
	}
	fmt.Printf("Send: %s\n", *msg)
}

func receive(ws *websocket.Conn) {
	for {
		var message2 string
		if err := websocket.Message.Receive(ws, &message2); err != nil {
			log.Printf("receive error %v", err)
		}
		fmt.Printf("Receive: %s\n", message2)
		guess <- 1
	}
}

func EncodeNewGameMsg(data *msg.NewGameMsg) string {
	header := []byte(fmt.Sprintf("%v|", msg.NewGameType))
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	//binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return string(p)
}
func EncodeGuessMsg(data *msg.GuessMsg) string {
	header := []byte(fmt.Sprintf("%v|", msg.GuessType))
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	//binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return string(p)
}
