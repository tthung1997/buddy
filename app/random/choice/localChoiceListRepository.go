package choice

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tthung1997/buddy/core/random"
)

type LocalChoiceListRepository struct {
	filePath string
}

func NewLocalChoiceListRepository(filePath string) (*LocalChoiceListRepository, error) {
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		choiceLists := make(map[string]random.ChoiceList)
		data, err := json.MarshalIndent(choiceLists, "", " ")
		if err != nil {
			return nil, err
		}

		_, err = file.Write(data)
		if err != nil {
			return nil, err
		}
	}

	return &LocalChoiceListRepository{filePath: filePath}, nil
}

func (r *LocalChoiceListRepository) GetChoiceList(id string) (random.ChoiceList, error) {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return random.ChoiceList{}, err
	}

	var choiceLists map[string]random.ChoiceList
	err = json.Unmarshal(file, &choiceLists)
	if err != nil {
		return random.ChoiceList{}, err
	}

	choiceList, exists := choiceLists[id]
	if !exists {
		return random.ChoiceList{}, fmt.Errorf("ChoiceList with ID %s not found", id)
	}

	return choiceList, nil
}

func (r *LocalChoiceListRepository) CreateOrUpdateChoiceList(choiceList random.ChoiceList) error {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var choiceLists map[string]random.ChoiceList
	if len(file) == 0 {
		choiceLists = make(map[string]random.ChoiceList)
	} else {
		err = json.Unmarshal(file, &choiceLists)
		if err != nil {
			return err
		}
	}

	choiceLists[choiceList.Id] = choiceList

	file, err = json.MarshalIndent(choiceLists, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(r.filePath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
