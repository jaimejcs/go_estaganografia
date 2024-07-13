package cmd

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var decrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Oculta a mensagem dentro da imagem <path/to/file>",
	Long: `Oculta a mensagem dentro da imagem <path/to/file>
			Exemplo: encrypt imagem.png -m 'Nosso segredo morre aqui' -o imagem_secreta.png`,
	Run: func(cmd *cobra.Command, args []string) {
		var input = args[0]
		var output string

		if output = cmd.Flag("output").Value.String(); len(strings.Trim(output, " ")) == 0 {
			output = "out.png"
		}

		inFile, err := os.Open(input) // Opens input file provided in the flags
		if err != nil {
			panic(err)
		}
		defer inFile.Close()

		reader := bufio.NewReader(inFile)
		img, _, err := image.Decode(reader)
		if err != nil {
			panic(err)
		}

		sizeOfMessage := GetMessageSizeFromImage(img) // Uses the library to check the message size

		msg := Decode(sizeOfMessage, img) // Read the message from the picture file

		if len(msg) != 0 {
			fmt.Println(string(msg))
		} else {
			fmt.Println("No message found")
		}
	},
}

func init() {
	decrypt.Flags().StringP("output", "o", "", "path to the output .PNG file. Default value is out.png")

	rootCmd.AddCommand(decrypt)
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
