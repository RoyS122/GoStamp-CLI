package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/chai2010/webp"
)

func StampIMG(buffer, src, dst *os.File, params CompressionParameter) {

	var (
		src_img, buffer_img image.Image
		dst_drawable        draw.Image
		buffer_rect         image.Rectangle
		buffer_pos          image.Point
		opacity_mask        *image.Uniform
		err                 error
		n                   string
	)

	src_img, n, err = image.Decode(src)
	if err != nil {
		fmt.Println("Error,", err, n)
		return
	}

	dst_drawable = image.NewRGBA(src_img.Bounds())
	draw.Draw(dst_drawable, src_img.Bounds(), src_img, image.Point{}, draw.Src)

	buffer_img, n, err = image.Decode(buffer)
	if err != nil {
		fmt.Println("Error,", err, "name", n)
	}

	buffer_pos = image.Pt(int(dst_drawable.Bounds().Max.X-buffer_img.Bounds().Min.X/2), int(dst_drawable.Bounds().Max.Y-buffer_img.Bounds().Min.Y/2))
	buffer_rect = buffer_img.Bounds()
	buffer_rect.Add(buffer_pos)

	opacity_mask = image.NewUniform(color.Alpha{params.MarkAlpha})

	draw.DrawMask(dst_drawable, buffer_rect, buffer_img, image.Point{}, opacity_mask, image.Point{}, draw.Over)
	switch params.Format {
	case "jpg", "jpeg":
		jpeg.Encode(dst, dst_drawable, &jpeg.Options{Quality: 100})
	case "png":
		png.Encode(dst, dst_drawable)
	case "webp":
		webp.Encode(dst, dst_drawable, &webp.Options{Lossless: true})
	}

}
