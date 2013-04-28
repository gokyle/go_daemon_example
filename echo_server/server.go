package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var msgchan chan string
var errchan chan error
var shutdown chan int

func msglogger() {
	for {
                msg, ok := <-msgchan
                if !ok {
                        break
                }
                fmt.Println("client message:", msg)
	}
        fmt.Println("shutting down message logger")
        shutdown <- 1
}

func errlogger() {
        for {
                err, ok := <-errchan
                if !ok {
                        break
                }
                fmt.Println("[!]", err.Error())
        }
        fmt.Println("shutting down error logger")
        shutdown <- 1
}

func echo(conn net.Conn) {
	defer conn.Close()

	msg, err := ioutil.ReadAll(conn)
	if err != nil {
		errchan <- err
		return
	}
	msgchan <- string(msg)

        _, err = conn.Write(msg)
        if err != nil {
                errchan <- err
                return
        }
}

func listener() {
	srv, err := net.Listen("tcp", ":4141")
	if err != nil {
		fmt.Println("[!] failed to set up server:", err.Error())
		os.Exit(1)
	}

        fmt.Println("listening on :4141")
	for {
		conn, err := srv.Accept()
		if err != nil {
			errchan <- err
		}

		go echo(conn)
	}

}

func main() {
	errchan = make(chan error, 16)
	msgchan = make(chan string, 16)
        shutdown = make(chan int, 2)

        go errlogger()
        go msglogger()
        go listener()

	sigc := make(chan os.Signal, 1)
        signal.Notify(sigc, os.Kill, os.Interrupt, syscall.SIGTERM)
	<-sigc
        fmt.Println("shutting down...")
        close(errchan)
        close(msgchan)

        // wait for shutdown signal from the two loggers
        var closed = 0
        for {
                <-shutdown
                closed++
                if closed == 2 {
                        break
                }
        }
        fmt.Println("shutdown complete.")
}
