package cmanager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)


const (
	TEMPLATE_FILE = "compose.template"
	CONFIG_FILE = "config.json"
)
var (
	included []string 
	excluded []string 
	image string
	tag string
	dry bool
	templateFile string
	configFile string
	config Config
	all bool

	rootCmd = &cobra.Command{
		Long: fmt.Sprintf(`Manage compose services in %s`, config.UsersDir),
		PersistentPreRun: func(cmd *cobra.Command, args []string){
			if templateFile == TEMPLATE_FILE{
				templateFile = getFilePath(templateFile)
			}
			if configFile == CONFIG_FILE{
				configFile = getFilePath(configFile)
			}
			if c, err := LoadConfig(configFile); err != nil{
				log.Fatalf("%v", err)
			}else{
				config = c
			}
		},
  	}
  	stop = &cobra.Command{
		Use:   "stop",
		Short: "Stop compose",
		Long:  `Run docker compose 'down' for  user compose file`,
		Run: func(cmd *cobra.Command, args []string) {
			users := FilterUsers(config.Users, included, excluded)
			Stop(users, dry)
		},
  	}
	start = &cobra.Command{
		Use:   "start",
		Short: "Start compose",
		Long:  `Run docker compose 'up -d' for  user compose file`,
		Run: func(cmd *cobra.Command, args []string) {
			users := FilterUsers(config.Users, included, excluded)
			Start(users, dry)
		},
  	}

	restart = &cobra.Command{
		Use:   "restart",
		Short: "Restart compose",
		Long:  `Run docker compose 'down' and 'up -d' for  user compose file`,
		Run: func(cmd *cobra.Command, args []string) {
			users := FilterUsers(config.Users, included, excluded)
			Restart(users, dry)
		},
  	}
	generate = &cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Long:  `Generate compose yaml file for users based on template file`,
		Run: func(cmd *cobra.Command, args []string) {
			users := FilterUsers(config.Users, included, excluded)
			GenerateCompose(image, tag, users, templateFile, dry)
		},
  	}
)

// Returns path for file dir path is executable dir
func getFilePath(fileName string) string{
	ex, _ := os.Executable()
	dir := filepath.Dir(ex)
	return filepath.Join(dir, fileName)
}

func init(){
	rootCmd.AddCommand(stop)
	rootCmd.AddCommand(start)
	rootCmd.AddCommand(generate)
	rootCmd.AddCommand(restart)
	rootCmd.PersistentFlags().StringArrayVarP(&included, "users", "u", nil, "user name to include")
	rootCmd.PersistentFlags().StringArrayVarP(&excluded, "eusers", "e", nil, "user name to exclude")
	rootCmd.PersistentFlags().BoolVarP(&dry, "dry", "d", false, "Dry test run")
	rootCmd.PersistentFlags().StringVarP(&configFile, "cfg", "", CONFIG_FILE, "Users config file")
	rootCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "Run for all users")
	generate.Flags().StringVarP(&image, "image", "i", "", "Image name")
	generate.Flags().StringVarP(&tag, "tag", "t", "", "Tag name")
	generate.Flags().StringVarP(&templateFile, "template", "", TEMPLATE_FILE, "Compose file template")
	generate.MarkFlagRequired("image")
	generate.MarkFlagRequired("tag")
	rootCmd.MarkFlagsOneRequired("users", "eusers", "all")

	
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
  }