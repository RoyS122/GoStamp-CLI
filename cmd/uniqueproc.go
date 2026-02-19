package cmd

import (
	"fmt"
	"os"
)

func UniqueProcess(file_info os.FileInfo, outputpath string, output_size, origin_size *int64, args []string, cp CompressionParameter) {
	*origin_size += file_info.Size()

	dst, err := os.Create(outputpath)
	if err != nil {
		fmt.Println("Error,", err)
		return
	}

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Println("Error,", err)
		return
	}

	if cp.Watermark != "" {
		MarkFile, errMFile := os.Open(cp.Watermark)
		if errMFile != nil {
			fmt.Println("Err,", errMFile)
			return
		}
		StampIMG(MarkFile, file, dst, cp)
		CompressFile(dst, dst, cp)
	} else {
		CompressFile(file, dst, cp)
	}

	s, _ := dst.Stat()
	*output_size += s.Size()

	file.Close()
	dst.Close()

}
