package goclient

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

func ReadKey(keyFile string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privKeyHex, err := os.ReadFile(keyFile)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	privateKey, err := crypto.HexToECDSA(string(privKeyHex)[2:])
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("error casting public key to ECDSA")
		return nil, nil, err
	}
	return privateKey, publicKeyECDSA, nil
}

func SignTransaction(key *ecdsa.PrivateKey, chainId *big.Int, tx *types.Transaction) (*types.Transaction, error) {
	// keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	if chainId == nil {
		return nil, errors.New("not authorized to sign this account")
	}

	signer := types.LatestSignerForChainID(chainId)
	// if from != keyAddr {
	// 	return nil, ErrNotAuthorized
	// }
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(signer, signature)
}
