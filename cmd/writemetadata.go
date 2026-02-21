package cmd

import (
	"fmt"
	metapng "gostampcli/metadata/png"
)

type MetadataParameter struct {
	Title        string
	Author       string
	License      string
	Date         string
	OriginalDate string
}

func AddMetadata(dst string, mparams MetadataParameter) {

	if err := metapng.PushMetadata(dst, "title", mparams.Title); err != nil {
		fmt.Println("err,", err)
	}

}
