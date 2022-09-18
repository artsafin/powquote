package protocol

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ErrorLike(contains string) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
		if err == nil {
			if h, ok := t.(interface{ Helper() }); ok {
				h.Helper()
			}
			return assert.Fail(t, "An error is expected but got nil.", msgAndArgs...)
		}

		return assert.ErrorContains(t, err, contains)
	}
}

func TestChallengeFromBytes(t *testing.T) {
	tests := []struct {
		challenge []byte
		want      Challenge
		err       assert.ErrorAssertionFunc
	}{
		{
			challenge: []byte("111--222"),
			want: Challenge{
				Nonce:      111,
				Complexity: 222,
			},
			err: assert.NoError,
		},
		{
			challenge: []byte("111--222--333"),
			err: ErrorLike(`strconv.ParseInt: parsing "222--333": invalid syntax`),
		},
		{
			challenge: []byte("aaaa--222"),
			err: ErrorLike(`strconv.ParseUint: parsing "aaaa": invalid syntax`),
		},
		{
			challenge: []byte("111--bbbb"),
			err: ErrorLike(`strconv.ParseInt: parsing "bbbb": invalid syntax`),
		},
		{
			challenge: []byte("18446744073709551616--222"),
			err: ErrorLike(`strconv.ParseUint: parsing "18446744073709551616": value out of range`),
		},
		{
			challenge: []byte("1--9223372036854775808"),
			err: ErrorLike(`strconv.ParseInt: parsing "9223372036854775808": value out of range`),
		},
	}
	for name, tt := range tests {
		t.Run(strconv.Itoa(name), func(t *testing.T) {
			gotC, err := ChallengeFromBytes(tt.challenge)
			if tt.err(t, err) && err == nil {
				assert.Equal(t, tt.want, gotC)
			}
		})
	}
}

func TestChallenge_Bytes(t *testing.T) {
	tests := []struct {
		ch   Challenge
		want []byte
	}{
		{
			ch: Challenge{
				Nonce:      111,
				Complexity: 222,
			},
			want: []byte("111--222"),
		},
		{
			ch: Challenge{
				Nonce:      math.MaxUint64,
				Complexity: math.MaxInt,
			},
			want: []byte("18446744073709551615--9223372036854775807"),
		},
	}
	for name, tt := range tests {
		t.Run(strconv.Itoa(name), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ch.Bytes())
		})
	}
}
