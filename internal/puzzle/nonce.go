package puzzle

import (
	"context"
	"crypto/rand"
	"math"
	"math/big"
	"sync/atomic"
	"time"
)

type nonceGenerator struct {
	period time.Duration
	value  atomic.Value
}

func GenerateNonceOnce() uint64 {
	n := &nonceGenerator{}
	n.tick()
	return n.Current()
}

func NewNonceGenerator(period time.Duration) *nonceGenerator {
	return &nonceGenerator{
		period: period,
	}
}

func (n *nonceGenerator) tick() {
	v, err := rand.Int(rand.Reader, new(big.Int).SetUint64(math.MaxUint64))
	if err != nil {
		panic(err)
	}
	n.value.Store(v.Uint64())
}
func (n *nonceGenerator) Current() uint64 {
	if v := n.value.Load().(uint64); v == 0 {
		n.tick()
	}
	return n.value.Load().(uint64)
}

func (n *nonceGenerator) Start(ctx context.Context) {
	ticker := time.NewTicker(n.period)

	n.tick()

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-ticker.C:
			if !ok {
				return
			}
			n.tick()
		}
	}
}
