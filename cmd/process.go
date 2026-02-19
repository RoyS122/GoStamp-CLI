/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	outputpath string
	format     string
	watermark  string
	markalpha  uint8
	license    string
	quality    int
	cpulimit   int
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Run the compressing of the file or the directory",
	Long:  `Test long process for compressing`,
	Run: func(cmd *cobra.Command, args []string) {
		var origin_size, output_size int64
		fmt.Println("process called")
		start := time.Now()

		var cp CompressionParameter = CompressionParameter{
			Watermark: watermark,
			MarkAlpha: markalpha,
			Quality:   quality,
			Format:    format,
			License:   license,
		}
		info_output, err := os.Stat(outputpath)
		if err != nil {
			fmt.Println("Error,", err)
		}
		if len(args) >= 1 {
			file_info, err := os.Stat(args[0])
			if err != nil {
				fmt.Println("err: ", err)
				return
			}

			if file_info.IsDir() {
				DirProcess(info_output, &output_size, &origin_size, args, cp)
			} else {
				UniqueProcess(file_info, outputpath, &output_size, &origin_size, args, cp)
			}

		}

		duration := time.Since(start)

		fmt.Printf("\n✨ Finish in %v !", duration.Round(time.Millisecond))
		fmt.Printf("\n📂 Input size:  %s", FormatSize(origin_size))
		fmt.Printf("\n📦 Output size: %s", FormatSize(output_size))

		// Petit bonus : afficher le gain réel
		if origin_size > 0 {
			gain := float64(origin_size-output_size) / float64(origin_size) * 100
			fmt.Printf("\n🚀 Efficiency:  %.1f%% saved\n", gain)
		}

	},
}

func init() {
	rootCmd.AddCommand(processCmd)

	processCmd.Flags().StringVarP(&outputpath, "output", "o", "", "Path of the output")
	processCmd.Flags().StringVarP(&format, "format", "f", "", "Name the format, (i.e: jpeg, webp, png)")
	processCmd.Flags().StringVarP(&watermark, "mark", "m", "", "Path to the watermark logo")
	processCmd.Flags().IntVarP(&quality, "quality", "q", 80, "Quality of compression (1-100)")
	processCmd.Flags().StringVarP(&license, "license", "l", "", "The license used to protect the file")
	processCmd.Flags().IntVarP(&cpulimit, "cpu-limit", "c", 50, "The percentage of the cpu power (1-100)")
	processCmd.Flags().Uint8VarP(&markalpha, "watermark-alpha", "a", uint8(75), "The opacity of the watermark (0-255)")

}
