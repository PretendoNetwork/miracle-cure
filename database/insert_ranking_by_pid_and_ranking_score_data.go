package database

import (
	"time"

	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"

	"github.com/PretendoNetwork/nex-go/v2/types"
)

func InsertRankingByPIDAndRankingScoreData(pid types.PID, rankingScoreData ranking_types.RankingScoreData, uniqueID types.UInt64) error {
	now := time.Now().UnixNano()
	if rankingScoreData.Score == 0 {
		return nil
	}

	_, err := Postgres.Exec(`
		INSERT INTO rankings (
			owner_pid,
			category,
			score,
			order_by,
			update_mode,
			groups,
			param,
			unique_id,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (owner_pid, unique_id, category) DO UPDATE
		SET score = excluded.score;`,
		pid,
		rankingScoreData.Category,
		rankingScoreData.Score,
		rankingScoreData.OrderBy,
		rankingScoreData.UpdateMode,
		rankingScoreData.Groups,
		rankingScoreData.Param&0x7FFFFFFFFFFFFFFF,
		uniqueID,
		now,
	)

	return err
}
