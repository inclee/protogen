package cmd

import (
	"github.com/inclee/protogen/internal/gen_go"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "remove all generated code files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gen_go.Clean(protoPath, codePath); err != nil {
			panic(err)
		}
	},
}

func init() {
	cleanCmd.Flags().StringVarP(&protoPath, "rpath", "r", protoPath, "proto file path dir default: ./internal/proto")
	cleanCmd.Flags().StringVarP(&codePath, "wpath", "w", codePath, "code file path dir default: ./internal/")
	rootCmd.AddCommand(cleanCmd)
}
