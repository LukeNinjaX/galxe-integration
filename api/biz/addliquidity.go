package biz

import "database/sql"

type AddLiquidityReq struct {
	db *sql.DB
}

func (s *AddLiquidityReq) doSend(limit int) error {
	/**

	task, err := GetAddLiquidityTask(s.db, limit)
	if err != nil {
		return err
	}
	for i, addressTask := range task {
		// 排队发送交易
		// 更新 set task_status='2'  by   task[0].TaskId
		// /UpdateTask(s.db, &UpdateTaskQuery{TaskId:})
	}
	*
	*/

	return nil
}
