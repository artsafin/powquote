package protocol

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
)

const (
	maxSolutionLength = 1024 * 64
)

const (
	quoteRequestFieldServerID = iota
	quoteRequestFieldClientID
	quoteRequestFieldServerNonce
	quoteRequestFieldClientNonce
	quoteRequestFieldSolution
	quoteRequestEOF
)

type HashData struct {
	ClientID    string
	NonceServer uint64
	NonceClient uint64
	Solution    []byte
}

type QuoteRequest struct {
	ServerID string
	HashData
}

// ParseQuoteRequest parses solution consisting of S, C, Ns, Nc, X separated by colon
func ParseQuoteRequest(solution []byte) (qr QuoteRequest, err error) {
	if solutionLen := len(solution); solutionLen > maxSolutionLength {
		return qr, fmt.Errorf("solution is too long: %v", solutionLen)
	}
	fields := bytes.SplitN(solution, []byte(separator), quoteRequestEOF)
	if len(fields) != quoteRequestEOF {
		return qr, fmt.Errorf("number of fields in solution is invalid: %v, expected %v", len(fields), quoteRequestEOF)
	}
	for field := 0; field < len(fields); field++ {
		switch field {
		case quoteRequestFieldServerID:
			qr.ServerID = string(fields[field])
		case quoteRequestFieldClientID:
			qr.ClientID = string(fields[field])
		case quoteRequestFieldServerNonce:
			qr.NonceServer, err = strconv.ParseUint(string(fields[field]), 10, 64)
		case quoteRequestFieldClientNonce:
			qr.NonceClient, err = strconv.ParseUint(string(fields[field]), 10, 64)
		case quoteRequestFieldSolution:
			qr.Solution, err = base64.StdEncoding.DecodeString(string(fields[field]))
		}
		if err != nil {
			return qr, fmt.Errorf("field: %v: %v", field, err)
		}
	}
	return
}

func (r *QuoteRequest) Bytes() []byte {
	var buf bytes.Buffer

	buf.WriteString(r.ServerID)
	buf.WriteString(separator)
	buf.WriteString(r.ClientID)
	buf.WriteString(separator)
	buf.WriteString(strconv.FormatUint(r.NonceServer, 10))
	buf.WriteString(separator)
	buf.WriteString(strconv.FormatUint(r.NonceClient, 10))
	buf.WriteString(separator)
	buf.WriteString(base64.StdEncoding.EncodeToString(r.Solution))

	return buf.Bytes()
}
