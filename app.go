package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

// App struct
type App struct {
	ctx context.Context
}

var semaphore = make(chan struct{}, 10) // 限制最大并发 Goroutines 数量为 10

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Scrape(url string) *BookInfo {
	return Scrape(url, 0)
}

// 打开本地文件或文件夹
func (a *App) OpenFileOrFolder(path string) error {
	var cmd *exec.Cmd

	// 仅考虑windows系统
	cmd = exec.Command("explorer", path)

	return cmd.Start()
}

func (a *App) Update() {
	err := ApplyUpdate()
	if err != nil {
		return
	}
}

func (a *App) CheckUpdate() string {
	update, err := CheckUpdate()
	if err != nil {

	}
	return update
}

func (a *App) GetPdf(pages []int, bookName string, quality int) error {

	// 创建临时文件夹
	tempDir, err := os.MkdirTemp("", bookName)
	if err != nil {
		log.Errorf("无法创建临时文件夹: %v", err)
		return err
	}
	defer os.RemoveAll(tempDir) // 在函数结束时删除临时文件夹

	log.Infof("创建了临时文件夹: %s\n", tempDir)
	var wg sync.WaitGroup
	var mu sync.Mutex // 保护共享资源

	downloaded := 0
	totalPages := len(pages)
	errorChan := make(chan error, totalPages)

	// 并发下载每个页面的图片
	for i, page := range pages {
		wg.Add(1)
		go func(i, page int) {
			defer wg.Done()

			semaphore <- struct{}{}        // 占用一个并发槽位
			defer func() { <-semaphore }() // 释放并发槽位

			pageStr := strconv.Itoa(page)
			imageURL := "https://docs.historyrussia.org/pages/" + pageStr + "/zooms/" + strconv.Itoa(quality)
			filePath := filepath.Join(tempDir, fmt.Sprintf("%05d.jpg", i))

			// 调用下载函数
			err := downloadFileWithRetry(imageURL, filePath)
			if err != nil {
				log.Errorf("下载页面 %s 时发生错误: %v\n", pageStr, err)
				errorChan <- err
				return
			}

			mu.Lock()
			downloaded++
			progress := float64(downloaded) / float64(totalPages) * 100
			mu.Unlock()

			// 发送下载进度到前端
			runtime.EventsEmit(a.ctx, "progress", fmt.Sprintf("%.2f%%", progress))
		}(i, page)
	}

	wg.Wait() // 等待所有下载任务完成
	close(errorChan)

	// 检查是否有下载错误
	if len(errorChan) > 0 {
		log.Errorf("下载过程中出现错误")
		return fmt.Errorf("下载过程中出现错误")
	}

	// 获取临时文件夹中的所有图片路径
	inputPaths, err := getImagePaths(tempDir)
	if err != nil {
		log.Errorf("读取临时文件夹时发生错误: %v", err)
		return err
	}
	log.Infof("临时文件夹中共有%d张图片", len(inputPaths))

	// 将图片合成 PDF
	outputPDF := bookName + ".pdf"
	path := imagesToPdf(inputPaths, outputPDF)

	// 通知前端任务完成
	runtime.EventsEmit(a.ctx, "complete", path)
	return nil
}

// getImagePaths 获取指定文件夹中的所有图片路径
func getImagePaths(dir string) ([]string, error) {
	var inputPaths []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dir, file.Name())
			if filepath.Ext(filePath) == ".jpg" || filepath.Ext(filePath) == ".png" {
				inputPaths = append(inputPaths, filePath)
			}
		}
	}
	return inputPaths, nil
}
