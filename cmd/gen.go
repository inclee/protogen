package cmd

import (
	"github.com/inclee/protogen/internal/gen_go"
	"github.com/spf13/cobra"
)

var protoPath = "internal/proto"
var codePath = "internal/"
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate code files based on the protocol file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gen_go.Gen(protoPath, codePath); err != nil {
			panic(err)
		}
	},
}

func init() {
	genCmd.Flags().StringVarP(&protoPath, "rpath", "r", protoPath, "proto file path dir default: ./internal/proto")
	genCmd.Flags().StringVarP(&codePath, "wpath", "w", codePath, "code file path dir default: ./internal/")
	rootCmd.AddCommand(genCmd)
}
