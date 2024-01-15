package rank

type RankingEntryPair struct {
	First  RankingEntry
	Second RankingEntry
}

type PairwiseRankResult string

const (
	First  PairwiseRankResult = "FIRST"
	Second PairwiseRankResult = "SECOND"
)

type IPairwiseRankEngine interface {
	Initialize([]string) ([]RankingEntry, error)
	GetNextPair([]RankingEntry) (RankingEntryPair, error)
	UpdateRankings([]RankingEntry, RankingEntryPair, PairwiseRankResult) ([]RankingEntry, error)
}
