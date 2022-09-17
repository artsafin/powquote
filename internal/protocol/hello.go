package protocol

import (
	"bytes"
	"strconv"
)

const separator = "--"

var Hello = []byte("HELLO")

type ChallengeRequest struct{}

type Challenge struct {
	Nonce      uint64
	Complexity int
}

func ChallengeFromBytes(bs []byte) (c Challenge, err error) {
	parts := bytes.SplitN(bs, []byte(separator), 2)
	c.Nonce, err = strconv.ParseUint(string(parts[0]), 10, 64)
	if err != nil {
		return
	}
	complexity, err := strconv.ParseInt(string(parts[1]), 10, 64)
	if err != nil {
		return
	}
	c.Complexity = int(complexity)
	return c, nil
}

func (c Challenge) Bytes() []byte {
	var bs []byte
	bs = append(bs, strconv.FormatUint(c.Nonce, 10)...)
	bs = append(bs, separator...)
	bs = append(bs, strconv.FormatInt(int64(c.Complexity), 10)...)
	return bs
}
