package webp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	wbp_off "golang.org/x/image/webp"
)

const (
	ChunkLengthSize  = 4
	ChunkTypeSize    = 4
	ChunkHeaderSize  = ChunkLengthSize + ChunkTypeSize
	ChunkCRCSize     = 4                   // Maybe not on webp format
	FileHeaderLength = 4 + ChunkHeaderSize // RIFF Code 4 Byte + Size of the file (uint32) except 8 first bytes for RIFF definitions 4 Byte + WEBP code 4 Byte
)

func PushMetadata(dst_path, field, content string) error { // Add metadata on a png file (as string)

	file, err := os.ReadFile(dst_path)
	if err != nil {
		return err
	}

	if len(file) < FileHeaderLength ||
		!bytes.Equal(file[:ChunkTypeSize], []byte{'R', 'I', 'F', 'F'}) ||
		!bytes.Equal(file[ChunkTypeSize+ChunkLengthSize:FileHeaderLength], []byte{'W', 'E', 'B', 'P'}) {
		return fmt.Errorf("not a valid WebP file")
	}
	f_io, err := os.Open(dst_path)
	if err != nil {
		return err
	}
	cfg, err := wbp_off.DecodeConfig(f_io)
	if err != nil {
		return err
	}
	f_io.Close()

	var chunk_type []byte = []byte{'X', 'M', 'P', ' '}

	var xmp_content []byte = []byte(fmt.Sprintf(`
	<x:xmpmeta xmlns:x='adobe:ns:meta/'>
		<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
			<rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/">
				<dc:%s> %s </dc:%s>
			</rdf:Description>
		</rdf:RDF>
	</x:xmpmeta>
`, field, content, field))
	var chunk_length uint32 = uint32(len(xmp_content))
	if len(xmp_content)%2 != 0 {
		xmp_content = append(xmp_content, uint8(0))
	}

	var chunk []byte = make([]byte, 0, ChunkHeaderSize+chunk_length)

	chunk = append(chunk, chunk_type...)
	chunk = binary.LittleEndian.AppendUint32(chunk, chunk_length)
	chunk = append(chunk, xmp_content...)

	// var fullchunklength uint32 = uint32(len(chunk))
	var VP8X_Chunk []byte
	if !(bytes.Contains([]byte("VP8X"), file) || bytes.Contains(file, []byte("VP8X"))) { // file does not have a VP8X
		VP8X_Chunk = append(VP8X_Chunk, []byte{'V', 'P', '8', 'X'}...)
		VP8X_Chunk = binary.LittleEndian.AppendUint32(VP8X_Chunk, 10) // 10 The size of the structure requiered (1 byte for flagging the features 3 Reserved + 6 sizing)
		VP8X_Chunk = append(VP8X_Chunk, byte(0b00010000))
		VP8X_Chunk = append(VP8X_Chunk, make([]byte, 3)...)
		bufW := make([]byte, 4)
		bufH := make([]byte, 4)

		binary.LittleEndian.PutUint32(bufH, uint32(cfg.Height-1))
		binary.LittleEndian.PutUint32(bufW, uint32(cfg.Width-1))

		VP8X_Chunk = append(VP8X_Chunk, bufW[:3]...)
		VP8X_Chunk = append(VP8X_Chunk, bufH[:3]...)

	}

	var newFile []byte = make([]byte, 0, uint32(len(file))+ChunkHeaderSize+chunk_length)
	if len(VP8X_Chunk) != 0 {
		newFile = append(newFile, file[:FileHeaderLength]...)
		newFile = append(newFile, VP8X_Chunk...)
		newFile = append(newFile, file[FileHeaderLength:]...)
	} else {
		newFile = append(newFile, file...)
	}

	newFile = append(newFile, chunk...)

	binary.LittleEndian.PutUint32(newFile[ChunkTypeSize:ChunkHeaderSize], uint32(len(newFile)-ChunkHeaderSize))

	if err := os.WriteFile(dst_path, newFile, 0644); err != nil {
		return err
	}

	return nil
} // It's work :)
