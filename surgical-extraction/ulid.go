package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

func NewID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).String()
}

var globalRand = rand.New(&lockedSource{
	src: rand.NewSource(time.Now().UnixNano()),
})

type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
