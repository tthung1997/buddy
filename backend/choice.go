package main

import (
	"log"

	appChoice "github.com/tthung1997/buddy/app/random/choice"
	grpcChoice "github.com/tthung1997/buddy/framework/grpc/choice"
)

const choiceListFile = ".db/choiceLists.json"

func main() {
	// Use LocalChoiceListRepository to store it to a file
	repo, err := appChoice.NewLocalChoiceListRepository(choiceListFile)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// start the server
	server := grpcChoice.NewChoiceServer(repo)
	server.Run()
}
