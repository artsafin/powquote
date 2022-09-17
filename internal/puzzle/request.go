package puzzle

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"

	"powquote/internal/protocol"
)

func ProtectionEnabled() bool {
	protected := os.Getenv("PROTECTED")
	if protected == "" {
		return true
	}
	value, err := strconv.ParseBool(protected)
	if err != nil {
		return true
	}
	return value
}

func ReadRequest(r io.Reader) (any, error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		token := sc.Bytes()
		if len(token) == 0 {
			continue
		}

		if bytes.EqualFold(token, protocol.Hello) {
			return protocol.ChallengeRequest{}, nil
		}

		quoteRequest, err := protocol.ParseQuoteRequest(token)
		if err != nil {
			return nil, err
		}
		return quoteRequest, nil
	}
	if sc.Err() != nil {
		return nil, sc.Err()
	}

	return nil, errors.New("invalid request")
}
