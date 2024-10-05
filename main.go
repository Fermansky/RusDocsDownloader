package main

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// 设置日志格式为JSON，便于日志处理和分析
	log.SetFormatter(&log.TextFormatter{})

	// 设置日志级别
	log.SetLevel(log.InfoLevel)

	// 创建一个日志文件，所有日志会写入这个文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		// 如果创建文件失败，日志输出到标准输出
		log.SetOutput(os.Stdout)
		log.Warn("Failed to log to file, using default stdout")
	}
}

func main() {

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "RusDocsDownloader",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 56, G: 0, B: 9, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
