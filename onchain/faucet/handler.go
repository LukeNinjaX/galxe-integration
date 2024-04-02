package faucet

import (
	"database/sql"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/contracts/rug"
	"github.com/artela-network/galxe-integration/onchain"

	log "github.com/sirupsen/logrus"
)

type Faucet struct {
	*onchain.Base
	rug  *rug.Rug
	conf *config.FaucetConfig
}

func NewFaucet(db *sql.DB, conf *config.FaucetConfig) (*Faucet, error) {
	conf.FillDefaults()

	base, err := onchain.NewBase(db, &conf.OnChain)
	if err != nil {
		return nil, err
	}

	f := &Faucet{
		Base: base,
		conf: conf,
	}

	f.refreshContract()

	base.RegisterGetTasks(f.mockGetTasks) // base.RegisterGetTasks(f.getTasks)
	base.RegisterSend(f.send)
	base.RegisterUpdateTask(f.mockUpdateTask) // base.RegisterUpdateTask(f.updateTask)
	base.RegisterRefreshNetwork(f.refreshContract)
	return f, nil
}

func (s *Faucet) refreshContract() bool {
	instance, err := rug.NewRug(common.HexToAddress(s.conf.RugAddress), s.Client())
	if err != nil {
		log.Debug("faucet module: reload rug contract failed", err)
		return false
	}
	s.rug = instance
	return true
}

func (s *Faucet) send(task biz.AddressTask) (hashs []common.Hash, err error) {
	hashTransfer, err := s.Client().Transfer(
		s.Privatekey(),
		common.HexToAddress(*task.AccountAddress),
		s.conf.TransferAmount,
		s.GetNonce(),
		&s.conf.TxConfig,
	)
	if err != nil {
		return nil, err
	}

	opts := s.DefaultOpts(&s.conf.TxConfig)
	fromAddress := crypto.PubkeyToAddress(*s.Publickey())
	toAddress := common.HexToAddress(*task.AccountAddress)
	rugAmount := big.NewInt(1).Mul(big.NewInt(s.conf.RugAmount), big.NewInt(onchain.UINT))
	log.Debugf("faucet module: transfering rug from %s to %s for %d, amount %d", fromAddress.Hex(), toAddress.Hex(), task.ID, rugAmount)
	txRug, err := s.rug.Transfer(opts, toAddress, rugAmount)
	if err != nil {
		return nil, err
	}
	return []common.Hash{hashTransfer, txRug.Hash()}, nil
}

func (s *Faucet) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetFaucetTask(s.DB(), count)
}

func (s *Faucet) updateTask(task biz.AddressTask, hashs []common.Hash, status *uint64) error {
	hexs := make([]string, len(hashs))
	for i, hash := range hashs {
		hexs[i] = hash.Hex()
	}
	memo := strings.Join(hexs, ",")

	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	req.Txs = &memo
	if status != nil {
		taskStatus := string(types.TaskStatusFail)
		if *status == 1 {
			taskStatus = string(types.TaskStatusSuccess)
		}
		req.TaskStatus = &taskStatus
		log.Debugf("faucet module: updating task, %d, hash %s, status %s\n", req.ID, *req.Txs, *req.TaskStatus)
	} else {
		log.Debugf("faucet module: updating task, %d, hash %s\n", req.ID, *req.Txs)
	}

	return biz.UpdateTask(s.DB(), req)
}
