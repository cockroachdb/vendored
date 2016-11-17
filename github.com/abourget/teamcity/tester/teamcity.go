package main

import (
	"fmt"

	"github.com/abourget/teamcity"
)

func main() {
	client := teamcity.New("myinstance.example.com", "username", "password")

	b, err := client.QueueBuild("Project_build_task", "master", nil)
	if err != nil {
		fmt.Printf("You're outta luck: %s\n", err)
		return
	}

	fmt.Printf("Build: %#v\n", b)
}
