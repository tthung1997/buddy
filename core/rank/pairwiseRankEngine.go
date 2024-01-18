package rank

type RankingEntryPair struct {
	First  RankingEntry
	Second RankingEntry
	Result PairwiseRankResult
}

type PairwiseRankResult string

const (
	First  PairwiseRankResult = "FIRST"
	Second PairwiseRankResult = "SECOND"
)

type IPairwiseRankEngine interface {
	Initialize([]string) ([]RankingEntry, error)
	GetNextPairs([]RankingEntry) ([]RankingEntry, []RankingEntryPair, error)
	UpdateRankings([]RankingEntry, []RankingEntryPair) ([]RankingEntry, error)
}
