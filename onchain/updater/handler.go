package updater

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/artela-network/galxe-integration/api/biz"
	"github.com/artela-network/galxe-integration/api/types"
	"github.com/artela-network/galxe-integration/config"
	"github.com/artela-network/galxe-integration/onchain"

	log "github.com/sirupsen/logrus"
)

type Updater struct {
	*onchain.Base
	conf *config.UpdaterConfig
}

func NewUpdater(db *sql.DB, conf *config.UpdaterConfig) (*Updater, error) {
	conf.FillDefaults()

	base, err := onchain.NewBase(db, &conf.OnChain)
	if err != nil {
		return nil, err
	}

	f := &Updater{
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

func (s *Updater) send(task biz.AddressTask) (hashs []common.Hash, err error) {
	if task.Txs == nil {
		return nil, errors.New("tasks is not valid, txs cannot be empty")
	}
	hash := common.HexToHash(*task.Txs)
	return []common.Hash{hash}, nil
}

func (s *Updater) refreshContract() bool {
	return true
}

func (s *Updater) getTasks(count int) ([]biz.AddressTask, error) {
	return biz.GetAddLiquidityTask(s.DB(), count)
}

func (s *Updater) updateTask(task biz.AddressTask, hashs []common.Hash, status *uint64) error {
	hexs := make([]string, len(hashs))
	for i, hash := range hashs {
		hexs[i] = hash.Hex()
	}
	memo := strings.Join(hexs, ",")

	req := &biz.UpdateTaskQuery{}
	req.ID = task.ID
	if status != nil {
		taskStatus := string(types.TaskStatusFail)
		if *status == 0 { // this task is expected to fail
			taskStatus = string(types.TaskStatusSuccess)
		}
		req.TaskStatus = &taskStatus
		log.Debugf("Updater module: updating task, %d, hash %s, status %s\n", req.ID, memo, *req.TaskStatus)
	} else {
		log.Debugf("Updater module: updating task, %d, hash %s\n", req.ID, memo)
	}

	return biz.UpdateTask(s.DB(), req)
}
