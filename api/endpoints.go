package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xeq_hub/tasks"
)

func Endpoints(r *gin.Engine) {

	r.GET("/page_visits", func(c *gin.Context) { c.JSON(200, map[string]interface{}{"data": TotalPageVisits}) })
	r.GET("/page_visits/all", func(c *gin.Context) { c.JSON(200, map[string]interface{}{"data": PageVisits}) })
	r.GET("/emission", func(c *gin.Context) { c.JSON(200, map[string]interface{}{"data": tasks.Supply}) })
	r.GET("/number_of_oracle_nodes", func(c *gin.Context) { c.JSON(200, map[string]interface{}{"data": tasks.NumberOfOracleNodes()}) })
	r.GET("/oracle_node_history", func(c *gin.Context) { c.JSON(200, map[string]interface{}{"data": tasks.NumberOfOracleNodes()}) })

	r.GET("/tx_pool", func(c *gin.Context) {
		c.JSON(http.StatusOK, tasks.GetTxPool())
	})

	r.GET("/tx/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		c.JSON(http.StatusOK, tasks.SearchTx(hash))
	})

	r.GET("/transactions/:start_height", func(c *gin.Context) {
		startHeight := c.Param("start_height")
		height, _ := strconv.Atoi(startHeight)
		c.JSON(http.StatusOK, tasks.Transactions(height))

	})

}
