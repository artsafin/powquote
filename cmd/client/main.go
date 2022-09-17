package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"powquote/internal/protocol"
	"powquote/internal/puzzle"
)

var ioTimeout = time.Second * 10

func main() {
	serverAddr := os.Getenv("SERVER")
	if serverAddr == "" {
		panic("SERVER variable must point to a server")
	}

	clientID, err := getLocalIP(serverAddr)
	if err != nil {
		log.Fatalf("unable to detect client id: %v", err)
	}

	challengeBs, err := say(serverAddr, protocol.Hello)
	if err != nil {
		log.Fatalf("error saying to server: %v", err)
	}

	challenge, err := protocol.ChallengeFromBytes(challengeBs)
	if err != nil {
		log.Fatalf("error parsing challenge: %v", err)
	}

	log.Println("solving challenge from server:", challenge)

	clientNonce := puzzle.GenerateNonceOnce()

	hashData := protocol.HashData{
		ClientID:    clientID,
		NonceServer: challenge.Nonce,
		NonceClient: clientNonce,
	}

	goodhash := solve(&hashData, challenge)
	log.Printf("found solution: %v", goodhash)

	quoteReq := protocol.QuoteRequest{
		ServerID: serverAddr,
		HashData: hashData,
	}
	log.Printf("making quote request: %q", quoteReq.Bytes())
	quote, err := say(serverAddr, quoteReq.Bytes())
	if err != nil {
		log.Fatalf("error sending solution: %v", err)
	}

	log.Printf("(👉ﾟヮﾟ)👉 %s", quote)
}

func solve(hashData *protocol.HashData, challenge protocol.Challenge) string {
	randBs := make([]byte, 64)
	var lasthash string

	for {
		if _, err := rand.Read(randBs); err != nil {
			log.Fatalf("error generating solution: %v", err)
		}
		hashData.Solution = randBs
		lasthash = puzzle.Hash(hashData)
		if puzzle.HashMatchesChallenge(lasthash, challenge) {
			break
		}
	}
	return lasthash
}

func getLocalIP(serverAddr string) (string, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return conn.LocalAddr().String(), nil
}

func say(serverAddr string, what []byte) ([]byte, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(ioTimeout)); err != nil {
		return nil, err
	}

	what = append(what, '\n')
	if _, err := conn.Write(what); err != nil {
		return nil, err
	}

	bs, err := io.ReadAll(conn)

	if err != nil {
		return nil, err
	}

	return bs, nil
}
