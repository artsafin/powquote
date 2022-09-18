package puzzle

import (
	"io"
	"strconv"
	"strings"
	"testing"

	"powquote/internal/protocol"

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

func TestReadRequest(t *testing.T) {
	tests := []struct {
		reader io.Reader
		want   any
		err    assert.ErrorAssertionFunc
	}{
		{
			reader: strings.NewReader(string(protocol.Hello)),
			want:   protocol.ChallengeRequest{},
			err:    assert.NoError,
		},
		{
			reader: strings.NewReader(`10.0.0.1:9999--10.1.0.1--111--222--eHl6`),
			want: protocol.QuoteRequest{
				ServerID: "10.0.0.1:9999",
				HashData: protocol.HashData{
					ClientID:    "10.1.0.1",
					NonceServer: 111,
					NonceClient: 222,
					Solution:    []byte("xyz"),
				},
			},
			err: assert.NoError,
		},
		{
			reader: strings.NewReader(`10.0.0.1:9999--10.1.0.1`),
			err: assert.Error,
		},
		{
			reader: strings.NewReader(""),
			err:    ErrorLike(`invalid request`),
		},
		{
			reader: strings.NewReader("\n\n"),
			err:    ErrorLike(`invalid request`),
		},
	}
	for name, tt := range tests {
		t.Run(strconv.Itoa(name), func(t *testing.T) {
			got, err := ReadRequest(tt.reader)
			if tt.err(t, err) && err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
