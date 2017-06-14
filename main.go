package main

import (
	"github.com/bitly/go-simplejson"
	"github.com/version_upload/sftp_client"
	"log"
	"os"
)

var config = make(map[string]string)

func init() {
	configFile, err := os.OpenFile("config.json", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	json, err := simplejson.NewFromReader(configFile)
	if err != nil {
		log.Fatal(err)
	}

	config["srcDir"], _ = json.Get("src_dir").String()
	config["dstDir"], _ = json.Get("dst_dir").String()
	config["host"], _ = json.Get("host").String()
	config["user"], _ = json.Get("user").String()
	config["passwd"], _ = json.Get("passwd").String()
}

func existDir(p string) (bool, error) {
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return true, err
		} else {
			return false, err
		}
	}

	return false, err
}

func main() {
	client, err := sftp_client.NewSftpClient(config["host"], config["user"], config["passwd"], true)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if exist, err := existDir(config["srcDir"]); exist {
		log.Fatal(err)
	}

	err1 := client.PutDir(config["srcDir"], config["dstDir"])
	if err1 != nil {
		log.Fatal(err1)
	}
}
