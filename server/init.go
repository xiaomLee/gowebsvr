package server

import (
	"fmt"
	"github.com/xiaomLee/gowebsvr/core/config"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"
)

func Init() error {
	rand.Seed(time.Now().UnixNano())

	// 参数初始化
	binDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("get bin dir err: %s", err.Error())
		return err
	}
	rootDir := path.Dir(binDir)

	// 配置初始化
	confDir := path.Join(rootDir, "/conf")

	config.LoadConfigFile("common", confDir, "toml")
	// 日志初始化
	//logDir := path.Join(rootDir, "/log")

	// i18n 初始化语言包

	fmt.Println("server init...")
	return nil
}
