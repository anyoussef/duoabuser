package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"duo-abuser/api"
)

func main() {
	key := os.Getenv("RIOT_API_KEY")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Enter tag: ")
	tag, _ := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	tag = strings.TrimSpace(tag)

	fmt.Print(name, tag)

	summoner, err := api.GetSummoner(name, tag, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Name:", summoner.GameName)
	fmt.Println("Tag:", summoner.TagLine)
}
