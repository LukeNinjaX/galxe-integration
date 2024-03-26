package goclient

import (
	"encoding/json"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	jose "github.com/dvsekhvalnov/jose2go"
	log "github.com/sirupsen/logrus"
)

// Item is a thing stored on the keyring
type Item struct {
	Key         string
	Data        []byte
	Label       string
	Description string

	// Backend specific config
	KeychainNotTrustApplication bool
	KeychainNotSynchronizable   bool
}

func readKeyStore() []byte {
	home := os.Getenv("HOME")
	filename := path.Join(home, "/.artelad/keyring-test/mykey.info")

	bytes, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		log.Error("file not exist", filename)
		return nil
	} else if err != nil {
		log.Error(err)
		return nil
	}

	payload, _, err := jose.Decode(string(bytes), "test")
	if err != nil {
		log.Error(err)
		return nil
	}

	decoded := &Item{}
	err = json.Unmarshal([]byte(payload), decoded)
	if err != nil {
		log.Error(err)
	}

	record := unmarshalRecord(decoded.Data)
	key := record.GetLocal().PrivKey
	return key.Value
}

func unmarshalRecord(data []byte) *keyring.Record {
	record := &keyring.Record{}
	if err := record.Unmarshal(data); err != nil {
		log.Error(err)
		return nil
	}
	return record
}
