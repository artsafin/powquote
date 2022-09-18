package puzzle

import (
	"net"
	"strconv"
	"testing"

	"powquote/internal/protocol"

	"github.com/stretchr/testify/assert"
)

func TestHashMatchesChallenge(t *testing.T) {
	tests := []struct {
		hash      string
		challenge protocol.Challenge
		want      bool
	}{
		{
			hash: "0000004c751b61da69a82eb3b5067a2494b3cd2e",
			challenge: protocol.Challenge{
				Complexity: 6,
			},
			want: true,
		},
		{
			hash: "0000014c751b61da69a82eb3b5067a2494b3cd2e",
			challenge: protocol.Challenge{
				Complexity: 6,
			},
			want: false,
		},
		{
			hash: "0asdf14c751b61da69a82eb3b5067a2494b3cd2e",
			challenge: protocol.Challenge{
				Complexity: 1,
			},
			want: true,
		},
		{
			hash: "wasdf14c751b61da69a82eb3b5067a2494b3cd2e",
			challenge: protocol.Challenge{
				Complexity: 0,
			},
			want: true,
		},
		{
			hash: "0000000000000000000000000000000000000000",
			challenge: protocol.Challenge{
				Complexity: -1,
			},
			want: true,
		},
		{
			hash: "00000000",
			challenge: protocol.Challenge{
				Complexity: 9,
			},
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(strconv.Itoa(name), func(t *testing.T) {
			got := HashMatchesChallenge(tt.hash, tt.challenge)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHash(t *testing.T) {
	tests := []struct {
		req  protocol.HashData
		want string
	}{
		{
			req: protocol.HashData{
				ClientID:    "10.1.0.1",
				NonceServer: 111,
				NonceClient: 222,
				Solution:    []byte("xyz"),
			},
			want: "a4e65a96591129e7bf2ac20b9b2d6f3809598c04",
		},
	}

	for name, tt := range tests {
		t.Run(strconv.Itoa(name), func(t *testing.T) {
			assert.Equal(t, tt.want, Hash(&tt.req))
		})
	}
}

func TestSolutionValid(t *testing.T) {
	type args struct {
		challenge  protocol.Challenge
		serverAddr net.Addr
		clientAddr net.Addr
		req        protocol.QuoteRequest
	}
	var validSolution = args{
		challenge: protocol.Challenge{
			Nonce:      4874918909949807476,
			Complexity: 6,
		},
		serverAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 2}, Port: 9999},
		clientAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 3}, Port: 7878},
		req: protocol.QuoteRequest{
			ServerID: "172.18.0.2:9999",
			HashData: protocol.HashData{
				ClientID:    "172.18.0.3",
				NonceServer: 4874918909949807476,
				NonceClient: 2190648595078496803,
				Solution:    []byte{0xe8, 0xc1, 0xcd, 0xfb, 0xa2, 0xbe, 0xed, 0xe4, 0x8b, 0xf6, 0x86, 0x74, 0x7b, 0xb4, 0x8a, 0x84, 0x9c, 0xb3, 0x9b, 0x35, 0x93, 0x6a, 0x6c, 0xb, 0x39, 0xd3, 0xc2, 0x6, 0xc7, 0x21, 0x4b, 0xe7, 0x21, 0x0, 0xa8, 0x76, 0xbf, 0xde, 0x8, 0xd9, 0xe0, 0xb9, 0xc9, 0x79, 0x2e, 0x63, 0x82, 0x36, 0xb2, 0xce, 0x21, 0x8a, 0xc1, 0x54, 0x54, 0x46, 0x7e, 0x33, 0x78, 0xba, 0x83, 0x9e, 0xf, 0x26, 0x31, 0x59, 0x1e, 0xd8, 0x46, 0x92, 0xd0, 0xbd, 0xc8, 0x8f, 0xa4, 0x45, 0x9a, 0x9f, 0xb7, 0xcd, 0x57, 0x59, 0xa9, 0xc1, 0x29, 0xb4, 0xc4, 0x5d, 0x32, 0x1a, 0x71, 0x53, 0x1c, 0x74, 0x2e, 0x49, 0xcd, 0x5b, 0x53, 0x7e, 0x58, 0x60, 0xc, 0x5c, 0x9, 0x82, 0xda, 0xd3, 0x28, 0x7c, 0xbf, 0xf8, 0xe0, 0x83, 0xf0, 0xb7, 0x6f, 0x57, 0x78, 0xf8, 0x8c, 0xff, 0x6, 0xbe, 0xd7, 0xda, 0x12, 0xd9},
			},
		},
	}

	tests := []struct {
		name string
		args args
		err  assert.ErrorAssertionFunc
	}{
		{
			name: "solution is always valid if complexity 0",
			args: args{
				challenge: protocol.Challenge{
					Nonce:      111,
					Complexity: 0,
				},
				serverAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 2}, Port: 9999},
				clientAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 3}, Port: 1234},
				req: protocol.QuoteRequest{
					ServerID: "172.18.0.2:9999",
					HashData: protocol.HashData{
						ClientID:    "172.18.0.3",
						NonceServer: 111,
						NonceClient: 222,
						Solution:    []byte("xyz"),
					},
				},
			},
			err: assert.NoError,
		},
		{
			name: "valid solution",
			args: validSolution,
			err: assert.NoError,
		},
		{
			name: "repeated valid solution is not accepted",
			args: validSolution,
			err: ErrorLike(`attempt exist: {172.18.0.3 4874918909949807476 2190648595078496803}`),
		},
		{
			name: "invalid solution",
			args: args{
				challenge: protocol.Challenge{
					Nonce:      1111,
					Complexity: 1,
				},
				serverAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 2}, Port: 9999},
				clientAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 3}, Port: 1234},
				req: protocol.QuoteRequest{
					ServerID: "172.18.0.2:9999",
					HashData: protocol.HashData{
						ClientID:    "172.18.0.3",
						NonceServer: 1111,
						NonceClient: 2222,
						Solution:    []byte("xyz"),
					},
				},
			},
			err: ErrorLike(`invalid hash solution`),
		},
		{
			name: "server nonce in challenge must match quote request",
			args: args{
				challenge: protocol.Challenge{
					Nonce:      111,
					Complexity: 0,
				},
				serverAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 2}, Port: 9999},
				clientAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 3}, Port: 1234},
				req: protocol.QuoteRequest{
					ServerID: "172.18.0.2:9999",
					HashData: protocol.HashData{
						ClientID:    "172.18.0.3",
						NonceServer: 111222,
						NonceClient: 222,
						Solution:    []byte("xyz"),
					},
				},
			},
			err: ErrorLike(`server nonce: 111222 != 111`),
		},
		{
			name: "server IP in challenge must match quote request",
			args: args{
				challenge: protocol.Challenge{
					Nonce:      111,
					Complexity: 0,
				},
				serverAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 222}, Port: 9999},
				clientAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 3}, Port: 1234},
				req: protocol.QuoteRequest{
					ServerID: "172.18.0.2:9999",
					HashData: protocol.HashData{
						ClientID:    "172.18.0.3",
						NonceServer: 111,
						NonceClient: 222,
						Solution:    []byte("xyz"),
					},
				},
			},
			err: ErrorLike(`server addr: 172.18.0.2:9999 != 172.18.0.222:9999`),
		},
		{
			name: "client IP in challenge must match quote request",
			args: args{
				challenge: protocol.Challenge{
					Nonce:      111,
					Complexity: 0,
				},
				serverAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 2}, Port: 9999},
				clientAddr: &net.TCPAddr{IP: []byte{172, 18, 0, 30}, Port: 1234},
				req: protocol.QuoteRequest{
					ServerID: "172.18.0.2:9999",
					HashData: protocol.HashData{
						ClientID:    "172.18.0.3",
						NonceServer: 111,
						NonceClient: 222,
						Solution:    []byte("xyz"),
					},
				},
			},
			err: ErrorLike(`client addr: 172.18.0.3 != 172.18.0.30`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SolutionValid(tt.args.challenge, tt.args.serverAddr, tt.args.clientAddr, tt.args.req)
			tt.err(t, err)
		})
	}
}
