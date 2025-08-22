package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Retry   int    `json:"retry"`
}

func loadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func main() {
	cfgPath := flag.String("config", "", "path to config.json")
	flag.Parse()

	var (
		cfg *Config
		err error
	)
	if cfg, err = loadConfig(*cfgPath); err != nil {
		fmt.Println("[config] 读取失败:", err)
	}

	fmt.Println("=== Go Binary Runner Demo ===")
	fmt.Println("这是一个Go交叉编译的示例程序")
	fmt.Println("运行在Android设备上")
	fmt.Println()

	if cfg != nil {
		fmt.Printf("[config] name=%q enabled=%v retry=%d\n", cfg.Name, cfg.Enabled, cfg.Retry)
		if !cfg.Enabled {
			fmt.Println("[config] enabled=false，程序不执行任务，直接退出")
			return
		}
	} else {
		fmt.Println("[config] 未提供配置，将使用默认行为")
	}

	// 任务数量：优先使用配置的 retry（>=1），否则默认 5
	tasks := 5
	if cfg != nil && cfg.Retry > 0 {
		tasks = cfg.Retry
	}

	// 显示一些系统信息
	fmt.Println("程序开始时间:", time.Now().Format("2006-01-02 15:04:05"))
	if cfg != nil && cfg.Name != "" {
		fmt.Println("任务名称:", cfg.Name)
	}

	// 模拟一些工作（逐秒输出，便于前端实时显示）
	for i := 1; i <= tasks; i++ {
		fmt.Printf("处理任务 %d/%d...\n", i, tasks)
		// 立刻刷新 stdout（在某些环境下有助于更快显示）
		os.Stdout.Sync()
		time.Sleep(1 * time.Second)
		fmt.Printf("任务 %d 完成\n", i)
		os.Stdout.Sync()
	}

	fmt.Println()
	fmt.Println("所有任务完成!")
	fmt.Println("程序结束时间:", time.Now().Format("2006-01-02 15:04:05"))
}