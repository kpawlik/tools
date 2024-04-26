package main

import (
	"github.pl/kpawlik/cmanager"
)
var(
	image = ""
	tag = ""
	operation = "restart"
	userName string
	test = false
	excludedUsers string
	eUsers []string
)

func init() {
	// flag.StringVar(&tag, "tag", "", "Tag to replace in users compose files")
	// flag.StringVar(&image, "image", "", "Image to replace in users compose files")
	// flag.StringVar(&userName, "user", "", "User name")
	// flag.StringVar(&excludedUsers, "euser", "", "Excluded user names, coma separated")

	// flag.StringVar(&operation, "operation", "restart", "[stop, start, restart, generate]")
	// flag.BoolVar(&test, "test", false, "Test output")

	// flag.Parse()
	// if len(image) > 0 && len(tag) == 0 {
	// 	log.Println("-tag value is required when change the image")
	// 	os.Exit(1)
	// }
	// if len(excludedUsers) > 0 {
	// 	eUsers = strings.Split(excludedUsers, ",")
	// }

}
func main() {
	// switch operation {
	// case "start":
	// 	cmanager.Start(userName, test)
	// case "stop":
	// 	cmanager.Stop(userName, test)
	// case "restart":
	// 	cmanager.Stop(userName, test)
	// 	cmanager.Start(userName, test)
	// case "generate":
	// 	cmanager.GenerateCompose(image, tag, userName, test)
	// }

	cmanager.Execute()

}