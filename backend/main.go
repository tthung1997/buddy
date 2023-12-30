package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	appRandom "github.com/tthung1997/buddy/app/random"
	"github.com/tthung1997/buddy/app/random/choice"
	coreRandom "github.com/tthung1997/buddy/core/random"
)

func main() {
	// Create a sample ChoiceList with 3 Choices
	choiceList := coreRandom.ChoiceList{
		Id: uuid.New().String(),
		Choices: []coreRandom.Choice{
			{Id: uuid.New().String(), Weight: 1, Value: "Choice 1", Color: "Red", UpdatedDateTime: time.Now()},
			{Id: uuid.New().String(), Weight: 2, Value: "Choice 2", Color: "Green", UpdatedDateTime: time.Now()},
			{Id: uuid.New().String(), Weight: 3, Value: "Choice 3", Color: "Blue", UpdatedDateTime: time.Now()},
		},
		UpdatedDateTime: time.Now(),
	}

	// Use LocalChoiceListRepository to store it to a file
	var repo coreRandom.IChoiceListRepository
	repo = choice.NewLocalChoiceListRepository("choiceLists.json")
	err := repo.CreateOrUpdateChoiceList(choiceList)
	if err != nil {
		fmt.Println("Error creating or updating ChoiceList:", err)
		return
	}

	// Use SimpleRandomizer to get a random Choice
	var randomizer coreRandom.IRandomizer
	randomizer = &appRandom.SimpleRandomizer{}
	randomChoice := randomizer.GetChoice(choiceList.Choices)

	// Print the random Choice
	fmt.Println("Random Choice:", randomChoice.Value)
}
