package cmd

import (
	"github.com/inclee/protogen/internal/gen_dart"
	"github.com/spf13/cobra"
)

var gDartCmd = &cobra.Command{
	Use:   "gdart",
	Short: "generate code files based on the protocol file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gen_dart.Gen(protoPath, codePath); err != nil {
			panic(err)
		}
	},
}

func init() {
	gDartCmd.Flags().StringVarP(&protoPath, "rpath", "r", protoPath, "proto file path dir default: ./internal/proto")
	gDartCmd.Flags().StringVarP(&codePath, "wpath", "w", codePath, "code file path dir default: ./")
	rootCmd.AddCommand(gDartCmd)
}
