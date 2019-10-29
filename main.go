package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

var (
	serveVia   = flag.String("serve_via", "127.0.0.1", "TCP listener host")
	connectVia = flag.String("connect_via", "127.0.0.1", "ffmpeg connection string")
	port       = flag.Int("port", 12345, "port number to run on")
)

func main() {
	flag.Parse()

	log.Printf("Running -serve_via=%s -connect_via=%s -port=%d", *serveVia, *connectVia, *port)

	serveAddress := fmt.Sprintf("%s:%d", *serveVia, *port)
	log.Printf("Serving via %s", serveAddress)
	lis, err := net.Listen("tcp", serveAddress)
	if err != nil {
		log.Fatalf("cannot open tcp port, %s", err)
	}
	defer lis.Close()
	go func() {
		logListen(lis)
	}()

	log.Printf("wait 5 seconds to make sure TCP is available")

	time.Sleep(time.Second * 5)
	connectToAddress := fmt.Sprintf("%s:%d", *connectVia, *port)
	callFfmpeg(connectToAddress)
	lis.Close()
}

func logListen(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err == nil {
			go ffmpegProgressReader(conn)
		} else {
			log.Printf("error handling TCP %s, ignoring since this happens on ffmpeg close", err)
			return
		}
	}
}

func ffmpegProgressReader(conn net.Conn) {
	defer conn.Close()
	// Read the incoming connection into the buffer.
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		t := scanner.Text()
		log.Printf("ffmpeg progress: %s", t)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("ffmpeg progress error: %s", err)
	}
}

func callFfmpeg(addr string) {
	args := []string{"-y", "-progress", "tcp://" + addr, "-i", "/app/small.ogv", "-strict", "-2", "/app/small.mp4"}
	log.Printf("calling ffmpeg with %+v", args)
	cmd := exec.Command("/app/ffmpeg", args...)
	out, err := cmd.CombinedOutput()
	log.Printf("%s", out)
	log.Printf("\n\n\nffmpeg has closed.")
	if err != nil {
		log.Fatalf("cannot call ffmpeg %s", err)
	}
}
