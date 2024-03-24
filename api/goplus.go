package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"

	"github.com/artela-network/galxe-integration/api/biz"
)

type UrlInput struct {
	AccountAddress string `form:"accountAddress" json:"accountAddress" binding:"required"`
}

func (s *Server) getTasks(c *gin.Context) {
	accountAddress := c.Param("accountAddress")
	taskId := c.Param("taskId")

	tasks, err := biz.GetAccountTaskInfo(s.db, accountAddress, taskId)
	if err != nil {
		log.Errorf("Failed to getTasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get tasks " + err.Error(),
		})
		return
	}
	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tasks,
	})
}
func (s *Server) newTasks(c *gin.Context) {
	input := &biz.InitTaskQuery{}
	if errA := c.ShouldBindBodyWith(&input, binding.JSON); errA != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to new tasks Should BindBody" + errA.Error(),
		})
		return
	}
	tasks, getErr := biz.GetTasks(s.db, &biz.TaskQuery{AccountAddress: input.AccountAddress, ChannelTaskId: input.ChannelTaskId})
	if getErr != nil {
		log.Errorf("Failed to getTasks: %v", getErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get tasks " + getErr.Error(),
		})
		return
	}

	if len(tasks) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Already have tasks",
			"data":    tasks,
		})
		return
	}

	getErr = biz.InitTask(s.db, input)
	// insert db
	if getErr != nil {
		log.Errorf("Failed to query database: %v", getErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to new task",
		})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
	return
}

func (s *Server) updateTask(c *gin.Context) {
	taskUpQuery := biz.UpdateTaskQuery{}

	if errA := c.ShouldBindBodyWith(&taskUpQuery, binding.JSON); errA == nil {
		err := biz.UpdateTask(s.db, taskUpQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to update task " + err.Error(),
			})
			return
		}
		// 返回查询结果
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})

		// 这时, 复用存储在上下文中的 body。
	} else {
		log.Errorf("Failed to bind body: %v", errA)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to bind body " + errA.Error(),
		})
		return
	}

}

func (s *Server) rugPullInfo(c *gin.Context) {
	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"completed": true,
	})
}

func (s *Server) syncStatus(c *gin.Context) {

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"completed": true,
	})
}
func (s *Server) faucet(c *gin.Context) {
	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"completed": true,
	})
}
