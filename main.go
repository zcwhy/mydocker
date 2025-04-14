package main

import (
	"mydocker/cmd"
	"mydocker/log"
)

func main() {
	log.LogInit()
	cmd.InitCmd()
}
