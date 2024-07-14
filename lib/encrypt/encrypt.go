package encrypt

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"

	"github.com/jaimejcs/go_esteganografia/lib/commonFunc"
)

// EncodeNRGBA encodes a given string into the input image using least significant bit encryption (LSB steganography)
// The minnimum image size is 24 pixels for one byte. For each additional byte, it is necessary 3 more pixels.
/*
	Input:
		writeBuffer *bytes.Buffer : the destination of the encoded image bytes
		pictureInputFile image.NRGBA : image data used in encoding
		message []byte : byte slice of the message to be encoded
	Output:
		bytes buffer ( io.writter ) to create file, or send data.
*/
func EncodeNRGBA(writeBuffer *bytes.Buffer, rgbImage *image.NRGBA, message []byte) error {

	var messageLength = uint32(len(message))

	var width = rgbImage.Bounds().Dx()
	var height = rgbImage.Bounds().Dy()
	var c color.NRGBA
	var bit byte
	var ok bool
	//var encodedImage image.Image
	if commonFunc.MaxEncodeSize(rgbImage) < messageLength+4 {
		return errors.New("message too large for image")
	}

	one, two, three, four := commonFunc.splitToBytes(messageLength)

	message = append([]byte{four}, message...)
	message = append([]byte{three}, message...)
	message = append([]byte{two}, message...)
	message = append([]byte{one}, message...)

	ch := make(chan byte, 100)

	go commonFunc.getNextBitFromString(message, ch)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbImage.NRGBAAt(x, y) // get the color at this pixel

			/*  RED  */
			bit, ok = <-ch
			if !ok { // if we don't have any more bits left in our message

				rgbImage.SetNRGBA(x, y, c)
				break
			}
			commonFunc.setLSB(&c.R, bit)

			/*  GREEN  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetNRGBA(x, y, c)
				break
			}
			commonFunc.setLSB(&c.G, bit)

			/*  BLUE  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetNRGBA(x, y, c)
				break
			}
			commonFunc.setLSB(&c.B, bit)

			rgbImage.SetNRGBA(x, y, c)
		}
	}

	err := png.Encode(writeBuffer, rgbImage)
	return err
}

// Encode encodes a given string into the input image using least significant bit encryption (LSB steganography)
// The minnimum image size is 23 pixels
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		writeBuffer *bytes.Buffer : the destination of the encoded image bytes
		message []byte : byte slice of the message to be encoded
		pictureInputFile image.Image : image data used in encoding
	Output:
		bytes buffer ( io.writter ) to create file, or send data.
*/
func Encode(writeBuffer *bytes.Buffer, pictureInputFile image.Image, message []byte) error {

	rgbImage := commonFunc.imageToNRGBA(pictureInputFile)

	return EncodeNRGBA(writeBuffer, rgbImage, message)

}
