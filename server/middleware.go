package server

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
)

func RequestStart() gin.HandlerFunc {
	return func(c *gin.Context) {
		rq := c.Request.URL.Query()
		fmt.Print(rq)
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Print(err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}