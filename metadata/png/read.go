package png

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gostampcli/metadata/utils"
	"os"
)

func ReadMetadata(path string) (r []utils.MetaField) {

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if !bytes.Equal(file[:PNGSigSize], pngSinature) {
		fmt.Println("Not a png file")
		return nil
	}

	var c uint32 = PNGSigSize

	var currentType string = string(file[c+ChunkLengthSize : c+ChunkHeaderSize])

	for currentType != "IEND" {
		if c+PNGMinimumPayload > uint32(len(file)) {
			fmt.Println("Corrupted png format ")
			return nil
		}
		chunkLength := binary.BigEndian.Uint32(file[c : c+ChunkLengthSize])

		switch currentType {
		case "iTXt", "tEXt", "zTXt":

			var value []byte = file[c+ChunkHeaderSize : c+ChunkHeaderSize+chunkLength]
			var vSplitted []string
			var vs []byte
			for _, b := range value {
				if b == 0x00 {
					vSplitted = append(vSplitted, string(vs))
					vs = []byte{}
					continue
				}
				vs = append(vs, b)
			}
			vSplitted = append(vSplitted, string(vs))

			r = append(r, utils.MetaField{Title: vSplitted[0], Content: vSplitted[1]})

		}

		c += ChunkHeaderSize + chunkLength + ChunkCRCSize
		currentType = string(file[c+ChunkLengthSize : c+ChunkHeaderSize])
	}

	return
}
