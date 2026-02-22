# Jpeg data structure

A jpeg file is build with differents segments identified as tags of four bytes "FF xx" 
The segment with a content is writed as four byte "FF xx" then two bytes of length data (the length count the bytes of lenght ex: FF E1 00 06 74 65 73 74 ("test" data, 2 bytes of definition, 2 bytes of length and 4 bytes of content)) 

## Segments

- `FF D8`: Start of the file
- `FF E0`->`FF EF`: AppSegments ( E1 for EXIF datas )
- `FF FE`: Comments
- `FF D9`: End of file 

## Exif format

Is like an indexed db

### TIFF Bloc Format


#### IFD


## XMP
```xml
<x:xmpmeta xmlns:x='adobe:ns:meta/'>
	<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
		<rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/">
			<dc:field1name> Field content </dc:field1name>
			<dc:field2name> Field 2 content </dc:field2name>
		</rdf:Description>
	</rdf:RDF>
...
</x:xmpmeta>
```