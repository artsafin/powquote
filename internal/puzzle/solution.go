package puzzle

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"sync"

	"powquote/internal/protocol"
)

type solutionAttempt struct {
	clientID    string
	nonceServer uint64
	nonceClient uint64
}

var solutionAttempts = make(map[solutionAttempt]struct{})
var attemptsMutex sync.Mutex

func stripPort(addr net.Addr) string {
	addrPort, err := netip.ParseAddrPort(addr.String())
	if err != nil {
		return addr.String()
	}

	return addrPort.Addr().String()
}

func SolutionValid(challenge protocol.Challenge, serverAddr net.Addr, clientAddr net.Addr, req protocol.QuoteRequest) error {
	if req.ServerID != serverAddr.String() {
		return fmt.Errorf("server addr: %v != %v", req.ServerID, serverAddr.String())
	}
	if req.ClientID != stripPort(clientAddr) {
		return fmt.Errorf("client addr: %v != %v", req.ClientID, clientAddr.String())
	}
	if req.NonceServer != challenge.Nonce {
		return fmt.Errorf("server nonce: %v != %v", req.NonceServer, challenge.Nonce)
	}

	attemptsMutex.Lock()
	defer attemptsMutex.Unlock()

	attempt := solutionAttempt{
		clientID:    req.ClientID,
		nonceServer: req.NonceServer,
		nonceClient: req.NonceClient,
	}
	if _, ok := solutionAttempts[attempt]; ok {
		return fmt.Errorf("attempt exist: %v", attempt)
	}

	calchash := Hash(&req.HashData)

	if !HashMatchesChallenge(calchash, challenge) {
		return fmt.Errorf("invalid hash solution %q: hash %v, complexity=%v", req.Bytes(), calchash, challenge.Complexity)
	}

	solutionAttempts[attempt] = struct{}{}

	return nil
}

func Hash(req *protocol.HashData) string {
	var buf bytes.Buffer
	buf.WriteString(req.ClientID)
	buf.WriteByte(';')
	buf.WriteString(strconv.FormatUint(req.NonceServer, 10))
	buf.WriteByte(';')
	buf.WriteString(strconv.FormatUint(req.NonceClient, 10))
	buf.WriteByte(';')
	buf.Write(req.Solution)

	sum := sha1.Sum(buf.Bytes())
	return hex.EncodeToString(sum[:])
}

func HashMatchesChallenge(hash string, challenge protocol.Challenge) bool {
	runes := []rune(hash)

	if len(runes) < challenge.Complexity {
		return false
	}
	for i := 0; i < challenge.Complexity; i++ {
		if runes[i] != '0' {
			return false
		}
	}

	return true
}
