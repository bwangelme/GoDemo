package main

import (
	"image"
	"image/draw"
	"os"

	"github.com/chai2010/webp"
)

func main() {
	dir := "/Users/mico/Github/GoDemo/imgpaste/img"
	os.Chdir(dir)
	// 读取背景图片
	bgFile, err := os.Open("background.webp")
	if err != nil {
		panic(err)
	}
	defer bgFile.Close()

	bgImg, err := webp.Decode(bgFile)
	if err != nil {
		panic(err)
	}

	// 读取前景图片
	fgFile, err := os.Open("foreground.webp")
	if err != nil {
		panic(err)
	}
	defer fgFile.Close()

	fgImg, err := webp.Decode(fgFile)
	if err != nil {
		panic(err)
	}

	// 创建一个与背景图片相同大小的新图像
	bounds := bgImg.Bounds()
	combined := image.NewRGBA(bounds)

	// 先绘制背景图片
	draw.Draw(combined, bounds, bgImg, image.Point{}, draw.Src)

	// 计算前景图片的位置（居中）
	offsetX := (bounds.Dx() - fgImg.Bounds().Dx()) / 2
	offsetY := (bounds.Dy() - fgImg.Bounds().Dy()) / 2
	offset := image.Pt(offsetX, offsetY)

	// 在前景图片上绘制叠加效果（使用Over操作符实现透明度叠加）
	draw.Draw(combined, fgImg.Bounds().Add(offset), fgImg, image.Point{}, draw.Over)

	// 保存合并后的图片
	output, err := os.Create("overlay.webp")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	// 设置WebP编码选项
	options := &webp.Options{Lossless: true, Quality: 90}

	if err := webp.Encode(output, combined, options); err != nil {
		panic(err)
	}
}
