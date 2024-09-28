package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// App struct
type App struct {
	ctx context.Context
}

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

func (a *App) GetPdf(startPage int, endPage int, bookName string) {
	fmt.Printf("开始下载书籍: %s, 从第 %d 页到第 %d 页\n", bookName, startPage, endPage)

	// 创建临时文件夹
	tempDir, err := os.MkdirTemp("", bookName)
	if err != nil {
		log.Fatalf("无法创建临时文件夹: %v", err)
	}
	defer os.RemoveAll(tempDir) // 在函数结束时删除临时文件夹

	fmt.Printf("创建了临时文件夹: %s\n", tempDir)

	// 下载图片并保存到临时文件夹
	for page := startPage; page <= endPage; page++ {
		pageStr := strconv.Itoa(page)

		// 图片 URL
		imageURL := "https://docs.historyrussia.org/pages/" + pageStr + "/zooms/" + strconv.Itoa(4)

		// 本地保存路径
		filePath := filepath.Join(tempDir, pageStr+".jpg")

		// 调用下载函数
		err := downloadFile(imageURL, filePath)
		if err != nil {
			fmt.Printf("下载页面 %s 时发生错误: %v\n", pageStr, err)
		} else {
			//fmt.Printf("页面 %s 下载成功!\n", pageStr)
			fmt.Printf("\r下载进度：%.2f%%", float64(page-startPage)/float64(endPage-startPage)*100)
			runtime.EventsEmit(a.ctx, "progress", fmt.Sprintf("%.2f%%", float64(page-startPage)/float64(endPage-startPage)*100))
		}
	}

	fmt.Println()

	// 获取临时文件夹中的所有图片路径
	var inputPaths []string
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		log.Fatalf("读取临时文件夹时发生错误: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(tempDir, file.Name())
			if filepath.Ext(filePath) == ".jpg" || filepath.Ext(filePath) == ".png" {
				inputPaths = append(inputPaths, filePath)
			}
		}
	}

	// 将图片合成 PDF
	outputPDF := bookName + ".pdf"
	imagesToPdf(inputPaths, outputPDF)

	runtime.EventsEmit(a.ctx, "complete")
}
