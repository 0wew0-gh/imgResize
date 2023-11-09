package mediaResize

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// MediaWH image width and height
type MediaWH struct {
	Width  int `json:"width"`  //宽
	Height int `json:"height"` //高
}

type ProbeData struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

// ========================
//
//	缩放并压缩媒体文件
//	path		string		媒体文件路径
//	fileType	string		媒体文件类型
//	返回值		*MediaWH	媒体文件宽高
//	返回值		error		错误信息
func DecodeFileWidthHeight(path string, fileType string) (*MediaWH, error) {
	Mediatypes := strings.Split(strings.ToLower(fileType), "/")
	fType := "image"
	extType := "jpg"
	if len(Mediatypes) > 1 {
		fType = Mediatypes[0]
		extType = Mediatypes[1]
	}
	fmt.Println("fType:", fType)
	fmt.Println("extType:", extType)
	switch fType {
	case "video":
		cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", path)

		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		var data ProbeData
		err = json.Unmarshal(output, &data)
		if err != nil {
			return nil, err
		}

		for _, stream := range data.Streams {
			fmt.Printf("Width: %d, Height: %d\n", stream.Width, stream.Height)
			return &MediaWH{
				Width:  stream.Width,
				Height: stream.Height,
			}, nil
		}
	case "image":
		img, err := imaging.Open(path)
		if err != nil {
			return nil, err
		}
		imgBytes, err := imageToBytes(img)
		if err != nil {
			return nil, err
		}
		var imgConf image.Config
		switch extType {
		case "jpg", "jpeg":
			imgConf, err = jpeg.DecodeConfig(bytes.NewReader(imgBytes))
		case "webp":
			imgConf, err = webp.DecodeConfig(bytes.NewReader(imgBytes))
		case "png":
			imgConf, err = png.DecodeConfig(bytes.NewReader(imgBytes))
		case "tif", "tiff":
			imgConf, err = tiff.DecodeConfig(bytes.NewReader(imgBytes))
		case "gif":
			imgConf, err = gif.DecodeConfig(bytes.NewReader(imgBytes))
		case "bmp":
			imgConf, err = bmp.DecodeConfig(bytes.NewReader(imgBytes))
		default:
			return nil, errors.New("unknown file type")
		}
		if err != nil {
			return nil, err
		}
		return &MediaWH{
			Width:  imgConf.Width,
			Height: imgConf.Height,
		}, nil
	}
	return nil, errors.New("unknown file type")
}

// ========================
//
//	使用`[]byte`解析图片的宽高信息
//	imgBytes	[]byte		图片字节
//	fileType	string		图片格式
//	返回值		*MediaWH	图片宽高
//	返回值		error		错误信息
func DecodeBytesWidthHeight(imgBytes []byte, fileType string) (*MediaWH, error) {
	var (
		imgConf image.Config
		err     error
	)
	switch strings.ToLower(fileType) {
	case "jpg", "jpeg":
		imgConf, err = jpeg.DecodeConfig(bytes.NewReader(imgBytes))
	case "webp":
		imgConf, err = webp.DecodeConfig(bytes.NewReader(imgBytes))
	case "png":
		imgConf, err = png.DecodeConfig(bytes.NewReader(imgBytes))
	case "tif", "tiff":
		imgConf, err = tiff.DecodeConfig(bytes.NewReader(imgBytes))
	case "gif":
		imgConf, err = gif.DecodeConfig(bytes.NewReader(imgBytes))
	case "bmp":
		imgConf, err = bmp.DecodeConfig(bytes.NewReader(imgBytes))
	default:
		return nil, errors.New("unknown file type")
	}
	if err != nil {
		return nil, err
	}
	return &MediaWH{
		Width:  imgConf.Width,
		Height: imgConf.Height,
	}, nil
}

// ========================
//
//	使用`image.Image`解析图片的宽高信息
//	img			image.Image	图片
//	fileType	string		图片格式
//	返回值		*MediaWH	图片宽高
//	返回值		error		错误信息
func DecodeImageWidthHeight(img image.Image, fileType string) (*MediaWH, error) {
	var (
		imgBytes []byte
		err      error
	)
	imgBytes, err = imageToBytes(img)
	if err != nil {
		return nil, err
	}
	return DecodeBytesWidthHeight(imgBytes, fileType)
}

func imageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func saveImage(img image.Image, path string, imgSize string, imgType string, opts int) (string, error) {
	var (
		dst *os.File
		err error
	)

	pathList := strings.Split(path, ".")
	pathList[len(pathList)-1] = imgType
	path = strings.Join(pathList, ".")
	path = strings.Replace(path, ".", "."+imgSize+".", -1)

	dst, err = os.Create(path)
	if err != nil {
		dst.Close()
		return path, err
	}

	switch strings.ToLower(imgType) {
	case "jpg", "jpeg":
		if opts > 100 {
			err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 100})
		} else if opts < 0 {
			err = jpeg.Encode(dst, img, nil)
		} else {
			err = jpeg.Encode(dst, img, &jpeg.Options{Quality: opts})
		}
	case "webp":
		if opts > 100 {
			err = webp.Encode(dst, img, &webp.Options{Lossless: true, Quality: 100})
		} else if opts < 0 {
			err = webp.Encode(dst, img, &webp.Options{Lossless: true})
		} else {
			err = webp.Encode(dst, img, &webp.Options{Lossless: true, Quality: float32(opts)})
		}
	case "png":
		err = png.Encode(dst, img)
	case "tif", "tiff":
		err = tiff.Encode(dst, img, nil)
	case "gif":
		err = gif.Encode(dst, img, nil)
	case "bmp":
		err = bmp.Encode(dst, img)
	default:
		dst.Close()
		return path, errors.New("unknown file type")
	}
	dst.Close()
	if err != nil {
		return path, err
	}
	return path, nil
}

// ========================
//
//	根据文件地址获取文件类型
//	path		string		文件地址
//	返回值		string		文件类型
//	返回值		error		错误信息
func DetectContentType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return "", fmt.Errorf("os.Open failed: %s", err)
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return "", fmt.Errorf("file.Stat failed: %s", err)
	}
	buffer := make([]byte, info.Size())
	_, err = file.Read(buffer)
	file.Close()
	if err != nil {
		return "", fmt.Errorf("file.Read failed: %s", err)
	}
	return http.DetectContentType(buffer), nil
}

// ========================
//
//	缩放并压缩媒体文件
//	path		string		原媒体文件路径
//	newPath		string		新媒体文件路径
//	contentType	string		媒体文件类型
//	codeRate	int		视频码率,-1为默认值:1500k
//	width		int		缩放宽度
//	height		int		缩放高度
//	返回值		image.Image	新媒体文件
//	返回值		error		错误信息
func Resize(path string, newPath string, contentType string, codeRate int, width int, height int) (image.Image, error) {
	Mediatypes := strings.Split(strings.ToLower(contentType), "/")
	fType := "image"
	if len(Mediatypes) > 1 {
		fType = Mediatypes[0]
	}
	switch fType {
	case "image":
		img, err := imaging.Open(path)
		if err != nil {
			return nil, err
		}
		newImage := imaging.Resize(img, width, height, imaging.Lanczos)
		return newImage, nil
	case "video":
		if height%2 != 0 {
			height++
		}
		scale := fmt.Sprintf("%dx%d", width, height)
		cRate := fmt.Sprintf("%dk", codeRate)
		cmd := exec.Command("ffmpeg", "-i", path, "-b:v", cRate, "-s", scale, "-acodec", "copy", newPath)
		if codeRate < 0 {
			cmd = exec.Command("ffmpeg", "-i", path, "-b:v", "1500k", "-s", scale, "-acodec", "copy", newPath)
		}
		if _, err := os.Stat(newPath); !os.IsNotExist(err) {
			fmt.Println("file exist:", newPath, "break")

			file, err := os.Open(newPath)
			if err != nil {
				file.Close()
				return nil, err
			}
			info, err := file.Stat()
			file.Close()
			if err != nil {
				return nil, err
			}
			fmt.Println("file:", newPath, "size:", info.Size())
			if info.Size() > 0 {
				return nil, nil
			} else {
				err = os.Remove(newPath)
				if err != nil {
					return nil, err
				}
				// cmd = exec.Command("ffmpeg", "-y", "-i", path,"-vf", "scale=-2:ih", "-b:v", cRate, "-s", scale, "-acodec", "copy", newPath)
			}
		}
		fmt.Println("Resize scale:", scale)
		fmt.Println("path:", path)
		fmt.Println("newPath:", newPath)
		fmt.Println("contentType:", contentType)
		fmt.Println("codeRate:", cRate)

		// err := cmd.Run()
		// if err != nil {
		// 	return nil, err
		// }
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(string(output))
			return nil, err
		}

	}
	return nil, nil
}

// ========================
//
//	根据宽高比例缩放媒体文件
//	mediaWidth	int		原媒体文件宽度
//	mediaHeght	int		原媒体文件高度
//	width		int		缩放宽度
//	height		int		缩放高度
//	返回值		int		新媒体文件宽度
//	返回值		int		新媒体文件高度
func calcResolutionRatio(mediaWidth int, mediaHeght int, width int, height int) (int, int) {
	var (
		newWidth  float64 = 0.0
		newHeight float64 = 0.0
	)
	if width < 0 || height < 0 {
		return mediaWidth, mediaHeght
	}
	if mediaWidth > mediaHeght {
		newWidth = float64(width)
		newHeight = float64(mediaHeght) * newWidth / float64(mediaWidth)
	} else {
		newHeight = float64(height)
		newWidth = float64(mediaWidth) * newHeight / float64(mediaHeght)
	}
	return int(newWidth), int(newHeight)
}
