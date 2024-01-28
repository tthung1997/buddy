package random

type IRandomizer interface {
	// GetChoice returns a random choice from the list of choices
	GetChoice([]Choice, int) []Choice
}
