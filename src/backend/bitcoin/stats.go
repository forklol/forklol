package bitcoin

import (
	"backend/db"
	"fmt"
	"sort"
	"sync"
	"time"
	"runtime"
)

// Very crude way to get data in parallel
func GetHistoricalDataDescending(coin Coin, to, from, interval uint64) *[]*ChainAverage {
	data := make([]*ChainAverage, 0)

	lock := sync.Mutex{}

	var done uint64 = 0
	target := (to - from) / interval

	running := 0
	threads := runtime.NumCPU()

	for n := to; n > from; n -= interval {

		for {
			if running < threads {
				running++
				break
			}
			time.Sleep(time.Millisecond)
		}

		go func(coin Coin, n, interval uint64) {
			stats := GetLastAverages(coin, n, interval)
			lock.Lock()
			data = append(data, stats)
			done++
			running--
			lock.Unlock()
		}(coin, n, interval)
	}

	for {
		if done == target {
			break;
		}
		time.Sleep(time.Millisecond)
	}

	sort.Slice(data, func(a, b int) bool {
		return data[a].Timestamp < data[b].Timestamp
	})

	return &data
}

func GetLastAverages(coin Coin, from, span uint64) *ChainAverage {
	return &ChainAverage{
		Coin:        coin.Symbol,
		Height:      statGetMaxHeight(coin, from, span),
		Blocks:      statGetBlockCount(coin, from, span),
		AvgDiff:     statGetAvgDifficulty(coin, from, span),
		Work:        statGetWork(coin, from, span),
		Rate:        statGetRateCalulated(coin, from, 24*3600),
		Rate3h:      statGetRateCalulated(coin, from, 3*3600),
		Rate12h:     statGetRateCalulated(coin, from, 12*3600),
		Rate1d:      statGetRateCalulated(coin, from, 24*3600),
		Rate3d:      statGetRateCalulated(coin, from, 3*24*3600),
		Rate7d:      statGetRateCalulated(coin, from, 7*24*3600),
		RateP:       statGetRateP(coin, from, span),
		TotalReward: statGetTotalReward(coin, from, span),
		AvgReward:   statGetAvgReward(coin, from, span),
		CoinPrice:   statGetAvgPrice(coin, from, span),
		DARI:        statGetAvgDARI(coin, from, span),
		Size:        statGetAvgSize(coin, from, span),
		Txs:         statGetAvgTxs(coin, from, span),
		CDD:         statGetAvgCDD(coin, from, span),
		Timestamp:   from,
	}
}

const restrict = " WHERE b.coin = ? AND b.timestamp <= ? AND b.timestamp > ?"

func statGetMaxHeight(coin Coin, ts, interval uint64) uint64 {
	var r uint64
	err := db.GetDB().Get(&r, "SELECT MAX(height) from blocks WHERE coin = ? AND timestamp <= ?", coin.Symbol, ts)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	return r
}

func statGetBlockCount(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT COUNT(*) from blocks b"+restrict, coin.Symbol, ts, ts-interval)
	return r
}

func statGetAvgDifficulty(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT AVG(difficulty) from blocks b"+restrict, coin.Symbol, ts, ts-interval)
	return r
}

func statGetWork(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT MAX(work) from blocks WHERE coin = ? AND timestamp <= ?", coin.Symbol, ts)
	return r
}

func statGetTotalReward(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT SUM(reward) from blocks b"+restrict, coin.Symbol, ts, ts-interval)
	return r
}

func statGetAvgReward(coin Coin, ts, interval uint64) float64 {
	blocks := statGetBlockCount(coin, ts, interval)
	if blocks == 0 {
		return 0
	}
	return statGetTotalReward(coin, ts, interval) / blocks
}

func statGetAvgPrice(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT AVG(p.price) from blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height"+restrict, coin.Symbol, ts, ts-interval)
	return r
}

func statGetAvgDARI(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT AVG(dari) from blocks b"+restrict, coin.Symbol, ts, ts-interval)
	return r
}

func statGetAvgSize(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT AVG(size) from blocks b"+restrict+" AND b.txs > 1", coin.Symbol, ts, ts-interval)
	return r
}

func statGetAvgTxs(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT AVG(txs) from blocks b"+restrict+" AND b.txs > 1", coin.Symbol, ts, ts-interval)
	return r
}

func statGetAvgCDD(coin Coin, ts, interval uint64) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT AVG(cdd) from blocks b"+restrict+" AND b.txs > 1", coin.Symbol, ts, ts-interval)
	return r
}

func statGetRate(coin Coin, ts, interval uint64, span string) float64 {
	var r float64
	db.GetDB().Get(&r, "SELECT MAX("+span+") from blocks b"+restrict, coin.Symbol, ts, ts-interval)
	return r
}

func statGetRateCalulated(coin Coin, ts, span uint64) float64 {
	r := struct {
		Count      float64 `db:"blks"`
		Difficulty float64 `db:"diff"`
	}{}

	err := db.GetDB().Get(&r, "SELECT COUNT(*) as blks, AVG(difficulty) as diff from blocks b"+restrict, coin.Symbol, ts, ts-span)

	if r.Count == 0 || err != nil {
		return 0.0
	}

	expected := span / 600
	return r.Count / float64(expected) * r.Difficulty
}

func statGetRateP(coin Coin, ts, interval uint64) float64 {
	blocks := statGetBlockCount(coin, ts, interval)

	if blocks == 0 {
		return 0
	}

	var diff float64
	db.GetDB().Get(&diff, "SELECT AVG(difficulty) from blocks b"+restrict, coin.Symbol, ts, ts-interval)
	ratio := blocks / (float64(interval) / 600)

	return ratio * diff
}
