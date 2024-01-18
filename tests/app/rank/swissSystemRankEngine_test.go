package rank

import (
	"testing"

	appRank "github.com/tthung1997/buddy/app/rank"
	coreRank "github.com/tthung1997/buddy/core/rank"
)

func TestSwiss_Initialize_Success(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Call Initialize
	entries := []string{"A", "B", "C"}
	rankings, err := swiss.Initialize(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the rankings are what we expect
	if len(rankings) != 3 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Id != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Rank != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Rank != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Rank != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "A" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Value != "C" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Additionals != "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}
}

func TestSwiss_GetNextPairs_Success(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Call Initialize
	entries := []string{"A", "B", "C"}
	rankings, err := swiss.Initialize(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Call GetNextPairs
	rankings, pairs, err := swiss.GetNextPairs(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that pairs is what we expect
	if len(rankings) != 3 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "C" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":1,\"loseCount\":0,\"byeCount\":1}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if len(pairs) != 1 {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].First.Id != 0 {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].Second.Id != 1 {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].First.Value != "A" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].Second.Value != "B" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].First.Additionals != "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].Second.Additionals != "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}
}

func TestSwiss_GetNextPairs_WithBye(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Initialize
	rankings := []coreRank.RankingEntry{
		{
			Id:          0,
			Rank:        0,
			Value:       "A",
			Additionals: "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}",
		},
		{
			Id:          1,
			Rank:        0,
			Value:       "B",
			Additionals: "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}",
		},
		{
			Id:          2,
			Rank:        0,
			Value:       "C",
			Additionals: "{\"winCount\":0,\"loseCount\":0,\"byeCount\":1}",
		},
	}

	// Call GetNextPairs
	rankings, pairs, err := swiss.GetNextPairs(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that pairs is what we expect
	if len(rankings) != 3 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":1,\"loseCount\":0,\"byeCount\":1}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if len(pairs) != 1 {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].First.Id != 2 {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].Second.Id != 0 {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].First.Value != "C" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].Second.Value != "A" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].First.Additionals != "{\"winCount\":0,\"loseCount\":0,\"byeCount\":1}" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}

	if pairs[0].Second.Additionals != "{\"winCount\":0,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected pairs: %v", pairs)
	}
}

func TestSwiss_GetNextPairs_NoPair(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Initialize
	rankings := []coreRank.RankingEntry{
		{
			Id:          0,
			Rank:        0,
			Value:       "A",
			Additionals: "{\"winCount\":3,\"loseCount\":0,\"byeCount\":0}",
		},
		{
			Id:          1,
			Rank:        0,
			Value:       "B",
			Additionals: "{\"winCount\":2,\"loseCount\":0,\"byeCount\":0}",
		},
		{
			Id:          2,
			Rank:        0,
			Value:       "C",
			Additionals: "{\"winCount\":1,\"loseCount\":0,\"byeCount\":0}",
		},
	}

	// Call GetNextPairs
	rankings, pairs, err := swiss.GetNextPairs(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that pairs is what we expect
	if len(rankings) != 3 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "A" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":3,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":2,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Id != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Value != "C" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Additionals != "{\"winCount\":1,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if len(pairs) != 0 {
		t.Fatalf("unexpected pairs: %v", pairs)
	}
}

func TestSwiss_UpdateRankings_FirstWin(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Call Initialize
	entries := []string{"A", "B", "C"}
	rankings, err := swiss.Initialize(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Call GetNextPairs
	rankings, pairs, err := swiss.GetNextPairs(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Set the result of the first pair to First
	pairs[0].Result = coreRank.First

	// Call UpdateRankings
	rankings, err = swiss.UpdateRankings(rankings, pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that rankings is what we expect
	if len(rankings) != 3 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "A" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":1,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Id != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Value != "C" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":1,\"loseCount\":0,\"byeCount\":1}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Additionals != "{\"winCount\":0,\"loseCount\":1,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}
}

func TestSwiss_UpdateRankings_SecondWin(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Call Initialize
	entries := []string{"A", "B", "C"}
	rankings, err := swiss.Initialize(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Call GetNextPairs
	rankings, pairs, err := swiss.GetNextPairs(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Set the result of the first pair to First
	pairs[0].Result = coreRank.Second

	// Call UpdateRankings
	rankings, err = swiss.UpdateRankings(rankings, pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that rankings is what we expect
	if len(rankings) != 3 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":1,\"loseCount\":0,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Id != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Value != "C" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":1,\"loseCount\":0,\"byeCount\":1}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Id != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Value != "A" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Additionals != "{\"winCount\":0,\"loseCount\":1,\"byeCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}
}
