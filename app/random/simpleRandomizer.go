package random

import (
	"math/rand"

	"github.com/tthung1997/buddy/core/random"
)

type SimpleRandomizer struct {
}

func (r *SimpleRandomizer) GetChoice(choices []random.Choice) random.Choice {
	// return a random choice from the list of choices taking weight into account
	totalWeight := 0
	for _, choice := range choices {
		totalWeight += choice.Weight
	}

	randomNumber := rand.Intn(totalWeight)

	for _, choice := range choices {
		randomNumber -= choice.Weight
		if randomNumber < 0 {
			return choice
		}
	}

	return random.Choice{} // return an empty choice if the list is empty
}
