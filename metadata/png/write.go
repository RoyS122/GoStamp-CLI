package png

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
)

var pngSinature []byte = []byte{137, 80, 78, 71, 13, 10, 26, 10}

const (
	ChunkLengthSize = 4
	ChunkTypeSize   = 4
	ChunkHeaderSize = ChunkLengthSize + ChunkTypeSize
	ChunkCRCSize    = 4

	PNGMinimumPayload = ChunkHeaderSize + ChunkCRCSize
	PNGSigSize        = 8
)

func PushMetadata(dst_path, field, content string) error { // Add metadata on a png file (as string)

	file, err := os.ReadFile(dst_path)
	if err != nil {
		return err
	}

	if len(file) < 8 || !bytes.Equal(file[:8], pngSinature) {
		return fmt.Errorf("not a valid PNG file")
	}

	// It's a reel png

	var (
		content_type  []byte = []byte("tEXt")
		chunk_content []byte = []byte(field + "\x00" + content)

		content_size int = len(chunk_content)

		crcBuffer []byte = append([]byte("tEXt"), chunk_content...)

		crcHash uint32 = crc32.ChecksumIEEE(crcBuffer) // CRC (Check of the integrity of datas)

		chunk []byte = make([]byte, ChunkHeaderSize+len(chunk_content)+ChunkCRCSize) // Chunk array
	)

	binary.BigEndian.PutUint32(chunk[:4], uint32(content_size))               // Put the content size at his part of chunk
	copy(chunk[ChunkLengthSize:ChunkHeaderSize], content_type)                // Put the content type in
	copy(chunk[ChunkHeaderSize:ChunkHeaderSize+content_size], chunk_content)  // Put The content
	binary.BigEndian.PutUint32(chunk[ChunkHeaderSize+content_size:], crcHash) // Put the hash

	// Here we will run through the file to get the append point of the file

	var c_pointer uint32 = PNGSigSize
	fileSize := uint32(len(file))
	for {
		if c_pointer+PNGMinimumPayload > fileSize {
			return fmt.Errorf("corrupted PNG structure")
		}

		currentType := string(file[c_pointer+ChunkLengthSize : c_pointer+ChunkHeaderSize])

		if currentType == "IEND" {
			fmt.Println("IEND found")
			break
		}

		sizeofchunk := binary.BigEndian.Uint32(file[c_pointer : c_pointer+ChunkLengthSize])
		c_pointer += PNGMinimumPayload + sizeofchunk
	}

	newFile := make([]byte, 0, len(file)+len(chunk))
	newFile = append(newFile, file[:c_pointer]...)
	newFile = append(newFile, chunk...)
	newFile = append(newFile, file[c_pointer:]...)

	if err := os.WriteFile(dst_path, newFile, 0644); err != nil {
		return err
	}

	fmt.Println("Just pushed : ", field, content)

	return nil
}
