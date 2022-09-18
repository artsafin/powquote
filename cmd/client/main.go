package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"os"
	"strconv"
	"time"

	"powquote/internal/protocol"
	"powquote/internal/puzzle"
)

var ioTimeout = time.Second * 10

func main() {
	verboseVar := os.Getenv("VERBOSE")
	verbose, err := strconv.ParseBool(verboseVar)
	if err != nil {
		verbose = true
	}

	serverAddr := os.Getenv("SERVER")
	if serverAddr == "" {
		log.Fatal("SERVER variable must point to a server")
	}

	clientID, serverID, err := getIPs(serverAddr)
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

	if verbose {
		log.Println("solving challenge from server:", challenge)
	}

	clientNonce := puzzle.GenerateNonceOnce()

	hashData := protocol.HashData{
		ClientID:    clientID,
		NonceServer: challenge.Nonce,
		NonceClient: clientNonce,
	}

	goodhash := puzzle.Solve(&hashData, challenge)
	if verbose {
		log.Printf("found solution: %v", goodhash)
	}

	quoteReq := protocol.QuoteRequest{
		ServerID: serverID,
		HashData: hashData,
	}
	if verbose {
		log.Printf("making quote request: %q", quoteReq.Bytes())
	}
	quote, err := say(serverAddr, quoteReq.Bytes())
	if err != nil {
		log.Fatalf("error sending solution: %v", err)
	}

	if verbose {
		log.Printf("(👉ﾟヮﾟ)👉 %s", quote)
	} else {
		fmt.Printf("%s", quote)
	}
}

func getIPs(serverAddr string) (string, string, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return "", "", err
	}
	defer func(){
		_ = conn.Close()
	}()

	locAddr, err := netip.ParseAddrPort(conn.LocalAddr().String())
	if err != nil {
		return "", "", err
	}
	return locAddr.Addr().String(), conn.RemoteAddr().String(), nil
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
