package cmd

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/chai2010/webp"
)

type CompressionParameter struct {
	Watermark string
	MarkAlpha uint8
	Quality   int
	Format    string
	License   string
}

func CompressFile(origin, dst *os.File, params CompressionParameter) error {

	img, _, err := image.Decode(origin)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}

	q := params.Quality
	if q == 0 {
		q = 80
	}
	if q > 100 {
		q = 100
	}

	switch strings.ToLower(params.Format) {
	case "webp":
		err = webp.Encode(dst, img, &webp.Options{
			Lossless: false,
			Quality:  float32(q),
		})
	case "png":

		enc := png.Encoder{CompressionLevel: png.BestCompression}

		err = enc.Encode(dst, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(dst, img, &jpeg.Options{Quality: int(q)})
	default:
		err = webp.Encode(dst, img, &webp.Options{Quality: float32(q)})
	}

	if err != nil {
		return fmt.Errorf("erreur encodage destination: %w", err)
	}

	return nil
}

func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
