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
	GetNextPair() (RankingEntryPair, error)
	SetResult(RankingEntryPair, PairwiseRankResult) error
	GetRanking() ([]RankingEntry, error)
}
