package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	admin "go-5m3Mirco/app/shop/admin"
)

func main() {
	//创建随机数源
	source := rand.NewSource(time.Now().UnixNano())
	//创建随机数生成器
	_ = rand.New(source)

	//自动将程序的并发性能配置到当前硬件的最佳状态
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	//启动
	admin.NewApp("admin-server").Run()
}
