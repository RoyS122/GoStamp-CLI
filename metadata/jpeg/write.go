package jpeg

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	TagSize       = 2
	SegLengthSize = 2
)

func PushMetadata(dst_path, field, content string) error { // maybe later add some choice of format (XMP OK ! )
	file, err := os.ReadFile(dst_path)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if string(file[:TagSize]) != string([]byte{0xFF, 0xD8}) {
		fmt.Println("This is not a jpeg file")
		return fmt.Errorf("This is not a jpeg file")
	}

	var xmp_identifier []byte = []byte("http://ns.adobe.com/xap/1.0/\x00")

	var content_arr []byte = []byte(fmt.Sprintf(`
	<x:xmpmeta xmlns:x='adobe:ns:meta/'>
		<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
			<rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/">
				<dc:%s> %s </dc:%s>
			</rdf:Description>
		</rdf:RDF>
	</x:xmpmeta>
`, field, content, field))

	var fullsegementlength uint32 = uint32(TagSize + len(content_arr) + SegLengthSize + len(xmp_identifier))

	var segment []byte = make([]byte, 0, fullsegementlength)
	segment = append(segment, 0xFF, 0xE1) // Add AppSegment E1( Exif or xmp)

	var lenBuff []byte = make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuff, uint16(len(content_arr)+len(xmp_identifier)+SegLengthSize))
	segment = append(segment, lenBuff...)

	segment = append(segment, xmp_identifier...) // Add XMP HEADER
	segment = append(segment, content_arr...)    // Add content

	var newf []byte = make([]byte, 0, TagSize+uint32(len(file)+len(segment)))

	newf = append(newf, 0xFF, 0xD8) // Add SOI

	newf = append(newf, segment...) // Add the segment

	newf = append(newf, file[TagSize:]...) // Append the following content of the file

	var c_pointer uint32 = TagSize
	var marker uint16 // maybe usefull
	for c_pointer+TagSize < uint32(len(file)) {
		if file[c_pointer] != 0xFF {
			return fmt.Errorf("Lost in the file at byte: %d, last marker: %x", c_pointer, marker)
		}

		marker = binary.BigEndian.Uint16(file[c_pointer : c_pointer+TagSize])

		switch marker {
		case 0xFFD0, 0xFFD1, 0xFFD2, 0xFFD3, 0xFFD4, 0xFFD5, 0xFFD6, 0xFFD7, 0xFF01: // no length tags
			c_pointer += 2
		case 0xFFDA: // end of metadata
			err := os.WriteFile(dst_path, newf, 0644)
			return err

		default:
			seg_length := binary.BigEndian.Uint16(file[c_pointer+TagSize : c_pointer+TagSize+SegLengthSize])
			c_pointer += uint32(TagSize + seg_length)
		}

	}

	return fmt.Errorf("Reach the end of the file before the marker")
}
