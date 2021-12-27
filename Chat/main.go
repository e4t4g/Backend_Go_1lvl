package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	nickname := nick()
	who := conn.RemoteAddr().String() + " " + nickname
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)

	for input.Scan() {
		messages <- nickname + ": " + input.Text()
	}
	leaving <- ch
	messages <- nickname + " has left"
	err := conn.Close()
	if err != nil {
		return
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//name request from terminal after connection
func nick() string {
	fmt.Println("Enter your nickname: ")
	nick1, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.Trim(nick1, "\n")

}
