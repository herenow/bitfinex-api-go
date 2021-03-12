package utils

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// v2 types

type NonceGenerator interface {
	GetNonce() string
}

type EpochNonceGenerator struct {
	nonce uint64
	lock  *sync.Mutex
}

// GetNonce is a naive nonce producer that takes the current Unix nano epoch
// and counts upwards.
// This is a naive approach because the nonce bound to the currently used API
// key and as such needs to be synchronised with other instances using the same
// key in order to avoid race conditions.
func (u *EpochNonceGenerator) GetNonce() string {
	u.lock.Lock()
	defer u.lock.Unlock()

	n := getNonceFromTime()
	if n <= u.nonce {
		n = u.nonce + 1
	}

	u.nonce = n
	return strconv.FormatUint(n, 10)
}

func NewEpochNonceGenerator() *EpochNonceGenerator {
	return &EpochNonceGenerator{
		nonce: getNonceFromTime(),
		lock:  &sync.Mutex{},
	}
}

func getNonceFromTime() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}

// v1 support

var nonce uint64

func init() {
	nonce = uint64(time.Now().UnixNano()) * 1000000
}

// GetNonce is a naive nonce producer that takes the current Unix nano epoch
// and counts upwards.
// This is a naive approach because the nonce bound to the currently used API
// key and as such needs to be synchronised with other instances using the same
// key in order to avoid race conditions.
func GetNonce() string {
	return strconv.FormatUint(atomic.AddUint64(&nonce, 1), 10)
}
