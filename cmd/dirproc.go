package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func DirProcess(info_output os.FileInfo, output_size, origin_size *int64, args []string, cp CompressionParameter) {

	var wg sync.WaitGroup
	var (
		err  error
		file *os.File
	)

	var maxworkers int = int(math.Ceil(float64(cpulimit) / 100 * float64(runtime.NumCPU())))

	compressions := make(chan struct{}, maxworkers)
	sizeChan := make(chan int64)

	go func() {
		for s := range sizeChan {
			*output_size += s
		}
	}()

	dir_entries, err := os.ReadDir(args[0])
	for _, e := range dir_entries {
		if e.IsDir() {
			continue
		}
		i, _ := e.Info()
		*origin_size += i.Size()
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
			filepathtocompress := filepath.Join(args[0], e.Name())
			file, err = os.Open(filepathtocompress)

			if err != nil {
				fmt.Println(err)
				continue
			}
			var dst *os.File
			if info_output.IsDir() {
				baseName := filepath.Base(e.Name())
				nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

				targetExt := format
				if targetExt == "" {
					targetExt = "webp"
				}

				finalName := fmt.Sprintf("%s.%s", nameWithoutExt, targetExt)

				finalPath := filepath.Join(outputpath, finalName)

				dst, err = os.Create(finalPath)
				if err != nil {

					fmt.Println(err)
					continue
				}
			}
			compressions <- struct{}{}
			wg.Add(1)
			go func(f, d *os.File) {

				if cp.Watermark != "" {
					MarkFile, errMFile := os.Open(cp.Watermark)
					if errMFile != nil {
						fmt.Println("Err,", errMFile)
						return
					}
					StampIMG(MarkFile, f, d, cp)
					CompressFile(d, d, cp)
				} else {
					CompressFile(f, d, cp)
				}

				s, _ := d.Stat()
				sizeChan <- s.Size()
				f.Close()
				d.Close()
				<-compressions
				wg.Done()
			}(file, dst)

		}

	}
	wg.Wait()
}
