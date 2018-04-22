package republisher

import (
	"context"
	"errors"
	"time"

	ic "github.com/AtlantPlatform/go-ipfs/go-libp2p-crypto"
	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	keystore "github.com/AtlantPlatform/go-ipfs/keystore"
	namesys "github.com/AtlantPlatform/go-ipfs/namesys"
	pb "github.com/AtlantPlatform/go-ipfs/namesys/pb"
	path "github.com/AtlantPlatform/go-ipfs/path"
	ds "unknown/go-datastore"
	dshelp "unknown/go-ipfs-ds-help"
	recpb "unknown/go-libp2p-record/pb"
	routing "unknown/go-libp2p-routing"
	logging "unknown/go-log"
	proto "unknown/gogo-protobuf/proto"
	goprocess "unknown/goprocess"
	gpctx "unknown/goprocess/context"
)

var errNoEntry = errors.New("no previous entry")

var log = logging.Logger("ipns-repub")

// DefaultRebroadcastInterval is the default interval at which we rebroadcast IPNS records
var DefaultRebroadcastInterval = time.Hour * 4

// InitialRebroadcastDelay is the delay before first broadcasting IPNS records on start
var InitialRebroadcastDelay = time.Minute * 1

// FailureRetryInterval is the interval at which we retry IPNS records broadcasts (when they fail)
var FailureRetryInterval = time.Minute * 5

// DefaultRecordLifetime is the default lifetime for IPNS records
const DefaultRecordLifetime = time.Hour * 24

type Republisher struct {
	r    routing.ValueStore
	ds   ds.Datastore
	self ic.PrivKey
	ks   keystore.Keystore

	Interval time.Duration

	// how long records that are republished should be valid for
	RecordLifetime time.Duration
}

// NewRepublisher creates a new Republisher
func NewRepublisher(r routing.ValueStore, ds ds.Datastore, self ic.PrivKey, ks keystore.Keystore) *Republisher {
	return &Republisher{
		r:              r,
		ds:             ds,
		self:           self,
		ks:             ks,
		Interval:       DefaultRebroadcastInterval,
		RecordLifetime: DefaultRecordLifetime,
	}
}

func (rp *Republisher) Run(proc goprocess.Process) {
	timer := time.NewTimer(InitialRebroadcastDelay)
	defer timer.Stop()
	if rp.Interval < InitialRebroadcastDelay {
		timer.Reset(rp.Interval)
	}

	for {
		select {
		case <-timer.C:
			timer.Reset(rp.Interval)
			err := rp.republishEntries(proc)
			if err != nil {
				log.Error("Republisher failed to republish: ", err)
				if FailureRetryInterval < rp.Interval {
					timer.Reset(FailureRetryInterval)
				}
			}
		case <-proc.Closing():
			return
		}
	}
}

func (rp *Republisher) republishEntries(p goprocess.Process) error {
	ctx, cancel := context.WithCancel(gpctx.OnClosingContext(p))
	defer cancel()

	err := rp.republishEntry(ctx, rp.self)
	if err != nil {
		return err
	}

	if rp.ks != nil {
		keyNames, err := rp.ks.List()
		if err != nil {
			return err
		}
		for _, name := range keyNames {
			priv, err := rp.ks.Get(name)
			if err != nil {
				return err
			}
			err = rp.republishEntry(ctx, priv)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (rp *Republisher) republishEntry(ctx context.Context, priv ic.PrivKey) error {
	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return err
	}

	log.Debugf("republishing ipns entry for %s", id)

	// Look for it locally only
	_, ipnskey := namesys.IpnsKeysForID(id)
	p, seq, err := rp.getLastVal(ipnskey)
	if err != nil {
		if err == errNoEntry {
			return nil
		}
		return err
	}

	// update record with same sequence number
	eol := time.Now().Add(rp.RecordLifetime)
	err = namesys.PutRecordToRouting(ctx, priv, p, seq, eol, rp.r, id)
	if err != nil {
		return err
	}

	return nil
}

func (rp *Republisher) getLastVal(k string) (path.Path, uint64, error) {
	ival, err := rp.ds.Get(dshelp.NewKeyFromBinary([]byte(k)))
	if err != nil {
		// not found means we dont have a previously published entry
		return "", 0, errNoEntry
	}

	val := ival.([]byte)
	dhtrec := new(recpb.Record)
	err = proto.Unmarshal(val, dhtrec)
	if err != nil {
		return "", 0, err
	}

	// extract published data from record
	e := new(pb.IpnsEntry)
	err = proto.Unmarshal(dhtrec.GetValue(), e)
	if err != nil {
		return "", 0, err
	}
	return path.Path(e.Value), e.GetSequence(), nil
}
