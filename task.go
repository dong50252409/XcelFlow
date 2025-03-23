package main

import (
	"fmt"
	"sync"
	"time"
	"xCelFlow/flags"
)

// 自定义配置区域 >>>>>>>>>>>>>>>>>>>>>>
const (
	maxConcurrent = 10 // 最大并发协程数
)

func processFile(path string) error {
	return Run(path, flags.SchemaName)
}

type result struct {
	FilePath string // 文件路径
	Error    error  // 错误信息
}

func runTask(filepathList []string) {
	start := time.Now()

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)
	results := make(chan result, len(filepathList))

	// 启动并发任务
	for _, filePath := range filepathList {
		wg.Add(1)
		sem <- struct{}{}

		go func(path string) {
			defer func() {
				<-sem
				wg.Done()
				if r := recover(); r != nil {
					results <- result{FilePath: path, Error: fmt.Errorf("%v", r)}
				}
			}()

			if err := processFile(path); err != nil {
				results <- result{FilePath: path, Error: err}
			}
		}(filePath)
	}

	// 等待任务完成
	go func() {
		wg.Wait()
		close(results)
		close(sem)
	}()

	// 处理结果
	success, total := 0, len(filepathList)
	for res := range results {
		if res.Error != nil {
			fmt.Printf("[错误] %s → %v\n", res.FilePath, res.Error)
			continue
		}
		success++
	}

	// 输出统计
	fmt.Printf("已处理：%d/%d个文件，耗时：%v\n", success, total, time.Since(start))
}
