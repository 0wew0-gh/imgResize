package mediaResize

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// ========================
//
//	批量处理图片
//	paths		[]string	原图片路径
//	newPaths	[]string	新图片路径
//	formats		[]string	图片格式
//	maxWHs		[]MediaWH	图片宽高
//	quality		int		图片质量
//	isPrint		bool		是否打印错误及提示信息
//	返回值		[]string	新图片路径
//	返回值		error		错误信息
func VideoResizes(paths []string, newPaths []string, formats []string, maxWHs []MediaWH, quality int, isPrint bool) ([][]string, [][]string, [][]string, error) {

	newImagePath := [][]string{}
	newsizes := [][]string{}
	newformats := [][]string{}

	for i := 0; i < len(paths); i++ {
		newpaths, sizes, formats, err := VideoResize(paths[i], newPaths[i], formats, maxWHs, quality, isPrint)
		if err != nil {
			return newImagePath, newsizes, newformats, err
		}
		newImagePath = append(newImagePath, newpaths)
		newsizes = append(newsizes, sizes)
		newformats = append(newformats, formats)
	}
	return newImagePath, newsizes, newformats, nil
}

// ========================
//
//	处理图片
//	paths		string		原图片路径
//	newPaths	string		新图片路径
//	formats		[]string	图片格式
//	maxWHs		[]MediaWH	图片宽高
//	quality		int		图片质量
//	isPrint		bool		是否打印错误及提示信息
//	返回值		[]string	新图片路径
//	返回值		error		错误信息
func VideoResize(path string, newPath string, formats []string, maxWHs []MediaWH, codeRate int, isPrint bool) ([]string, []string, []string, error) {
	newImagePath := []string{}
	sizes := []string{}
	newformats := []string{}

	file, err := os.Open(path)
	if err != nil {
		// fmt.Println("os.Open failed:", err)
		file.Close()
		return newImagePath, sizes, newformats, err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return newImagePath, sizes, newformats, err
	}

	buffer := make([]byte, info.Size())
	_, err = file.Read(buffer)
	file.Close()
	if err != nil {
		return newImagePath, sizes, newformats, err
	}

	contentType := http.DetectContentType(buffer)
	fmt.Println("contentType:", contentType)

	// Mediatypes := strings.Split(contentType, "/")
	// fExt := "mp4"
	// if len(Mediatypes) > 1 {
	// 	fExt = strings.ToLower(Mediatypes[1])
	// }

	// 解析图片宽高后，进行图片缩放
	videowh, err := DecodeFileWidthHeight(path, contentType)
	if err != nil {
		if isPrint {
			fmt.Println("1 DecodeFileWidthHeight failed:", err)
		}
		return newImagePath, sizes, newformats, err
	}

	exists := map[string]bool{}
	sizeNamei := 0
	for i := 0; i < len(maxWHs); i++ {
		var (
			videoSize string = ""
			w         int    = videowh.Width
			h         int    = videowh.Height
		)

		switch sizeNamei {
		case 0:
			videoSize = "S"
		case 1:
			videoSize = "M"
		case 2:
			videoSize = "L"
		default:
			for ii := 2; ii < i; ii++ {
				videoSize += "X"
			}
			videoSize += "L"
		}
		fmt.Println("videoSize:", videoSize)

		if isPrint && i == 0 {
			fmt.Println("Image size is:", videowh.Width, "x", videowh.Height)
		}
		if maxWHs[i].Width < 0 || maxWHs[i].Height < 0 {
			// 不进行图片缩放
			videoSize = "R"
			sizeNamei--
			sizes = append(sizes, videoSize)
		} else {
			w, h = calcResolutionRatio(videowh.Width, videowh.Height, maxWHs[i].Width, maxWHs[i].Height)
			fmt.Println(">>>>> calcResolutionRatio", w, h)
			sizes = append(sizes, videoSize)
		}

		sizeNamei++

		// isRformat := false

		// 保存图片
		for _, v := range formats {
			// if fExt == v {
			// 	isRformat = true
			// }

			pathList := strings.Split(newPath, ".")
			pathList[len(pathList)-1] = v
			resizePath := strings.Join(pathList, ".")
			resizePath = strings.Replace(resizePath, ".", "."+videoSize+".", -1)

			_, err = Resize(path, resizePath, contentType, codeRate, w, h)
			if err != nil {
				if isPrint {
					fmt.Println("Resize failed:", err)
				}
				return newImagePath, sizes, newformats, err
			}

			if isPrint {
				fmt.Println("saveImage:", path)
			}
			newImagePath = append(newImagePath, resizePath)

			if _, ok := exists[v]; !ok {
				newformats = append(newformats, v)
				exists[v] = true
			}
			// if i+1 == len(formats) && !isRformat {
			// 	pathList[len(pathList)-1] = fExt
			// 	resizePath = strings.Join(pathList, ".")
			// 	resizePath = strings.Replace(resizePath, ".", "."+videoSize+".", -1)
			// 	_, err = Resize(path, resizePath, contentType, codeRate, w, h)
			// 	if err != nil {
			// 		if isPrint {
			// 			fmt.Println("Resize failed:", err)
			// 		}
			// 		return newImagePath, sizes, newformats, err
			// 	}
			// 	if isPrint {
			// 		fmt.Println("saveImage 2:", path)
			// 	}
			// 	newImagePath = append(newImagePath, resizePath)
			// 	if _, ok := exists[fExt]; !ok {
			// 		newformats = append(newformats, fExt)
			// 		exists[fExt] = true
			// 	}
			// }
		}
	}
	return newImagePath, sizes, newformats, nil
}
