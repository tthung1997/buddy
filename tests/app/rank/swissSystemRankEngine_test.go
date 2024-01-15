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

	if rankings[0].Additionals != "{\"winCount\":0,\"loseCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":0,\"loseCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[2].Additionals != "{\"winCount\":0,\"loseCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}
}

func TestSwiss_GetNextPair_Success(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Call Initialize
	entries := []string{"A", "B", "C"}
	rankings, err := swiss.Initialize(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Call GetNextPair
	pair, err := swiss.GetNextPair(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the pair is what we expect
	if pair.First.Id != 0 {
		t.Fatalf("unexpected pair: %v", pair)
	}

	if pair.Second.Id != 1 {
		t.Fatalf("unexpected pair: %v", pair)
	}
}

func TestSwiss_GetNextPair_NotEnoughPair(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Call Initialize
	entries := []string{"A"}
	rankings, err := swiss.Initialize(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Call GetNextPair
	pair, err := swiss.GetNextPair(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pair != nil {
		t.Fatalf("unexpected pair: %v", pair)
	}
}

func TestSwiss_GetNextPair_NoPair(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Initialize rankings
	rankings := []coreRank.RankingEntry{
		{
			Id:          0,
			Rank:        0,
			Value:       "A",
			Additionals: "{\"winCount\":1,\"loseCount\":0}",
		},
		{
			Id:          1,
			Rank:        0,
			Value:       "B",
			Additionals: "{\"winCount\":0,\"loseCount\":1}",
		},
	}

	// Call GetNextPair
	pair, err := swiss.GetNextPair(rankings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pair != nil {
		t.Fatalf("unexpected pair: %v", pair)
	}
}

func TestSwiss_UpdateRankings_FirstWin(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Initialize rankings
	rankings := []coreRank.RankingEntry{
		{
			Id:          0,
			Rank:        0,
			Value:       "A",
			Additionals: "{\"winCount\":1,\"loseCount\":0}",
		},
		{
			Id:          1,
			Rank:        0,
			Value:       "B",
			Additionals: "{\"winCount\":1,\"loseCount\":0}",
		},
	}

	// Initialize pair
	pair := coreRank.RankingEntryPair{
		First:  rankings[0],
		Second: rankings[1],
	}

	// Call UpdateRankings
	rankings, err := swiss.UpdateRankings(rankings, pair, coreRank.First)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the rankings are what we expect
	if len(rankings) != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Rank != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Rank != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "A" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":2,\"loseCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":1,\"loseCount\":1}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}
}

func TestSwiss_UpdateRankings_SecondWin(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Initialize rankings
	rankings := []coreRank.RankingEntry{
		{
			Id:          0,
			Rank:        0,
			Value:       "A",
			Additionals: "{\"winCount\":1,\"loseCount\":0}",
		},
		{
			Id:          1,
			Rank:        0,
			Value:       "B",
			Additionals: "{\"winCount\":1,\"loseCount\":0}",
		},
	}

	// Initialize pair
	pair := coreRank.RankingEntryPair{
		First:  rankings[0],
		Second: rankings[1],
	}

	// Call UpdateRankings
	rankings, err := swiss.UpdateRankings(rankings, pair, coreRank.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the rankings are what we expect
	if len(rankings) != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Id != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Rank != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Rank != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Value != "A" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":2,\"loseCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":1,\"loseCount\":1}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}
}

func TestSwiss_UpdateRankings_NotDonw(t *testing.T) {
	// Create a SwissSystemRankEngine
	swiss := appRank.SwissSystemRankEngine{}

	// Initialize rankings
	rankings := []coreRank.RankingEntry{
		{
			Id:          0,
			Rank:        0,
			Value:       "A",
			Additionals: "{\"winCount\":1,\"loseCount\":0}",
		},
		{
			Id:          1,
			Rank:        0,
			Value:       "B",
			Additionals: "{\"winCount\":1,\"loseCount\":0}",
		},
		{
			Id:          2,
			Rank:        0,
			Value:       "C",
			Additionals: "{\"winCount\":1,\"loseCount\":1}",
		},
	}

	// Initialize pair
	pair := coreRank.RankingEntryPair{
		First:  rankings[0],
		Second: rankings[1],
	}

	// Call UpdateRankings
	rankings, err := swiss.UpdateRankings(rankings, pair, coreRank.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the rankings are what we expect
	if len(rankings) != 2 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Id != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Id != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Rank != 0 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Rank != 1 {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Value != "B" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Value != "A" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[0].Additionals != "{\"winCount\":2,\"loseCount\":0}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}

	if rankings[1].Additionals != "{\"winCount\":1,\"loseCount\":1}" {
		t.Fatalf("unexpected rankings: %v", rankings)
	}
}
