package rank

import "time"

type RankingEntry struct {
	Id              int
	Rank            int
	UpdatedDateTime time.Time
	Value           string
	Additionals     string
}
