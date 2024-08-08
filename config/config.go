package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type AppConfig struct {
    App struct {
        Name    string `json:"name"`
        Version string `json:"version"`
    } `json:"app"`
    Server struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"server"`
    Nginx struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"nginx"`
}

var Config AppConfig
var ServerUrl string
var NginxUrl string

func LoadConfig(env string) {
    filePath := fmt.Sprintf("../../config/%s.json", env)
    file, err := os.Open(filePath)
    if err != nil {
        panic(fmt.Sprintf("Error opening config file: %v", err))
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        panic(fmt.Sprintf("Error reading config file: %v", err))
    }

    err = json.Unmarshal(data, &Config)
    if err != nil {
        panic(fmt.Sprintf("Error unmarshaling config: %v", err))
    }

    ServerUrl= fmt.Sprintf("%s:%d", Config.Server.Host, Config.Server.Port)
    NginxUrl= fmt.Sprintf("%s:%d", Config.Nginx.Host, Config.Nginx.Port)

}