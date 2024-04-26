package cmanager

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)



var (
	included []string 
	excluded []string 
	image string
	tag string
	test bool
	templateFile string

	rootCmd = &cobra.Command{
		Long: fmt.Sprintf(`Manage compose services in %s`, usersDir),
	
  	}
  	stop = &cobra.Command{
		Use:   "stop",
		Short: "Stop compose",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			Stop(included, excluded, test)
		},
  	}
	start = &cobra.Command{
		Use:   "start",
		Short: "Start compose",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			Start(included, excluded, test)
		},
  	}
	generate = &cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			GenerateCompose(image, tag, included, excluded, test)
		},
  	}
)

func init(){
	rootCmd.AddCommand(stop)
	rootCmd.AddCommand(start)
	rootCmd.AddCommand(generate)
	rootCmd.PersistentFlags().StringArrayVarP(&included, "users", "u", nil, "user name to include")
	rootCmd.PersistentFlags().StringArrayVarP(&excluded, "eusers", "e", nil, "user name to exclude")
	rootCmd.PersistentFlags().BoolVarP(&test, "dry", "d", false, "Dry test run")
	generate.Flags().StringVarP(&image, "image", "i", "", "Image name")
	generate.Flags().StringVarP(&tag, "tag", "t", "", "Tag name")
	// generate.Flags().StringVarP(&templateFile, "template-file", "f", "", "Tag name")
	generate.MarkFlagRequired("image")
	generate.MarkFlagRequired("tag")
	generate.MarkFlagRequired("template-file")
	// generate.MarkFlagFilename("template-file")
	
	
}
  func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
  }