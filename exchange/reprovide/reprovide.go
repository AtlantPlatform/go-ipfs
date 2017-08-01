package reprovide

import (
	"context"
	"fmt"
	"time"

	backoff "gx/ipfs/QmPJUtEJsm5YLUWhF6imvyCH8KZXRJa9Wup7FDMwTy5Ufz/backoff"
	routing "gx/ipfs/QmPjTrrSfE6TzLv6ya6VWhGcCgPrUAdcgrDcQyRDX2VyW1/go-libp2p-routing"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	cid "gx/ipfs/QmTprEaAA2A9bst5XH7exuyi5KzNMK3SEDNN8rBDnKWcUS/go-cid"
)

var log = logging.Logger("reprovider")

type KeyChanFunc func(context.Context) (<-chan *cid.Cid, error)

type Reprovider struct {
	ctx     context.Context
	trigger chan context.CancelFunc

	// The routing system to provide values through
	rsys routing.ContentRouting

	keyProvider KeyChanFunc
}

func NewReprovider(ctx context.Context, rsys routing.ContentRouting, keyProvider KeyChanFunc) *Reprovider {
	return &Reprovider{
		ctx:     ctx,
		trigger: make(chan context.CancelFunc),

		rsys: rsys,
		keyProvider: keyProvider,
	}
}

func (rp *Reprovider) ProvideEvery(tick time.Duration) {
	// dont reprovide immediately.
	// may have just started the daemon and shutting it down immediately.
	// probability( up another minute | uptime ) increases with uptime.
	after := time.After(time.Minute)
	var done context.CancelFunc
	for {
		select {
		case <-rp.ctx.Done():
			return
		case done = <-rp.trigger:
		case <-after:
		}

		err := rp.Reprovide(rp.ctx)
		if err != nil {
			log.Debug(err)
		}

		if done != nil {
			done()
		}
		after = time.After(tick)
	}
}

func (rp *Reprovider) Reprovide(ctx context.Context) error {
	keychan, err := rp.keyProvider(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get key chan: %s", err)
	}
	for c := range keychan {
		op := func() error {
			err := rp.rsys.Provide(ctx, c, true)
			if err != nil {
				log.Debugf("Failed to provide key: %s", err)
			}
			return err
		}

		// TODO: this backoff library does not respect our context, we should
		// eventually work contexts into it. low priority.
		err := backoff.Retry(op, backoff.NewExponentialBackOff())
		if err != nil {
			log.Debugf("Providing failed after number of retries: %s", err)
			return err
		}
	}
	return nil
}

func (rp *Reprovider) Trigger(ctx context.Context) error {
	progressCtx, done := context.WithCancel(ctx)
	select {
	case <-rp.ctx.Done():
		return context.Canceled
	case <-ctx.Done():
		return context.Canceled
	case rp.trigger <- done:
		<-progressCtx.Done()
		return nil
	}
}
