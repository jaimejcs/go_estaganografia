package cmd

import (
	"bufio"
	"bytes"
	"image"
	"os"
	"strings"

	"github.com/jaimejcs/go_esteganografia/lib/encrypt"
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Oculta a mensagem dentro da imagem <path/to/file>",
	Long: `Oculta a mensagem dentro da imagem <path/to/file>
Exemplo: encrypt imagem.png -m 'Nosso segredo morre aqui' -o imagem_secreta.png`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var input = args[0]
		mensagem, _ := cmd.Flags().GetString("message")
		var output string

		if output = cmd.Flag("output").Value.String(); len(strings.TrimSpace(output)) == 0 {
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
		err = encrypt.Encode(encodedImg, img, []byte(mensagem)) // Calls library and Encodes the message into a new buffer
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
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("message", "m", "", "message to be encoded")
	encryptCmd.MarkFlagRequired("message")
	encryptCmd.Flags().StringP("output", "o", "", "path to the output .PNG file. Default value is out.png")
}
