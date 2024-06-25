package database

import (
	"time"

	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"

	"github.com/PretendoNetwork/nex-go/v2/types"
)

func InsertRankingByPIDAndRankingScoreData(pid *types.PID, rankingScoreData *ranking_types.RankingScoreData, uniqueID *types.PrimitiveU64) error {
	now := time.Now().UnixNano()
	if(rankingScoreData.Score.Value == 0){
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
		pid.Value(),
		rankingScoreData.Category.Value,
		rankingScoreData.Score.Value,
		rankingScoreData.OrderBy.Value,
		rankingScoreData.UpdateMode.Value,
		rankingScoreData.Groups.Value,
		rankingScoreData.Param.Value & 0x7FFFFFFFFFFFFFFF,
		uniqueID.Value,
		now,
	)

	return err
}