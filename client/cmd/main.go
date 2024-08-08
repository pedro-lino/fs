package main

import (
	"bufio"
	"fmt"
	client "fs/client/internal"
	"fs/config"
	"os"
	"strings"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	config.LoadConfig(env)
	for {
		fmt.Println("Choose an option: (1) Upload directory, (2) Download file, (exit to quit): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Enter directory path to upload (or type 'exit' to quit): ")
			dirPath, _ := reader.ReadString('\n')
			dirPath = strings.TrimSpace(dirPath)
			if dirPath == "exit" {
				break
			}
			client.ProcessPaths(config.NginxUrl, dirPath, env)
		case "2":
			fmt.Println("Enter the name of the file to download (or type 'exit' to quit): ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)
			if fileName == "exit" {
				break
			}
            fmt.Println("file name:", fileName)
			client.DownloadFile(config.NginxUrl, fileName)
		case "exit":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}