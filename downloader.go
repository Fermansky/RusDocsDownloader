package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type progressReader struct {
	io.Reader
	total   int64
	current int64
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

	// 获取文件大小
	total := response.ContentLength
	if total <= 0 {
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
