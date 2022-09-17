package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"powquote/internal/protocol"
	"powquote/internal/puzzle"
	"powquote/internal/quotes"
)

var nonces = puzzle.NewNonceGenerator(time.Minute * 5)

var complexity = 6

var ioTimeout = time.Second * 30

func main() {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		panic("invalid LISTEN variable")
	}

	ln, err := net.Listen("tcp", listen)
	if err != nil {
		panic(err)
	}

	rootctx := context.Background()

	go nonces.Start(rootctx)

	log.Printf("begin listening on %v; DoS protected = %v, complexity = %v", ln.Addr(), puzzle.ProtectionEnabled(), complexity)

	for {
		conn, err := ln.Accept()
		if err := conn.SetDeadline(time.Now().Add(ioTimeout)); err != nil {
			log.Printf("error setting deadline: %v", err)
		}
		if err != nil {
			log.Printf("error accepting connection: %v", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		addr := conn.RemoteAddr()
		if err := conn.Close(); err != nil {
			log.Printf("(%v) error closing connection: %v", addr, err)
			return
		}
		log.Printf("(%v) connection closed by server", addr)
	}()

	log.Printf("(%v) connected", conn.RemoteAddr())

	if !puzzle.ProtectionEnabled() {
		writeResponse(conn, []byte(quotes.Next()))
		return
	}

	req, err := puzzle.ReadRequest(conn)
	if err != nil {
		log.Printf("(%v) error processing request: %v", conn.RemoteAddr(), err)
		writeResponse(conn, []byte("send hello request to begin client puzzle"))
		return
	}

	challenge := protocol.Challenge{
		Nonce:      nonces.Current(),
		Complexity: complexity,
	}

	switch req := req.(type) {
	case protocol.ChallengeRequest:
		log.Printf("(%v) challenge request", conn.RemoteAddr())
		writeResponse(conn, challenge.Bytes())
	case protocol.QuoteRequest:
		log.Printf("(%v) quote request", conn.RemoteAddr())
		if err := puzzle.SolutionValid(challenge, conn.LocalAddr(), req); err == nil {
			log.Printf("(%v) solution correct %v", conn.RemoteAddr(), puzzle.Hash(&req.HashData))
			writeResponse(conn, []byte(quotes.Next()))
		} else {
			log.Printf("(%v) invalid solution: %v", conn.RemoteAddr(), err)
			writeResponse(conn, []byte("invalid solution"))
		}
	}
}

func writeResponse(conn net.Conn, bs []byte) {
	log.Printf("(%v) writing: %s", conn.RemoteAddr(), bs)
	if _, err := conn.Write(bs); err != nil {
		log.Printf("error writing response: %v", err)
	}
}
