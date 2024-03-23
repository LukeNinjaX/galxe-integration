package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"

	"github.com/artela-network/galxe-integration/api/dao"
)

func (s *Server) getTasks(c *gin.Context) {
	addr := c.Param("address")
	tasks, err := dao.GetTasks(s.db, addr)
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
	address := c.Query("address")
	tasks, getErr := dao.GetTasks(s.db, address)
	if getErr != nil {
		log.Errorf("Failed to getTasks: %v", getErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get tasks " + getErr.Error(),
		})
		return
	}

	if tasks != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Already have tasks",
			"data":    tasks,
		})
		return
	}

	// insert db
	_, err := s.db.Exec("INSERT INTO address_tasks (account_address, task_name, task_status) VALUES ($1, 'AddLiquidity', '0'),($1, 'AspectPull', '0'),($1, 'RugPull', '0');", address)
	if err != nil {
		log.Errorf("Failed to query database: %v", err)
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
}

func (s *Server) updateTask(c *gin.Context) {
	taskUpQuery := dao.UpdateTaskQuery{}

	if errA := c.ShouldBindBodyWith(&taskUpQuery, binding.JSON); errA == nil {

		c.String(http.StatusOK, `the body should be UpdateTaskQuery`)
		err := dao.UpdateTask(s.db, taskUpQuery)
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
