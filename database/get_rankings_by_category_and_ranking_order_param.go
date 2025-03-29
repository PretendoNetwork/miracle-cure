package database

import (
	"database/sql"

	"github.com/PretendoNetwork/miracle-cure/globals"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"

	"github.com/PretendoNetwork/nex-go/v2/types"
)

func GetRankingsAndCountByCategoryAndRankingOrderParam(category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error) {
	rankings := types.NewList[ranking_types.RankingRankData]()

	rows, err := Postgres.Query(`
		SELECT
		owner_pid,
		score,
		groups,
		param
		FROM rankings WHERE category=$1 ORDER BY score DESC LIMIT $2 OFFSET $3`,
		category,
		rankingOrderParam.Length,
		rankingOrderParam.Offset,
	)
	if err != nil {
		return rankings, 0, err
	}

	row := 1
	for rows.Next() {
		rankingRankData := ranking_types.NewRankingRankData()
		rankingRankData.UniqueID = types.NewUInt64(0)
		rankingRankData.Order = types.NewUInt32(uint32(row))
		rankingRankData.Category = category

		var pid uint64
		err := rows.Scan(
			&pid,
			&rankingRankData.Score,
			&rankingRankData.Groups,
			&rankingRankData.Param,
		)

		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
		}

		commonDataRows, err := Postgres.Query(`
			SELECT
			common_data
			FROM common_datas WHERE owner_pid=$1 AND unique_id=$2`,
			pid,
			rankingRankData.UniqueID,
		)
		if err != nil {
			globals.Logger.Critical(err.Error())
			return rankings, 0, err
		}
		commonDataRows.Next()
		err = commonDataRows.Scan(
			&rankingRankData.CommonData,
		)

		if err != nil && err != sql.ErrNoRows {
			globals.Logger.Critical(err.Error())
		}

		rankingRankData.PrincipalID = types.NewPID(pid)

		if err == nil {
			rankings = append(rankings, rankingRankData)

			row += 1
		}
	}

	return rankings, uint32(len(rankings)), nil
}
