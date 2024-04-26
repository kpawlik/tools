package cmanager

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"text/template"
)

const IMAGE="harbor.delivery.iqgeo.cloud/stedin/capture7031-is7052"
const TAG="latest"

const COMPOSE=`
name: "stedin-{{.UserName}}"
services:
  memcache-{{.UserName}}:
    container_name: memcache-{{.UserName}}
    image: memcached:1.6.19
    restart: always
    ports:
      - {{.MemcachePort}}:11211
  stedin-{{.UserName}}:
    container_name: stedin-{{.UserName}}
    image: {{.ImageName}}:{{.ImageTag}}
    restart: always
    depends_on:
      - memcache-{{.UserName}}
    environment:
      DEBUG: "true"
      PGHOST: std-flex-psql-t01.privatelink.postgres.database.azure.com
      PGPORT: 5432
      PGUSER: adminstdflexpsqlt01
      PGPASSWORD: kKDdkeQsOEk1eMLL
      PGDATABASE: iqgeo-dev-7.0
      MYW_DB_HOST: std-flex-psql-t01.privatelink.postgres.database.azure.com
      MYW_DB_PORT: 5432
      MYW_DB_USERNAME: adminstdflexpsqlt01
      MYW_DB_PASSWORD: kKDdkeQsOEk1eMLL
      MYW_DB_NAME: iqgeo-dev-7.0
      WSGI_THREADS: 5
      WSGI_PROCESSES: 2 
      BEAKER_SESSION_TYPE: ext:memcached
      BEAKER_SESSION_URL: std-eo-mbmr-d01:{{.MemcachePort}}
      REPLICATION_SYNC_URL: http://std-eo-mbmr-d01:{{.ApachePort}}/
    volumes:
      - ./../modules/custom:/opt/iqgeo/platform/WebApps/myworldapp/modules/custom
      - /usr/iqgeo/users/shared-data:/shared-data
      - ./pre-apache-extra-endpoint-7.0.sh:/opt/iqgeo/pre-apache-extra-endpoint.sh
    ports:
      - {{.ApachePort}}:8080
`

type User struct{
	UserName string
	ImageName string
	ImageTag string
	MemcachePort string
	ApachePort string
}

var (
	Users = []User{
		User{
			UserName: "test",
			MemcachePort: "11390",
			ApachePort: "8090",
		},
		User{
			UserName: "jaligato",
			MemcachePort: "11391",
			ApachePort: "8091",
		},
		User{
			UserName: "marenas",
			MemcachePort: "11392",
			ApachePort: "8092",
		},
		User{
			UserName: "opastrana",
			MemcachePort: "11393",
			ApachePort: "8093",
		},
		User{
			UserName: "psonawane",
			MemcachePort: "11394",
			ApachePort: "8094",
		},
		User{
			UserName: "rllagas",
			MemcachePort: "11395",
			ApachePort: "8095",
		},
		User{
			UserName: "kpawlik",
			MemcachePort: "11396",
			ApachePort: "8096",
		},
	}
	usersDir = "/usr/iqgeo/users/"
	templ = template.New("compose")
	eUsers []string
)



func getComposePath(user string) string {
	userDir := filepath.Join(usersDir, user, "IQGeo%20platform","local_dev")
	composeFile := fmt.Sprintf("compose.%s.yaml", user)
	return filepath.Join(userDir, composeFile)
}
func userNameOk(user string,  included []string, excluded []string) bool{
	if slices.Contains(excluded, user){
		return false
	}
	if len(included) > 0 && !slices.Contains(included, user){
		return false 
	}
	return true

}
func GenerateCompose(image string, tag string, included []string, excluded []string, test bool){
	var (
		err error
		composeFile io.WriteCloser
	)
	if len(image) == 0 && len(tag) == 0{
		log.Printf("image and tag are required to generate compose")
		flag.PrintDefaults()
		return
	}
	if templ, err = templ.Parse(COMPOSE); err != nil{
		log.Fatalf("Parse template error %v", err)
	}
	for _, user := range Users {
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
		
		if err = templ.Execute(composeFile, user); err != nil{
			log.Printf("Error write file %s %v", composePath, err)
		}
		log.Printf("Done %s", user.UserName)
		
	}
}

func Stop(included []string, excluded []string, test bool){
	var (
		err error
		cmd *exec.Cmd
		stderr []byte
	)
	for _, user := range Users {
		if !userNameOk(user.UserName, included, excluded){
			continue
		}
		composePath := getComposePath(user.UserName)
		params := []string{"compose", "-f", composePath, "down"}
		log.Printf("Stop %s", composePath)
		if test{
			continue
		}
		cmd = exec.Command("docker", params...); 
		stderr, err = cmd.CombinedOutput()
		if err != nil{
			log.Printf("Error %v, %s", err, string(stderr))
		}
	}
}

func Start(userNames []string, excluded []string, test bool){
	var (
		err error
		cmd *exec.Cmd
		stderr []byte
	)
	for _, user := range Users {
		if !userNameOk(user.UserName, included, excluded){
			continue
		}
		composePath := getComposePath(user.UserName)
		params := []string{"compose", "-f", composePath, "up", "-d"}
		log.Printf("Start %s", composePath)
		if test{
			continue
		}
		cmd = exec.Command("docker", params...); 
		stderr, err = cmd.CombinedOutput()
		if err != nil{
			log.Printf("Error %v, %s", err, string(stderr))
		}
	}
	
}
