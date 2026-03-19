package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"duo-abuser/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("RIOT_API_KEY")

	if key == "" {
		log.Fatal("RIOT_API_KEY not set")
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Enter tag: ")
	tag, _ := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	tag = strings.TrimSpace(tag)

	curr_summoner, err := api.GetSummoner(name, tag, key)
	println(curr_summoner.Puuid)
	println(curr_summoner.GameName)
	println(curr_summoner.TagLine)
}
