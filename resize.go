package imgResize

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/imaging"
)

// /////////
// 批量处理图片
// paths: 原图片路径
// newPaths: 新图片路径
// formats: 图片格式
// maxWidth: 最大宽度
// maxHeight: 最大高度
// quality: 图片质量,-1表示不压缩
// 返回值: 错误信息
func ImgResizes(paths []string, newPaths []string, formats []string, maxWHs []ImageWH, quality int) ([][]string, error) {

	newImagePath := [][]string{}

	for i := 0; i < len(paths); i++ {
		newpaths, err := ImgResize(paths[i], newPaths[i], formats, maxWHs, quality)
		if err != nil {
			return newImagePath, err
		}
		newImagePath = append(newImagePath, newpaths)
	}
	return newImagePath, nil
}

func ImgResize(path string, newPath string, formats []string, maxWHs []ImageWH, quality int) ([]string, error) {

	newImagePath := []string{}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("os.Open failed:", err)
		file.Close()
		return newImagePath, err
	}
	// 读取图像文件的配置信息
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("image.DecodeConfig failed:", err)
		file.Close()
		return newImagePath, err
	}
	file.Close()

	// 打印图像的格式
	fmt.Println("Image format is:", format)

	tempImage, err := imaging.Open(path)
	if err != nil {
		fmt.Println("imaging.Open failed:", err)
		return newImagePath, err
	}

	for i := 0; i < len(maxWHs); i++ {
		newImage := tempImage
		isResize := false
		format = "jpg"
		imagewh, err := DecodeImageWidthHeight(newImage, format)
		if err != nil {
			fmt.Println("DecodeImageWidthHeight failed:", err)
			return newImagePath, err
		}
		fmt.Println("Image size is:", imagewh.Width, "x", imagewh.Height)

		if imagewh.Width > maxWHs[i].Width {
			isResize = true
			newImage = imaging.Resize(newImage, maxWHs[i].Width, 0, imaging.Lanczos)
		}

		imagewh, err = DecodeImageWidthHeight(newImage, format)
		if err != nil {
			fmt.Println("DecodeImageWidthHeight failed:", err)
			return newImagePath, err
		}

		if imagewh.Height > maxWHs[i].Height {
			isResize = true
			newImage = imaging.Resize(newImage, 0, maxWHs[i].Height, imaging.Lanczos)
		}

		imgSize := ""
		switch i {
		case 0:
			imgSize = "S"
		case 1:
			imgSize = "M"
		case 2:
			imgSize = "L"
		default:
			for ii := 2; ii < i; ii++ {
				imgSize += "X"
			}
			imgSize += "L"
		}
		for _, v := range formats {
			path, err = saveImage(newImage, newPath, imgSize, v, quality)
			if err != nil {
				fmt.Println("saveImage failed:", err)
				return newImagePath, err
			}
			newImagePath = append(newImagePath, path)
		}
		if !isResize {
			break
		}
	}
	return newImagePath, nil
}
