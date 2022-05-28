package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
)

var black = color.RGBA{0, 0, 0, 255}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 16, 8))
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			if x%2 == 0 {
				img.Set(x, y, image.White)
			} else {
				img.SetRGBA(x, y, black)
			}
		}
	}

	bytes := convertImageToBits(img)
	for i := 0; i < len(bytes); i++ {
		fmt.Printf("%x\n", bytes[i])
	}
}

func convertImageToBits(img image.Image) []byte {
	wh := img.Bounds()
	b := make([]byte, wh.Max.X*wh.Max.Y/8)
	for y := 0; y < wh.Max.Y; y++ {
		for x := 0; x < wh.Max.X; x++ {
			if img.At(x, y) == black {
				continue
			}
			fmt.Println(img.At(x, y))
			byteIndex := (y * wh.Max.X / 8) + (x / 8)
			bitIndex := x % 8
			fmt.Println(byteIndex, bitIndex, strconv.FormatInt(int64(b[byteIndex]), 2), strconv.FormatInt(1<<bitIndex, 2))
			b[byteIndex] |= (1 << bitIndex)
		}
	}
	return b
}
