package main

import (
	"cos-storager/cmd"

	"qiniupkg.com/x/log.v7"
)

func main() {
	log.Println("Start Server")
	cmd.Execute()
}
