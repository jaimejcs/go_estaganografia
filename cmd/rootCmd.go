package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "EsteganoGo",
	Short: "EsteganoGo é uma CLI de esteganografia feita em Golang",
	Long: `É uma ferramenta de esteganografia que usa da técnica LSB
para ocultar mensagens em arquivos png
Feito como trabalho para a matéria Linguagens e Paradigmas de Progamação`,
	Args: cobra.MinimumNArgs(1),
	Run:  func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
