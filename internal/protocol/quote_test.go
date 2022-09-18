package protocol

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteRequest_Bytes(t *testing.T) {
	tests := []struct {
		req  QuoteRequest
		want []byte
	}{
		{
			req: QuoteRequest{
				ServerID: "10.0.0.1:9999",
				HashData: HashData{
					ClientID:    "10.1.0.1",
					NonceServer: 111,
					NonceClient: 222,
					Solution:    []byte("xyz"),
				},
			},
			want: []byte("10.0.0.1:9999--10.1.0.1--111--222--eHl6"),
		},
	}
	for name, tt := range tests {
		t.Run(strconv.Itoa(name), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.req.Bytes())
		})
	}
}

func TestParseQuoteRequest(t *testing.T) {
	tests := []struct {
		solution []byte
		wantQr   QuoteRequest
		err      assert.ErrorAssertionFunc
	}{
		{
			solution: []byte("10.0.0.1:9999--10.1.0.1--111--222--eHl6"),
			wantQr: QuoteRequest{
				ServerID: "10.0.0.1:9999",
				HashData: HashData{
					ClientID:    "10.1.0.1",
					NonceServer: 111,
					NonceClient: 222,
					Solution:    []byte("xyz"),
				},
			},
			err: assert.NoError,
		},
		{
			solution: bytes.Repeat([]byte{'a'}, maxSolutionLength + 1),
			err: ErrorLike(`solution is too long: 65537`),
		},
		{
			solution: []byte("10.0.0.1:9999--10.1.0.1--111--222--eHl6--foo"),
			err: ErrorLike(`field: 4: illegal base64 data`),
		},
		{
			solution: []byte("10.0.0.1:9999--10.1.0.1"),
			err: ErrorLike(`number of fields in solution is invalid: 2, expected 5`),
		},
		{
			solution: []byte("10.0.0.1:9999--10.1.0.1--18446744073709551616--222--eHl6--foo"),
			err: ErrorLike(`strconv.ParseUint: parsing "18446744073709551616": value out of range`),
		},
	}
	for name, tt := range tests {
		t.Run(strconv.Itoa(name), func(t *testing.T) {
			gotQr, err := ParseQuoteRequest(tt.solution)
			if tt.err(t, err) && err == nil {
				assert.Equal(t, tt.wantQr, gotQr)
			}
		})
	}
}
