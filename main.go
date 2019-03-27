package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiaomLee/gowebsvr/core/config"
	"github.com/xiaomLee/gowebsvr/server"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var db = make(map[string]string)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := server.Init(); err != nil {
		fmt.Println("octopus server startup failed! error:", err.Error())
		os.Exit(1)
	}
	Run()
}

func Run() {
	// 框架初始化
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(server.RequestStart()) // 请求开始，勿调整顺序
	engine.Use(gin.Logger())
	engine.Use(server.Recovery())

	//if config.Instance().Base.Mock {
	//	// 初始化第三方mock接口
	//	engine.Use(server.MockThirdPartCall())
	//}
	addUrls(engine)
	//启动服务
	var svrConf struct {
		Name      string
		Port      int
		TimeoutMs int64
	}
	config.ReadConfig("common", "server", &svrConf)
	//engine.Run(":" + strconv.Itoa(svrConf.Port))
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(svrConf.Port),
		Handler:        engine,
		ReadTimeout:    time.Duration(svrConf.TimeoutMs) * time.Millisecond,
		WriteTimeout:   time.Duration(svrConf.TimeoutMs) * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}

func addUrls(r *gin.Engine) {
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

}
