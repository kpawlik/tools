package cmanager

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"text/template"
)

// return full path for compose file
func getComposePath(user string) string {
	userDir := filepath.Join(config.UsersDir, user, config.RepoName,"local_dev")
	composeFile := fmt.Sprintf("compose.%s.yaml", user)
	return filepath.Join(userDir, composeFile)
}

// Check if user is in excluded or included list
func userNameOk(user string,  included []string, excluded []string) bool{
	if slices.Contains(excluded, user){
		return false
	}
	if len(included) > 0 && !slices.Contains(included, user){
		return false 
	}
	return true
}

func FilterUsers(users []User, included []string, excluded[] string) []User{
	res := []User{}
	for _, user := range users{
		if userNameOk(user.UserName, included, excluded){
			res = append(res, user)
		}
	}
	return res
}

func GenerateCompose(image string, tag string, users []User, templatePath string, test bool)(err error){
	var (
		composeFile io.WriteCloser
		buff []byte
		composeT = template.New("compose")
	)
	if buff, err = os.ReadFile(templatePath); err !=nil{
		err = fmt.Errorf("error reading template error %v", err)
		return
	}
	if composeT, err = composeT.Parse(string(buff)); err != nil{
		err = fmt.Errorf("parse template error %v", err)
		return
	}
	for _, user := range users {
		if !userNameOk(user.UserName, included, excluded){
			continue
		}
		if len(image)>0{
			user.ImageName = image
		}
		if len(tag)>0{
			user.ImageTag = tag
		}
		composePath := getComposePath(user.UserName)
		if test{
			composeFile = os.Stdout
		}else{
			if composeFile, err = os.OpenFile(composePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil{
				log.Printf("Error opening file %s, %v", composePath, err)
			}
		 	defer composeFile.Close()
		}
		
		if err = composeT.Execute(composeFile, user); err != nil{
			log.Printf("Error write file %s %v", composePath, err)
		}
		log.Printf("Done %s", user.UserName)
	}
	return
}

func Stop(users []User, test bool){
	composeCommand("stop", users, test)
}

func Start(users []User, test bool){
	composeCommand("start", users, test)
}

func Restart(users []User, test bool){
	composeCommand("stop", users, test)
	composeCommand("start", users, test)
}

func composeCommand(operation string, users []User, dry bool) {
	var (
		err error
		cmd *exec.Cmd
		stderr []byte
		paramsMap = map[string][]string{
			"start": {"compose", "-f", "", "up", "-d"},
			"stop": {"compose", "-f", "", "down"},
		}
	)
	for _, user := range users {
		composePath := getComposePath(user.UserName)
		params := paramsMap[operation]
		log.Printf("%s compose %s", operation, composePath)
		if dry{
			continue
		}
		cmd = exec.Command("docker", params...); 
		stderr, err = cmd.CombinedOutput()
		if err != nil{
			log.Printf("Error running %s on %s. %v", operation, composePath, err)
		}
		log.Printf("%s", string(stderr))
	}
}
