package utils

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

// ImageInfo 获取图像信息的结构体
type ImageInfo struct {
	Bytes  []byte      // 图像字节数据
	Width  int         // 宽度
	Height int         // 高度
	Mode   string      // 模式 (RGB, RGBA, etc.)
	Format string      // 格式 (JPEG, PNG, etc.)
	Image  image.Image // 原始图像对象
}

// 实现String()方法以便可以打印ImageInfo结构体
func (info ImageInfo) String() string {
	return fmt.Sprintf("ImageInfo{Width: %d, Height: %d, Mode: %s, Format: %s, BytesLength: %d}",
		info.Width, info.Height, info.Mode, info.Format, len(info.Bytes))
}

// ImageFromBytes 从字节数组中创建图像
func ImageFromBytes(data []byte, width, height int, mode string) image.Image {
	switch mode {
	case "RGB":
		// RGB格式：每像素3字节 (R, G, B)
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := range height {
			for x := range width {
				idx := (y*width + x) * 3
				if idx+2 < len(data) {
					r := data[idx]
					g := data[idx+1]
					b := data[idx+2]
					img.Set(x, y, color.RGBA{r, g, b, 255})
				}
			}
		}
		return img
	case "RGBA":
		// RGBA格式：每像素4字节 (R, G, B, A)
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		copy(img.Pix, data)
		return img
	default:
		return nil
	}
}

// ResizeImage 调整图像大小
func ResizeImage(src image.Image, newWidth, newHeight int) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 使用CatmullRom插值算法进行高质量缩放
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}

// DetermineColorMode 确定颜色模式
func DetermineColorMode(img image.Image) string {
	bounds := img.Bounds()
	x, y := bounds.Min.X, bounds.Min.Y

	// 检查一个像素点来确定颜色模式
	c := img.At(x, y)
	_, _, _, a := c.RGBA()

	// 如果alpha值不是255，则是RGBA模式
	if a>>8 != 255 {
		return "RGBA"
	}
	return "RGB"
}

// GetImageBytes 从图像获取字节数据
func GetImageBytes(img image.Image) []byte {
	bounds := img.Bounds()

	// 转换为RGBA格式以统一处理
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 提取字节数据 (RGBA格式: 每像素4字节)
	return rgba.Pix
}

// GetImageInfo 获取图像信息
func GetImageInfo(img image.Image, format string) ImageInfo {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 获取字节数据
	bytes := GetImageBytes(img)

	// 确定颜色模式
	mode := DetermineColorMode(img)

	return ImageInfo{
		Bytes:  bytes,
		Width:  width,
		Height: height,
		Mode:   mode,
		Format: format,
		Image:  img,
	}
}
