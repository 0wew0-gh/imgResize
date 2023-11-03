package imgResize

import (
	"fmt"
	"testing"
)

func TestImgResize(t *testing.T) {
	imageList := []string{"image/01.png", "image/02.jpg", "image/03.png", "image/04.png"}
	newImgs := []string{"new/001.png", "new/002.jpg", "new/003.png", "new/004.png"}
	maxWHs := []ImageWH{
		{Width: 200, Height: 200},
		{Width: 500, Height: 500},
		{Width: 1000, Height: 1000},
	}

	newImagePath, err := ImgResizes(imageList, newImgs, []string{"jpg", "png", "webp"}, maxWHs, -1)
	if err != nil {
		fmt.Println("ImgResizes failed:", err)
	}
	fmt.Println(newImagePath)
}
