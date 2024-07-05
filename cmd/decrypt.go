package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/spf13/cobra"
)

var decrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Oculta a mensagem dentro da imagem",
	Long:  "Oculta a mensagem dentro dos bits menos significativos do RGB da imagem",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("decrypt executado")
	},
}

// decodeNRGBA gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
/*
	Input:
		startOffset uint32 : number of bytes used to declare size of message
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.NRGBA : image data used in decoding
	Output:
		message []byte decoded from image
*/
func decodeNRGBA(startOffset uint32, msgLen uint32, rgbImage *image.NRGBA) (message []byte) {

	var byteIndex uint32
	var bitIndex uint32

	width := rgbImage.Bounds().Dx()
	height := rgbImage.Bounds().Dy()

	var c color.NRGBA
	var lsb byte

	message = append(message, 0)

	// iterate through every pixel in the image and stitch together the message bit by bit
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbImage.NRGBAAt(x, y) // get the color of the pixel

			/*  RED  */
			lsb = getLSB(c.R)                                                    // get the least significant bit from the red component of this pixel
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb) // add this bit to the message
			bitIndex++

			if bitIndex > 7 { // when we have filled up a byte, move on to the next byte
				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}

			/*  GREEN  */
			lsb = getLSB(c.G)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {

				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}

			/*  BLUE  */
			lsb = getLSB(c.B)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {
				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}
		}
	}
	return
}

// decode gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		startOffset uint32 : number of bytes used to declare size of message
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.Image : image data used in decoding
	Output:
		message []byte decoded from image
*/
func decode(startOffset uint32, msgLen uint32, pictureInputFile image.Image) (message []byte) {

	rgbImage := imageToNRGBA(pictureInputFile)
	return decodeNRGBA(startOffset, msgLen, rgbImage)

}

// Decode gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.Image : image data used in decoding
	Output:
		message []byte decoded from image
*/
func Decode(msgLen uint32, pictureInputFile image.Image) (message []byte) {
	return decode(4, msgLen, pictureInputFile) // the offset of 4 skips the "header" where message length is defined

}

// imageToNRGBA converts image.Image to image.NRGBA
func imageToNRGBA(src image.Image) *image.NRGBA {
	bounds := src.Bounds()

	var m *image.NRGBA
	var width, height int

	width = bounds.Dx()
	height = bounds.Dy()

	m = image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), src, bounds.Min, draw.Src)
	return m
}
