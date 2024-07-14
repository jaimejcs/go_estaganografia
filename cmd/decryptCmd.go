package cmd

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/jaimejcs/go_esteganografia/lib/decrypt"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Exibe a mensagem dentro da imagem <path/to/file>",
	Long: `Exibe a mensagem dentro da imagem <path/to/file>
Exemplo: decrypt imagem.png`,
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

		sizeOfMessage := decrypt.GetMessageSizeFromImage(img) // Uses the library to check the message size

		msg := decrypt.Decode(sizeOfMessage, img) // Read the message from the picture file

		if len(msg) != 0 {
			fmt.Println(string(msg))
		} else {
			fmt.Println("No message found")
		}
	},
}

func init() {
	decryptCmd.Flags().StringP("output", "o", "", "path to the output .PNG file. Default value is out.png")

	rootCmd.AddCommand(decryptCmd)
}
