package api

import (
	"auth/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"xeq_hub/config"
	"time"
)

func API() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(count())
	r.Use(rateLimit())
	Endpoints(r)
	s := &http.Server{
		Addr:           config.APIPort,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
