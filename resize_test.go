package mediaResize

import (
	"fmt"
	"os"
	"testing"
)

func TestImgResize(t *testing.T) {
	var (
		imageList []string  = []string{"media/01.png", "media/02.jpg", "media/03.png", "media/04.png"}
		newImgs   []string  = []string{"new/001.png", "new/002.jpg", "new/003.png", "new/004.png"}
		maxWHs    []MediaWH = []MediaWH{
			{Width: -1, Height: -1},
			{Width: 200, Height: 200},
			{Width: 500, Height: 500},
			{Width: 1000, Height: 1000},
		}
		newPath string = "new"
		err     error
	)

	os.MkdirAll("new", 0777)
	// 保存文件到临时文件夹
	if _, err = os.Stat(newPath); os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("os.MkdirAll:", err)
			return
		}
	}

	newImagePath, newsizes, newformats, err := ImgResizes(imageList, newImgs, []string{"jpg", "webp"}, maxWHs, -1, true)
	if err != nil {
		fmt.Println("ImgResizes failed:", err)
	}
	fmt.Println(newImagePath)
	fmt.Println(newsizes)
	fmt.Println(newformats)
}

func TestVideoResize(t *testing.T) {

	var (
		imageList []string  = []string{"media/001.mp4"}
		newImgs   []string  = []string{"new/001.mp4"}
		maxWHs    []MediaWH = []MediaWH{
			{Width: -1, Height: -1},
			// {Width: 200, Height: 200},
			// {Width: 500, Height: 500},
			{Width: 1000, Height: 1000},
		}
		newPath string = "new"
		err     error
	)

	os.MkdirAll("new", 0777)
	// 保存文件到临时文件夹
	if _, err = os.Stat(newPath); os.IsNotExist(err) {
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			fmt.Println("os.MkdirAll:", err)
			return
		}
	}

	newImagePath, newsizes, newformats, err := VideoResizes(imageList, newImgs, []string{"mkv"}, maxWHs, 1500, true)
	if err != nil {
		fmt.Println("ImgResizes failed:", err)
	}
	fmt.Println(newImagePath)
	fmt.Println(newsizes)
	fmt.Println(newformats)
}
