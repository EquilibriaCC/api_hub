package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"teamAPI/config"
	"time"
)

var RateLimitCounter = make(map[string][]int64)

func rateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := RateLimitCounter[c.ClientIP()]; ok {
			if time.Now().Unix() > RateLimitCounter[c.ClientIP()][0] {
				RateLimitCounter[c.ClientIP()][0] = time.Now().Unix() + int64(time.Minute.Seconds() * config.RateLimitTime)
				RateLimitCounter[c.ClientIP()][1] = 1
			} else if RateLimitCounter[c.ClientIP()][0] > time.Now().Unix() && RateLimitCounter[c.ClientIP()][1] > config.RateLimitRequests {
				c.JSON(http.StatusGatewayTimeout, map[string]string{"error":"rate limit"})
				return
			}
			RateLimitCounter[c.ClientIP()][1]++
		} else {
			RateLimitCounter[c.ClientIP()] = []int64{time.Now().Unix() + int64(time.Minute.Seconds() * config.RateLimitTime), 0}
		}
		c.Next()
	}
}

var TotalPageVisits = 0
var PageVisits = make(map[string]int)

func count() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := PageVisits[c.FullPath()]; ok {
			PageVisits[c.FullPath()]++
		} else {
			PageVisits[c.FullPath()] = 1
		}
		TotalPageVisits++
		if TotalPageVisits % 250 == 0 {
			go writePageVisits()
		}
		c.Next()
	}
}

func writePageVisits() {

}
