/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	outputpath string
	format     string
	watermark  string
	license    string
	quality    int
	cpulimit   int
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Run the compressing of the file or the directory",
	Long:  `Test long process for compressing`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("process called")
		start := time.Now()
		var (
			origin_size, output_size int64
			file                     *os.File
			file_info                os.FileInfo
			dir_entries              []os.DirEntry
			err                      error
			// newimg      image.Image
		)
		var wg sync.WaitGroup
		info_output, err := os.Stat(outputpath)
		if len(args) >= 1 {
			file_info, err = os.Stat(args[0])
			if err != nil {
				fmt.Println("err: ", err)
				return
			}

			if file_info.IsDir() {
				var maxworkers int = int(math.Ceil(float64(cpulimit) / 100 * float64(runtime.NumCPU())))

				compressions := make(chan struct{}, maxworkers)
				sizeChan := make(chan int64)
				go func() {
					for s := range sizeChan {
						output_size += s
					}
				}()
				dir_entries, err = os.ReadDir(args[0])
				for _, e := range dir_entries {
					if e.IsDir() {
						continue
					}
					i, _ := e.Info()
					origin_size += i.Size()
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

							CompressFile(f, d, CompressionParameter{
								Quality: quality,
								Format:  format,
								License: license,
							})
							s, _ := d.Stat()
							sizeChan <- s.Size()
							f.Close()
							d.Close()
							<-compressions
							wg.Done()
						}(file, dst)

					}
				}

			}

		}

		wg.Wait()
		// À la fin de ton Run
		duration := time.Since(start)

		fmt.Printf("\n✨ Finish in %v !", duration.Round(time.Millisecond))
		fmt.Printf("\n📂 Input size:  %s", FormatSize(origin_size))
		fmt.Printf("\n📦 Output size: %s", FormatSize(output_size))

		// Petit bonus : afficher le gain réel
		if origin_size > 0 {
			gain := float64(origin_size-output_size) / float64(origin_size) * 100
			fmt.Printf("\n🚀 Efficiency:  %.1f%% saved\n", gain)
		}
		if watermark != "" {

		}

	},
}

func init() {
	rootCmd.AddCommand(processCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// processCmd.PersistentFlags().String("foo", "", "A help for foo")
	processCmd.Flags().StringVarP(&outputpath, "output", "o", "", "Path of the output")
	processCmd.Flags().StringVarP(&format, "format", "f", "", "Name the format, (i.e: jpeg, webp, png)")
	processCmd.Flags().StringVarP(&watermark, "mark", "m", "", "Path to the watermark logo")
	processCmd.Flags().IntVarP(&quality, "quality", "q", 80, "Quality of compression (1-100)")
	processCmd.Flags().StringVarP(&license, "license", "l", "", "The license used to protect the file")
	processCmd.Flags().IntVarP(&cpulimit, "cpu-limit", "c", 50, "The percentage of the cpu power")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
