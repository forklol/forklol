package bitcoin

import (
	"backend/db"
	"errors"
	"log"
)

type Block struct {
	Coin       string  `json:"-" db:"coin"`
	Height     uint64  `json:"height" db:"height"`
	Difficulty float64 `json:"difficulty" db:"difficulty"`
	Work       float64 `json:"work" db:"work"`
	Rate1d     float64 `json:"rate1d" db:"rate1d"`
	Rate3d     float64 `json:"rate3d" db:"rate3d"`
	Rate7d     float64 `json:"rate7d" db:"rate7d"`
	Rate3h     float64 `json:"rate3h" db:"rate3h"`
	Rate12h    float64 `json:"rate12h" db:"rate12h"`
	Reward     float64 `json:"reward" db:"reward"`
	Price      float64 `json:"price" db:"price"`
	DARI       float64 `json:"dari" db:"dari"`
	Size       float64 `json:"size" db:"size"`
	Txs        float64 `json:"txs" db:"txs"`
	CDD        float64 `json:"cdd" db:"cdd"`
	Timestamp  uint64  `json:"timestamp" db:"timestamp"`
}

type Coin struct {
	Symbol        string
	BlockchairURL string
}

type ChainWork struct {
	Blocks uint64  `db:"blocks" json:"blocks"`
	Work   float64 `db:"work" json:"work"`
}

type ChainWorkTotal struct {
	ChainWork
}

type HashFactor struct {
	ChainWork      ChainWork `json:"chainwork"`
	HashRateFactor float64   `json:"factor"`
	AverageWork    float64   `json:"work_avg"`
	ActualWork     float64   `json:"work_actual"`
}

const CHAINSPLIT_WORK = 32729585000856628.00
const CHAINSPLIT_HEIGHT = 478558

func GetWorkSinceSplit(coin Coin) (*ChainWork, error) {
	work := ChainWork{}
	if err := db.GetDB().Get(&work, "SELECT sum(difficulty) as work, count(*) as blocks FROM blocks WHERE coin = ? AND height > ?", coin.Symbol, CHAINSPLIT_HEIGHT); err != nil {
		return nil, err
	}

	return &work, nil
}

func GetWorkInPeriod(coin Coin, start, end uint64) (*ChainWork, error) {
	work := ChainWork{}
	if err := db.GetDB().Get(&work, "SELECT sum(difficulty) as work, count(*) as blocks FROM blocks WHERE coin = ? and timestamp >= ? and timestamp <= ?", coin.Symbol, start, end); err != nil {
		return nil, err
	}

	return &work, nil
}

func GetWorkTotal(coin Coin) (*ChainWorkTotal, error) {
	work, err := GetWorkSinceSplit(coin)
	if err != nil {
		return &ChainWorkTotal{ChainWork{0, CHAINSPLIT_WORK}}, nil
		return nil, err
	}

	work.Work += CHAINSPLIT_WORK
	work.Blocks += CHAINSPLIT_HEIGHT

	return &ChainWorkTotal{*work}, nil
}

func GetHashFactorInPeriod(coin Coin, start, end uint64) (*HashFactor, error) {
	work, err := GetWorkInPeriod(coin, start, end)

	if err != nil {
		return nil, err
	}

	if work.Blocks == 0 {
		log.Println("No blocks found to calculate hashrate in period")
		return nil, errors.New("No blocks found")
	}

	t := float64(end-start) / 3600 / 24 * 144
	factor := float64(work.Blocks) / t
	avgwork := work.Work / float64(work.Blocks)

	return &HashFactor{
		ChainWork:      *work,
		HashRateFactor: factor,
		AverageWork:    avgwork,
		ActualWork:     avgwork * factor,
	}, nil
}
