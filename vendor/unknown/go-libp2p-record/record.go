package record

import (
	"bytes"

	ci "github.com/AtlantPlatform/go-ipfs/go-libp2p-crypto"
	u "unknown/go-ipfs-util"
	pb "unknown/go-libp2p-record/pb"
	proto "unknown/gogo-protobuf/proto"
)

// MakePutRecord creates and signs a dht record for the given key/value pair
func MakePutRecord(sk ci.PrivKey, key string, value []byte, sign bool) (*pb.Record, error) {
	record := new(pb.Record)

	record.Key = proto.String(string(key))
	record.Value = value

	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	pkh := u.Hash(pkb)

	record.Author = proto.String(string(pkh))
	if sign {
		blob := RecordBlobForSig(record)

		sig, err := sk.Sign(blob)
		if err != nil {
			return nil, err
		}

		record.Signature = sig
	}
	return record, nil
}

// RecordBlobForSig returns the blob protected by the record signature
func RecordBlobForSig(r *pb.Record) []byte {
	k := []byte(r.GetKey())
	v := []byte(r.GetValue())
	a := []byte(r.GetAuthor())
	return bytes.Join([][]byte{k, v, a}, []byte{})
}
