package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Oculta a mensagem dentro da imagem <path/to/file>",
	Long: `Oculta a mensagem dentro da imagem <path/to/file>
			Exemplo: encrypt imagem.png -m 'Nosso segredo morre aqui' -o imagem_secreta.png`,
	Run: func(cmd *cobra.Command, args []string) {
		var input = args[0]
		var mensagem = cmd.Flag("message").Value.String()
		var output string

		if output = cmd.Flag("output").Value.String(); len(strings.Trim(output, " ")) == 0 {
			output = "out.png"
		}

		inFile, err := os.Open(input) // Opens input file provided in the flags
		if err != nil {
			panic(err)
		}
		defer inFile.Close()

		reader := bufio.NewReader(inFile) // Reads binary data from picture file
		img, _, err := image.Decode(reader)
		if err != nil {
			panic(err)
		}
		encodedImg := new(bytes.Buffer)
		err = Encode(encodedImg, img, []byte(mensagem)) // Calls library and Encodes the message into a new buffer
		if err != nil {
			panic(err)
		}
		outFile, err := os.Create(output) // Creates file to write the message into
		if err != nil {
			panic(err)
		}
		bufio.NewWriter(outFile).Write(encodedImg.Bytes()) // writes file to disk
	},
}

func init() {
	encrypt.Flags().StringP("message", "m", "", "message to be encoded")
	encrypt.MarkFlagRequired("message")
	encrypt.Flags().StringP("output", "o", "", "path to the output .PNG file. Default value is out.png")

	rootCmd.AddCommand(encrypt)
}

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
	if MaxEncodeSize(rgbImage) < messageLength+4 {
		return errors.New("message too large for image")
	}

	one, two, three, four := splitToBytes(messageLength)

	message = append([]byte{four}, message...)
	message = append([]byte{three}, message...)
	message = append([]byte{two}, message...)
	message = append([]byte{one}, message...)

	ch := make(chan byte, 100)

	go getNextBitFromString(message, ch)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbImage.NRGBAAt(x, y) // get the color at this pixel

			/*  RED  */
			bit, ok = <-ch
			if !ok { // if we don't have any more bits left in our message

				rgbImage.SetNRGBA(x, y, c)
				break
			}
			setLSB(&c.R, bit)

			/*  GREEN  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetNRGBA(x, y, c)
				break
			}
			setLSB(&c.G, bit)

			/*  BLUE  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetNRGBA(x, y, c)
				break
			}
			setLSB(&c.B, bit)

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

	rgbImage := imageToNRGBA(pictureInputFile)

	return EncodeNRGBA(writeBuffer, rgbImage, message)

}
