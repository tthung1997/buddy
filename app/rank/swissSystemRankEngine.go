package rank

import (
	"encoding/json"
	"errors"
	"sort"
	"time"

	coreRank "github.com/tthung1997/buddy/core/rank"
)

type SwissSystemRankEngine struct {
}

type SwissRankingEntryDescription struct {
	WinCount  int `json:"winCount"`
	LoseCount int `json:"loseCount"`
}

func toJson(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func fromJson(s string) (SwissRankingEntryDescription, error) {
	var v SwissRankingEntryDescription
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (r *SwissSystemRankEngine) Initialize(entries []string) ([]coreRank.RankingEntry, error) {
	rankings := make([]coreRank.RankingEntry, 0)

	for i, v := range entries {
		desc, err := toJson(SwissRankingEntryDescription{WinCount: 0, LoseCount: 0})
		if err != nil {
			return nil, err
		}

		entry := coreRank.RankingEntry{
			Id:              i,
			Rank:            0,
			UpdatedDateTime: time.Now(),
			Value:           v,
			Additionals:     desc,
		}

		rankings = append(rankings, entry)
	}

	return rankings, nil
}

func (r *SwissSystemRankEngine) GetNextPair(entries []coreRank.RankingEntry) (*coreRank.RankingEntryPair, error) {
	for i, entry := range entries {
		for j := i + 1; j < len(entries); j++ {
			if i == j {
				continue
			}

			var firstAdditional SwissRankingEntryDescription
			firstAdditional, err := fromJson(entry.Additionals)
			if err != nil {
				return nil, err
			}

			second := entries[j]
			secondAdditional, err := fromJson(second.Additionals)
			if err != nil {
				return nil, err
			}

			// Find pair with the same win/lose count
			if (firstAdditional.WinCount == secondAdditional.WinCount) && (firstAdditional.LoseCount == secondAdditional.LoseCount) {
				return &coreRank.RankingEntryPair{First: entry, Second: second}, nil
			}
		}
	}

	// If there is no pair left, return nil
	return nil, nil
}

func (r *SwissSystemRankEngine) UpdateRankings(entries []coreRank.RankingEntry, nextPair coreRank.RankingEntryPair, result coreRank.PairwiseRankResult) ([]coreRank.RankingEntry, error) {
	first := nextPair.First
	second := nextPair.Second

	firstAdditional, err := fromJson(first.Additionals)
	if err != nil {
		return nil, err
	}

	secondAdditional, err := fromJson(second.Additionals)
	if err != nil {
		return nil, err
	}

	// Update win/lose count
	if result == coreRank.First {
		firstAdditional.WinCount++
		secondAdditional.LoseCount++
	} else if result == coreRank.Second {
		firstAdditional.LoseCount++
		secondAdditional.WinCount++
	} else {
		return nil, errors.New("invalid result")
	}

	first.Additionals, err = toJson(firstAdditional)
	if err != nil {
		return nil, err
	}

	second.Additionals, err = toJson(secondAdditional)
	if err != nil {
		return nil, err
	}

	// Update rankings
	for i, entry := range entries {
		if entry.Id == first.Id {
			entries[i] = first
		} else if entry.Id == second.Id {
			entries[i] = second
		}
	}

	// Check if there is any pair left
	pair, err := r.GetNextPair(entries)
	if err != nil {
		return nil, err
	}

	// If there is no pair left, sort the entries
	if pair == nil {
		var sortErr error

		sort.Slice(entries, func(i, j int) bool {
			firstAdditional, err := fromJson(entries[i].Additionals)
			if err != nil {
				sortErr = err
				return false
			}

			secondAdditional, err := fromJson(entries[j].Additionals)
			if err != nil {
				sortErr = err
				return false
			}

			return (firstAdditional.WinCount > secondAdditional.WinCount) || ((firstAdditional.WinCount == secondAdditional.WinCount) && (firstAdditional.LoseCount < secondAdditional.LoseCount))
		})

		if sortErr != nil {
			return nil, sortErr
		}

		// Update rank
		for i, entry := range entries {
			entry.Rank = i
			entries[i] = entry
		}
	}

	return entries, nil
}
