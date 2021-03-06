package cmd

import (
	"fmt"
	"github.com/mkawserm/abesh/constant"
	"github.com/spf13/cobra"
	"os"
)

var abeshCMD = &cobra.Command{
	Use:   constant.Name,
	Short: constant.ShortDescription,
	Long:  constant.LongDescription,
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

var abeshVersion = &cobra.Command{
	Use:   "version",
	Short: "Abesh version",
	Long:  "Show abesh version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constant.Version)
		os.Exit(0)
	},
}

var abeshAuthors = &cobra.Command{
	Use:   "authors",
	Short: "Abesh authors",
	Long:  "Show abesh authors name",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constant.Authors)
		os.Exit(0)
	},
}

func init() {
	abeshCMD.AddCommand(abeshVersion, abeshAuthors)
}
