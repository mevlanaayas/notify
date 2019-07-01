package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadEnv() {
	var urls [] string
	envContent, err := ioutil.ReadFile("./.env")
	if err == nil {
		fmt.Println("reading env vars from file...")
		urls = strings.Split(string(envContent), "\n")
	} else {
		for _, url := range urls {
			temp := strings.Split(string(url), "=")
			key, value := temp[0], temp[1]
			err = os.Setenv(key, value)
		}
		fmt.Println("reading env file successful :)")
	}

}
