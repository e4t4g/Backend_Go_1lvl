package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	listener    net.Listener
	Connections chan net.Conn
}

var messagePrint = make(chan string)

func NewServer(address string) Server {
	lister, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic(err)
	}

	connChan := make(chan net.Conn)

	return Server{
		listener:    lister,
		Connections: connChan,
	}
}

func (s Server) Start() {
	log.Printf("server started on %s\n", s.listener.Addr())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		s.Connections <- conn
	}
}

func main() {
	srv := NewServer(":8001")
	go srv.Start()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			log.Println("start graceful")
			wg.Wait()
			log.Println("stop graceful")
			return
		case conn := <-srv.Connections:
			wg.Add(1)
			go handleConn(ctx, conn, &wg)
		}
	}
}

func handleConn(ctx context.Context, c net.Conn, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		c.Close()
	}()

	go msgSender()
	for {
		t := time.NewTicker(time.Second)

		select {
		case <-ctx.Done():
			io.WriteString(c, "Bye!")
			return
		case now := <-t.C:
			io.WriteString(c, now.Format("15:04:05\n\r"))
		case msg := <-messagePrint:
			io.WriteString(c, time.Now().Format("15:04:05")+" "+msg+"\n")

		}
	}
}

func msgSender() {
	fmt.Println("Enter your message:")

	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	textUpgrade := strings.Trim(text, "\n")
	messagePrint <- textUpgrade
}
