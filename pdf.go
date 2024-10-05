package main

import (
	"fmt"
	"github.com/signintech/gopdf"
	log "github.com/sirupsen/logrus"
	"image"
	"os"
	filepath2 "path/filepath"
)

// 函数用于读取图片的宽高
func getImageDimensions(imagePath string) (float64, float64, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	// 返回图片的宽度和高度
	return float64(img.Width), float64(img.Height), nil
}

func imagesToPdf(imagePaths []string, filepath string) string {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{})

	for index, imagePath := range imagePaths {

		// 读取图片的宽高
		width, height, err := getImageDimensions(imagePath)
		if err != nil {
			log.Errorf("Error getting image dimensions for %s: %v", imagePath, err)
			continue // 如果发生错误，跳过当前图片
		}

		width *= 72.0 / 150.0
		height *= 72.0 / 150.0

		// 添加新页面，并根据图片的尺寸设置页面大小
		pdf.AddPageWithOption(gopdf.PageOption{PageSize: &gopdf.Rect{W: width, H: height}})

		// 将图片插入 PDF 中
		err = pdf.Image(imagePath, 0, 0, &gopdf.Rect{W: width, H: height})
		if err != nil {
			log.Infof("Error inserting image %s into PDF: %v", imagePath, err)
			continue // 如果发生错误，跳过当前图片
		}

		fmt.Printf("\rPDF转换进度%.2f%%", float64(index+1)/float64(len(imagePaths))*100)
	}

	// 保存 PDF 文件
	err := pdf.WritePdf(filepath)
	if err != nil {
		log.Fatalf("Error saving PDF to %s: %v", filepath, err)
	}

	log.Infof("PDF 已生成: %s", filepath)
	path, _ := filepath2.Abs(filepath)

	return path
}
