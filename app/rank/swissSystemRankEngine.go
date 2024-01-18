package rank

import (
	"encoding/json"
	"sort"
	"time"

	coreRank "github.com/tthung1997/buddy/core/rank"
)

type SwissSystemRankEngine struct {
}

type SwissRankingEntryDescription struct {
	WinCount  int `json:"winCount"`
	LoseCount int `json:"loseCount"`
	ByeCount  int `json:"byeCount"`
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

func sortRankings(entries []coreRank.RankingEntry, forPairSelection bool) ([]coreRank.RankingEntry, error) {
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

		compareResult := (firstAdditional.WinCount > secondAdditional.WinCount) || ((firstAdditional.WinCount == secondAdditional.WinCount) && (firstAdditional.LoseCount < secondAdditional.LoseCount))

		if forPairSelection {
			return compareResult || ((firstAdditional.WinCount == secondAdditional.WinCount) && (firstAdditional.LoseCount == secondAdditional.LoseCount) && (firstAdditional.ByeCount > secondAdditional.ByeCount))
		} else {
			return compareResult || ((firstAdditional.WinCount == secondAdditional.WinCount) && (firstAdditional.LoseCount == secondAdditional.LoseCount) && (firstAdditional.ByeCount < secondAdditional.ByeCount))
		}
	})

	if sortErr != nil {
		return nil, sortErr
	}

	// Update rank
	for i, entry := range entries {
		entry.Rank = i
		entries[i] = entry
	}

	return entries, nil
}

func (r *SwissSystemRankEngine) Initialize(entries []string) ([]coreRank.RankingEntry, error) {
	rankings := make([]coreRank.RankingEntry, 0)

	for i, v := range entries {
		desc, err := toJson(SwissRankingEntryDescription{WinCount: 0, LoseCount: 0, ByeCount: 0})
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

func (r *SwissSystemRankEngine) GetNextPairs(entries []coreRank.RankingEntry) ([]coreRank.RankingEntry, []coreRank.RankingEntryPair, error) {
	entries, err := sortRankings(entries, true)
	if err != nil {
		return nil, nil, err
	}

	pairs := make([]coreRank.RankingEntryPair, 0)

	for i := 0; i < len(entries); i++ {
		entry := entries[i]

		additionals, err := fromJson(entry.Additionals)
		if err != nil {
			return nil, nil, err
		}

		if i+1 < len(entries) {
			nextEntry := entries[i+1]
			nextAdditionals, err := fromJson(nextEntry.Additionals)
			if err != nil {
				return nil, nil, err
			}

			if (additionals.WinCount == nextAdditionals.WinCount) && (additionals.LoseCount == nextAdditionals.LoseCount) {
				pairs = append(pairs, coreRank.RankingEntryPair{
					First:  entry,
					Second: nextEntry,
				})

				i++
				continue
			}
		}

		if i-1 >= 0 {
			prevEntry := entries[i-1]
			prevAdditionals, err := fromJson(prevEntry.Additionals)
			if err != nil {
				return nil, nil, err
			}

			if (additionals.WinCount == prevAdditionals.WinCount) && (additionals.LoseCount == prevAdditionals.LoseCount) {
				additionals.ByeCount++
				additionals.WinCount++

				entry.Additionals, err = toJson(additionals)
				if err != nil {
					return nil, nil, err
				}

				entries[i] = entry
			}
		}
	}

	entries, err = sortRankings(entries, false)
	if err != nil {
		return nil, nil, err
	}

	return entries, pairs, nil
}

func (r *SwissSystemRankEngine) UpdateRankings(entries []coreRank.RankingEntry, pairs []coreRank.RankingEntryPair) ([]coreRank.RankingEntry, error) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Id < entries[j].Id
	})

	for _, pair := range pairs {
		firstAdditionals, err := fromJson(pair.First.Additionals)
		if err != nil {
			return nil, err
		}

		secondAdditionals, err := fromJson(pair.Second.Additionals)
		if err != nil {
			return nil, err
		}

		if pair.Result == coreRank.First {
			firstAdditionals.WinCount++
			secondAdditionals.LoseCount++
		} else if pair.Result == coreRank.Second {
			firstAdditionals.LoseCount++
			secondAdditionals.WinCount++
		}

		entries[pair.First.Id].Additionals, err = toJson(firstAdditionals)
		if err != nil {
			return nil, err
		}

		entries[pair.Second.Id].Additionals, err = toJson(secondAdditionals)
		if err != nil {
			return nil, err
		}
	}

	entries, err := sortRankings(entries, false)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
