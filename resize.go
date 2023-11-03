package imgResize

import (
	"fmt"
	"image"
	"os"

	"github.com/disintegration/imaging"
)

// ========================
//
//	批量处理图片
//	paths		[]string	原图片路径
//	newPaths	[]string	新图片路径
//	formats		[]string	图片格式
//	maxWHs		[]ImageWH	图片宽高
//	quality		int		图片质量
//	isPrint		bool		是否打印错误及提示信息
//	返回值		[]string	新图片路径
//	返回值		error		错误信息
func ImgResizes(paths []string, newPaths []string, formats []string, maxWHs []ImageWH, quality int, isPrint bool) ([][]string, error) {

	newImagePath := [][]string{}

	for i := 0; i < len(paths); i++ {
		newpaths, err := ImgResize(paths[i], newPaths[i], formats, maxWHs, quality, isPrint)
		if err != nil {
			return newImagePath, err
		}
		newImagePath = append(newImagePath, newpaths)
	}
	return newImagePath, nil
}

// ========================
//
//	处理图片
//	paths		string		原图片路径
//	newPaths	string		新图片路径
//	formats		[]string	图片格式
//	maxWHs		[]ImageWH	图片宽高
//	quality		int		图片质量
//	isPrint		bool		是否打印错误及提示信息
//	返回值		[]string	新图片路径
//	返回值		error		错误信息
func ImgResize(path string, newPath string, formats []string, maxWHs []ImageWH, quality int, isPrint bool) ([]string, error) {
	newImagePath := []string{}
	file, err := os.Open(path)
	if err != nil {
		// fmt.Println("os.Open failed:", err)
		file.Close()
		return newImagePath, err
	}
	// 读取图像文件的配置信息
	_, rformat, err := image.DecodeConfig(file)
	if err != nil {
		// fmt.Println("image.DecodeConfig failed:", err)
		file.Close()
		return newImagePath, err
	}
	file.Close()
	if rformat == "jpeg" {
		rformat = "jpg"
	}

	tempImage, err := imaging.Open(path)
	if err != nil {
		if isPrint {
			fmt.Println("imaging.Open failed:", err)
		}
		return newImagePath, err
	}

	if isPrint {
		fmt.Println("open image: ", path, " format is:", rformat)
	}

	sizeNamei := 0
	for i := 0; i < len(maxWHs); i++ {
		var (
			newImage image.Image = tempImage
			isResize bool        = false
			format   string      = "jpg"
			imgSize  string      = ""
			imagewh  *ImageWH    = nil
		)

		switch sizeNamei {
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
		if isPrint && i == 0 {
			// 解析图片宽高后，进行图片缩放
			imagewh, err = DecodeImageWidthHeight(newImage, format)
			if err != nil {
				if isPrint {
					fmt.Println("DecodeImageWidthHeight failed:", err)
				}
				return newImagePath, err
			}
			fmt.Println("Image size is:", imagewh.Width, "x", imagewh.Height)
		}
		if maxWHs[i].Width < 0 || maxWHs[i].Height < 0 {
			// 不进行图片缩放
			imgSize = "R"
			isResize = true
			sizeNamei--
		} else {
			// 解析图片宽高后，进行图片缩放
			imagewh, err = DecodeImageWidthHeight(newImage, format)
			if err != nil {
				if isPrint {
					fmt.Println("DecodeImageWidthHeight failed:", err)
				}
				return newImagePath, err
			}

			if imagewh.Width >= imagewh.Height {
				if imagewh.Width > maxWHs[i].Width {
					isResize = true
					newImage = imaging.Resize(newImage, maxWHs[i].Width, 0, imaging.Lanczos)
				}
			} else {
				if imagewh.Height > maxWHs[i].Height {
					isResize = true
					newImage = imaging.Resize(newImage, 0, maxWHs[i].Height, imaging.Lanczos)
				}
			}

			if isPrint && isResize {
				imagewh, err = DecodeImageWidthHeight(newImage, format)
				if err != nil {
					if isPrint {
						fmt.Println("DecodeImageWidthHeight failed:", err)
					}
					return newImagePath, err
				}
				fmt.Println("Image resize is:", imagewh.Width, "x", imagewh.Height)
			}
		}
		sizeNamei++
		if !isResize {
			break
		}
		isRformat := false
		// 保存图片
		for i, v := range formats {
			if rformat == v || ((rformat == "tiff" || rformat == "tif") && (v == "tiff" || v == "tif")) {
				isRformat = true
			}
			path, err = saveImage(newImage, newPath, imgSize, v, quality)
			if err != nil {
				if isPrint {
					fmt.Println("saveImage failed:", err)
				}
				return newImagePath, err
			}
			newImagePath = append(newImagePath, path)
			fmt.Println(">>>>>>>>>>>>>")
			fmt.Println("i:", i, "v:", v, "rformat:", rformat, "isRformat:", isRformat)
			if i+1 == len(formats) && !isRformat {
				path, err = saveImage(newImage, newPath, imgSize, rformat, quality)
				if err != nil {
					if isPrint {
						fmt.Println("saveImage failed:", err)
					}
					return newImagePath, err
				}
				newImagePath = append(newImagePath, path)
			}
		}
	}
	return newImagePath, nil
}
