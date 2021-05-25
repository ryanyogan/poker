package main

import (
	"fmt"
	"log"
	"os"

	"yogan.dev/poker"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker!")
	fmt.Println("Type {name} wins to record a win!")
	poker.NewCLI(store, os.Stdin).PlayPoker()
}
