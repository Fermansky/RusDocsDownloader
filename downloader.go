package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

type progressReader struct {
	io.Reader
	total   int64
	current int64
}

const maxRetries = 3               // 最大重试次数
const retryDelay = 2 * time.Second // 每次重试之间的等待时间

// 下载文件，并添加重试机制
func downloadFileWithRetry(url string, filePath string) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = downloadFile(url, filePath)
		if err == nil {
			return nil // 下载成功
		}
		log.Errorf("下载失败 (尝试 %d/%d): %v", attempt, maxRetries, err)

		if attempt < maxRetries {
			time.Sleep(retryDelay) // 等待一段时间再重试
		}
	}
	return fmt.Errorf("下载失败，已达到最大重试次数: %v", err)
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	if err == nil {
		p.current += int64(n)
	}
	return n, err
}

func downloadFile(url, filepath string) error {
	// 发送 HTTP GET 请求获取图片
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}

	// 获取文件大小
	total := response.ContentLength
	if total <= 0 {
		log.Error("unable to determine file size")
		return fmt.Errorf("unable to determine file size")
	}

	// 创建本地文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用自定义的进度显示 reader
	progress := &progressReader{
		Reader: response.Body,
		total:  total,
	}

	// 将 HTTP 响应体中的数据写入文件
	_, err = io.Copy(file, progress)
	if err != nil {
		return err
	}

	return nil
}
