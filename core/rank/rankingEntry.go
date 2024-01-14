package rank

import "time"

type RankingEntry struct {
	Id              string
	Rank            int32
	UpdatedDateTime time.Time
	Description     string
}
