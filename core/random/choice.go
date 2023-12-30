package random

import "time"

type Choice struct {
	Id              string
	Value           string
	Weight          int32
	Color           string
	UpdatedDateTime time.Time
}

type ChoiceList struct {
	Id              string
	Choices         []Choice
	UpdatedDateTime time.Time
}

type IChoiceListRepository interface {
	GetChoiceList(string) (ChoiceList, error)
	CreateOrUpdateChoiceList(ChoiceList) error
}
