package imgResize

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/chai2010/webp"
)

// ImageWH image width and height
type ImageWH struct {
	Width  int `json:"width"`  //宽
	Height int `json:"height"` //高
}

// ========================
//
//	使用`[]byte`解析图片的宽高信息
//	imgBytes	[]byte		图片字节
//	fileType	string		图片格式
//	返回值		*ImageWH	图片宽高
//	返回值		error		错误信息
func DecodeBytesWidthHeight(imgBytes []byte, fileType string) (*ImageWH, error) {
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
	return &ImageWH{
		Width:  imgConf.Width,
		Height: imgConf.Height,
	}, nil
}

// ========================
//
//	使用`image.Image`解析图片的宽高信息
//	img			image.Image	图片
//	fileType	string		图片格式
//	返回值		*ImageWH	图片宽高
//	返回值		error		错误信息
func DecodeImageWidthHeight(img image.Image, fileType string) (*ImageWH, error) {
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
