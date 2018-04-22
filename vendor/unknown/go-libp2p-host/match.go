package host

import (
	"strings"

	"unknown/go-libp2p-protocol"
	semver "unknown/go-semver/semver"
)

func MultistreamSemverMatcher(base protocol.ID) (func(string) bool, error) {
	parts := strings.Split(string(base), "/")
	vers, err := semver.NewVersion(parts[len(parts)-1])
	if err != nil {
		return nil, err
	}

	return func(check string) bool {
		chparts := strings.Split(check, "/")
		if len(chparts) != len(parts) {
			return false
		}

		for i, v := range chparts[:len(chparts)-1] {
			if parts[i] != v {
				return false
			}
		}

		chvers, err := semver.NewVersion(chparts[len(chparts)-1])
		if err != nil {
			return false
		}

		return vers.Major == chvers.Major && vers.Minor >= chvers.Minor
	}, nil
}