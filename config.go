package app

import (
    "code.google.com/p/gcfg"
)

type Config struct {
    Db struct {
        Disable    bool
        Driver     string
        Datasource string
    }

    Net struct {
        Listen_host string
        Listen_port int
    }

    Site struct {
        Host        string
        Disabled    bool
        Title       string
        Author      string
        Description string
        Copyright   string
        Keywords    string
        Email       string
        Phone       string
        UploadPath  string
    }
}

func loadConfig(file string) (res Config) {

    if err := gcfg.ReadFileInto(&res, file); err != nil {
        panic("Config error: Readfile error: " + err.Error())
    }

    return
}
