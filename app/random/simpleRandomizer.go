package random

import (
	"math/rand"
	"time"

	"github.com/tthung1997/buddy/core/random"
)

type SimpleRandomizer struct {
	rand *rand.Rand
}

func NewSimpleRandomizer() *SimpleRandomizer {
	return &SimpleRandomizer{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *SimpleRandomizer) GetChoice(choices []random.Choice, count int) []random.Choice {
	// return count random choices from the list of choices taking weight into account
	// if count is greater than the length of the list, return the whole list
	// if count is less than or equal to 0, return an empty list
	if count <= 0 {
		return []random.Choice{}
	}

	if count >= len(choices) {
		return choices
	}

	// get the total weight of the list
	totalWeight := 0
	for _, choice := range choices {
		totalWeight += int(choice.Weight)
	}

	// get count random choices
	result := make([]random.Choice, 0, count)
	for i := 0; i < count; i++ {
		randomWeight := r.rand.Intn(totalWeight) + 1
		for _, choice := range choices {
			randomWeight -= int(choice.Weight)
			if randomWeight <= 0 {
				result = append(result, choice)
				totalWeight -= int(choice.Weight)
				choices = removeChoice(choices, choice)
				break
			}
		}
	}

	return result
}

func removeChoice(choices []random.Choice, choice random.Choice) []random.Choice {
	for i, c := range choices {
		if c == choice {
			return append(choices[:i], choices[i+1:]...)
		}
	}
	return choices
}
