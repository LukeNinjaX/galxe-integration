package rug

import (
	"database/sql"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/contracts/uniswapv2"
	"github.com/artela-network/galxe-integration/onchain"

	log "github.com/sirupsen/logrus"
)

type Rug struct {
	*onchain.Base
	swap *uniswapv2.UniswapV2
	conf *config.RugConfig
}

func NewRug(db *sql.DB, conf *config.RugConfig) (*Rug, error) {
	conf.FillDefaults()

	if !conf.Enable {
		return &Rug{}, nil
	}

	base, err := onchain.NewBase(db, &conf.OnChain, false)
	if err != nil {
		return nil, err
	}

	f := &Rug{
		Base: base,
		conf: conf,
	}

	f.refreshContract()

	base.RegisterGetTasks(f.getTasks)
	base.RegisterSend(f.send)
	base.RegisterUpdateTask(f.updateTask)
	base.RegisterRefreshNetwork(f.refreshContract)
	return f, nil
}

func (s *Rug) refreshContract() bool {
	contractAddress := common.HexToAddress(s.conf.ContractAddress)
	instance, err := uniswapv2.NewUniswapV2(contractAddress, s.Client())
	if err != nil {
		log.Error("rug module: load uniswapV2 failed", err)
		return false
	}
	s.swap = instance
	return true
}

func (s *Rug) send(task biz.AddressTask) (hashs []common.Hash, err error) {
	if task.AccountAddress == nil {
		log.Error("Base module: task AccountAddress is nil", task.ID)
		return nil, onchain.ErrInvalidTask
	}

	log.Debug("rug module: running rug for", task.ID)

	opts := s.DefaultOpts(&s.conf.TxConfig)
	if len(s.conf.Path) < 2 {
		panic("config .rug.path is not correct")
	}

	path := make([]common.Address, 2)
	path[0] = common.HexToAddress(s.conf.Path[0])
	path[1] = common.HexToAddress(s.conf.Path[1])

	fromAddress := crypto.PubkeyToAddress(*s.Publickey())
	toAddress := fromAddress // rug tokens to the sender
	amount := big.NewInt(1).Mul(big.NewInt(onchain.RugAmount), big.NewInt(onchain.UINT))

	log.Debugf("rug module: SwapETHForExactTokens, from %s, to %s, amount %d", fromAddress.Hex(), toAddress.Hex(), amount)
	tx, err := s.swap.SwapETHForExactTokens(
		opts, amount, path, toAddress, big.NewInt(int64(time.Now().Second())+10000),
	)
	if err != nil {
		log.Debug("rug module: submit rug failed", task.ID, err)
		return nil, err
	}

	return []common.Hash{tx.Hash()}, nil
}

func (s *Rug) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetAspectPullTask(s.DB(), count)
}

func (s *Rug) updateTask(task biz.AddressTask, hashs []common.Hash, status *uint64) error {
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
		if *status == 0 { // this task is expected to fail
			taskStatus = string(types.TaskStatusSuccess)
		}
		req.TaskStatus = &taskStatus
		log.Debugf("Rug module: updating task, %d, hash %s, status %s\n", req.ID, *req.Txs, *req.TaskStatus)
	} else {
		log.Debugf("Rug module: updating task, %d, hash %s\n", req.ID, *req.Txs)
	}

	return biz.UpdateTask(s.DB(), req)
}
