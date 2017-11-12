package bitcoin

import (
	"backend/db"
	"time"
	"gopkg.in/guregu/null.v3"
)

type CoinSummary struct {
	Price struct {
		Now   float64 `json:"now"`
		Day1  float64 `json:"day_1"`
		Day7  float64 `json:"day_7"`
		Day30 float64 `json:"day_30"`
	} `json:"price"`

	Blocks struct {
		PerHour struct {
			Hour6 float64 `json:"hour_6"`
			Day1  float64 `json:"day_1"`
			Day14 float64 `json:"day_14"`
		} `json:"per_hour"`
		Hour1 int `json:"hour_1"`
		Hour6 int `json:"hour_6"`
		Day1  int `json:"day_1"`
		Day14 int `json:"day_14"`
	} `json:"blocks"`

	DARI struct {
		Last  float64 `json:"last"`
		Blk3  float64 `json:"blk_3"`
		Hour6 float64 `json:"hour_6"`
		Day1  float64 `json:"day_1"`
		Day7  float64 `json:"day_7"`
	} `json:"dari"`

	Reward struct {
		Block struct {
			Coinbase  float64 `json:"coinbase"`
			Fees      float64 `json:"fees"`
			FeesCoin  float64 `json:"fees_coin"`
			FeePct    float64 `json:"fee_pct"`
			Total     float64 `json:"total"`
			TotalCoin float64 `json:"total_coin"`
			Blocks    int     `json:"blocks"`
		} `json:"block"`

		Day struct {
			Coinbase  float64 `json:"coinbase"`
			Fees      float64 `json:"fees"`
			FeesCoin  float64 `json:"fees_coin"`
			FeePct    float64 `json:"fee_pct"`
			Total     float64 `json:"total"`
			TotalCoin float64 `json:"total_coin"`
			Blocks    int     `json:"blocks"`
		} `json:"day"`

		Week struct {
			Coinbase  float64 `json:"coinbase"`
			Fees      float64 `json:"fees"`
			FeesCoin  float64 `json:"fees_coin"`
			FeePct    float64 `json:"fee_pct"`
			Total     float64 `json:"total"`
			TotalCoin float64 `json:"total_coin"`
			Blocks    int     `json:"blocks"`
		} `json:"week"`
	} `json:"reward"`

	Txs struct {
		Second float64 `json:"second"`
		Hour   float64 `json:"hour"`
		Day    float64 `json:"day"`
		Block  float64 `json:"block"`
	} `json:"txs"`

	TxFees struct {
		SatB  float64 `json:"sat_b"`
		UsdKB float64 `json:"usd_kb"`
	} `json:"tx_fees"`

	// @todo halvening
	Halvening struct {
		Seconds   int64   `json:"seconds"`
		Blocks    int     `json:"blocks"`
		Reward    float64 `json:"reward"`
		RewardNew float64 `json:"reward_new"`
	} `json:"halvening"`

	Blocksize struct {
		Day1  float64 `json:"day_1"`
		Day7  float64 `json:"day_7"`
		Day30 float64 `json:"day_30"`
		Day90 float64 `json:"day_90"`
	} `json:"blocksize"`

	Hashrate struct {
		Hour6 float64 `json:"hour_6"`
		Day1  float64 `json:"day_1"`
		Day7  float64 `json:"day_7"`
		Day30 float64 `json:"day_30"`
	} `json:"hashrate"`

	FeeMarket    []*FeeMarket `json:"fee_market"`
	DARIRelative []null.Float `json:"dari_relative"`
}

type FeeMarket struct {
	Timestamp  int64   `json:"timestamp"`
	Fees       float64 `json:"fees"`
	Percentage float64 `json:"percentage"`
}

type RewardPrice struct {
	Reward float64 `db:"reward"`
	Price  float64 `db:"price"`
}

func MakeCoinSummary(coin Coin) *CoinSummary {
	summary := &CoinSummary{}

	getPriceSummary(coin, summary)
	getBlockSummary(coin, summary)
	getRewardSummary(coin, summary)
	getTxsSummary(coin, summary)
	getDARISummary(coin, summary)
	getHalveningSummary(coin, summary)
	getBlocksizeSummary(coin, summary)
	getFeeMarketSummary(coin, summary)
	getTxFeesSummary(coin, summary)
	getHashrateSummary(coin, summary)
	getDARIRelativeSummary(coin, summary)

	return summary
}

func getDARIRelativeSummary(coin Coin, summary *CoinSummary) {
	rdari := make([]null.Float, 0)
	from := time.Now().Truncate(1800 * time.Second).Unix() + 1800

	for n := 336; n >= 0; n-- {
		var dari float64
		t := from - (int64(n) * 1800)
		f := t - 1800

		if err := db.GetDB().Get(&dari, "SELECT AVG(b.dari) from blocks b WHERE b.coin = ? AND b.timestamp >= ? AND b.timestamp < ?", coin.Symbol, f, t); err != nil {
			rdari = append(rdari, null.NewFloat(0, false))
		} else {
			rdari = append(rdari, null.NewFloat(dari, true))
		}
	}

	summary.DARIRelative = rdari
}

func getHashrateSummary(coin Coin, summary *CoinSummary) {
	now := uint64(time.Now().Unix())
	summary.Hashrate.Hour6 = statGetRateCalulated(coin, now, 6*3600)
	summary.Hashrate.Day1 = statGetRateCalulated(coin, now, 24*3600)
	summary.Hashrate.Day7 = statGetRateCalulated(coin, now, 7*24*3600)
	summary.Hashrate.Day30 = statGetRateCalulated(coin, now, 30*24*3600)
}

func getTxFeesSummary(coin Coin, summary *CoinSummary) {
	db.GetDB().Get(&(summary.TxFees.UsdKB), "SELECT (AVG(b.reward)-1250000000) / (AVG(b.size)) * 1000 / 100000000 * AVG(p.price) FROM blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp >= ?", coin.Symbol, time.Now().Unix()-(24*3600))
	db.GetDB().Get(&(summary.TxFees.SatB), "SELECT ((AVG(b.reward-1250000000)) / AVG(b.size)) FROM blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp >= ?", coin.Symbol, time.Now().Unix()-(24*3600))
}

func getFeeMarketSummary(coin Coin, summary *CoinSummary) {
	feemarket := make([]*FeeMarket, 0)
	from := time.Now().Truncate(24 * time.Hour).Unix()

	for n := 30; n >= 0; n-- {
		results := make([]RewardPrice, 0)

		f := from - (int64(n) * 24 * 3600)
		t := f + (24 * 3600)

		db.GetDB().Select(&results, "SELECT b.reward AS reward, p.price as price FROM blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp >= ? AND b.timestamp < ?", coin.Symbol, f, t)

		var fees float64 = 0
		var reward float64 = 0
		var sum float64 = 0

		for _, r := range results {
			fee := r.Reward/100000000 - 12.5
			fees += fee
			reward += r.Reward / 100000000
			sum += fee * r.Price
		}

		feemarket = append(feemarket, &FeeMarket{Timestamp: f, Fees: sum, Percentage: fees / reward * 100})
	}

	summary.FeeMarket = feemarket
}

func getBlocksizeSummary(coin Coin, summary *CoinSummary) {
	db.GetDB().Get(&(summary.Blocksize.Day1), "SELECT AVG(size)/1000 FROM blocks WHERE coin = ? AND timestamp > ? AND txs > 1", coin.Symbol, time.Now().Unix()-(24*3600))
	db.GetDB().Get(&(summary.Blocksize.Day7), "SELECT AVG(size)/1000 FROM blocks WHERE coin = ? AND timestamp > ? AND timestamp < ? AND txs > 1", coin.Symbol, time.Now().Unix()-(8*24*3600), time.Now().Unix()-(7*24*3600))
	db.GetDB().Get(&(summary.Blocksize.Day30), "SELECT AVG(size)/1000 FROM blocks WHERE coin = ? AND timestamp > ? AND timestamp < ? AND txs > 1", coin.Symbol, time.Now().Unix()-(31*24*3600), time.Now().Unix()-(30*24*3600))
	db.GetDB().Get(&(summary.Blocksize.Day90), "SELECT AVG(size)/1000 FROM blocks WHERE coin = ? AND timestamp > ? AND timestamp < ? AND txs > 1", coin.Symbol, time.Now().Unix()-(91*24*3600), time.Now().Unix()-(90*24*3600))
}

func getHalveningSummary(coin Coin, summary *CoinSummary) {
	height := getHeight(coin)
	summary.Halvening.Blocks = 630000 - height

	var d30blocks float64
	var reward float64
	db.GetDB().Get(&d30blocks, "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(30*24*3600))
	db.GetDB().Get(&reward, "SELECT AVG(reward) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(30*24*3600))

	summary.Halvening.Seconds = int64(float64(summary.Halvening.Blocks) / d30blocks * (30 * 24 * 3600))

	//@todo halvening
	summary.Halvening.RewardNew = reward/100000000 - 6.25
	summary.Halvening.Reward = reward / 100000000
}

func getDARISummary(coin Coin, summary *CoinSummary) {
	height := getHeight(coin)
	db.GetDB().Get(&(summary.DARI.Last), "SELECT dari FROM blocks WHERE coin = ? AND height = ?", coin.Symbol, height)
	db.GetDB().Get(&(summary.DARI.Blk3), "SELECT AVG(dari) FROM blocks WHERE coin = ? AND height >= ?", coin.Symbol, height-2)
	db.GetDB().Get(&(summary.DARI.Hour6), "SELECT AVG(dari) FROM blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(6*3600))
	db.GetDB().Get(&(summary.DARI.Day1), "SELECT AVG(dari) FROM blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(24*3600))
	db.GetDB().Get(&(summary.DARI.Day7), "SELECT AVG(dari) FROM blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(7*24*3600))
}

func getTxsSummary(coin Coin, summary *CoinSummary) {
	db.GetDB().Get(&(summary.Txs.Day), "SELECT SUM(txs) FROM blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(24*3600))

	summary.Txs.Hour = summary.Txs.Day / 24
	summary.Txs.Second = summary.Txs.Day / (24 * 3600)

	db.GetDB().Get(&(summary.Txs.Block), "SELECT AVG(txs) FROM blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(24*3600))
}

func getRewardSummary(coin Coin, summary *CoinSummary) {
	height := getHeight(coin)
	price := getPrice(coin, height)

	var reward float64
	db.GetDB().Get(&reward, "SELECT reward/100000000 FROM blocks WHERE coin = ? AND height = ?", coin.Symbol, height)

	// @todo halvening
	summary.Reward.Block.Coinbase = 12.5 * price
	summary.Reward.Block.Fees = (reward - 12.5) * price
	summary.Reward.Block.Total = reward * price
	summary.Reward.Block.TotalCoin = reward
	summary.Reward.Block.FeePct = summary.Reward.Block.Fees / summary.Reward.Block.Total * 100
	summary.Reward.Block.FeesCoin = reward - 12.5
	summary.Reward.Block.Blocks = 1

	getRewardInPeriod(
		coin,
		time.Now().Unix()-(24*3600),
		time.Now().Unix(),
		&(summary.Reward.Day.Coinbase),
		&(summary.Reward.Day.Fees),
		&(summary.Reward.Day.Total),
		&(summary.Reward.Day.TotalCoin),
		&(summary.Reward.Day.FeePct),
		&(summary.Reward.Day.FeesCoin),
		&(summary.Reward.Day.Blocks),
	)

	getRewardInPeriod(
		coin,
		time.Now().Unix()-(7*24*3600),
		time.Now().Unix(),
		&(summary.Reward.Week.Coinbase),
		&(summary.Reward.Week.Fees),
		&(summary.Reward.Week.Total),
		&(summary.Reward.Week.TotalCoin),
		&(summary.Reward.Week.FeePct),
		&(summary.Reward.Week.FeesCoin),
		&(summary.Reward.Week.Blocks),
	)
}

func getRewardInPeriod(coin Coin, from, to int64, coinbase, fees, total, totalcoin, feepct, feescoin *float64, blocks *int) {
	results := make([]RewardPrice, 0)
	db.GetDB().Select(&results, "SELECT b.reward/100000000 AS reward, p.price as price FROM blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp >= ? AND b.timestamp < ?", coin.Symbol, from, to)

	for _, rp := range results {
		*coinbase += 12.5 * rp.Price
		*fees += (rp.Reward - 12.5) * rp.Price
		*total += rp.Reward * rp.Price
		*feescoin += rp.Reward - 12.5
		*totalcoin += rp.Reward
	}

	*feepct = *feescoin / *totalcoin * 100
	*blocks = len(results)
}

func getBlockSummary(coin Coin, summary *CoinSummary) {
	db.GetDB().Get(&(summary.Blocks.PerHour.Hour6), "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(6*3600))
	db.GetDB().Get(&(summary.Blocks.PerHour.Day1), "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(24*3600))
	db.GetDB().Get(&(summary.Blocks.PerHour.Day14), "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(14*24*3600))

	summary.Blocks.PerHour.Hour6 /= 6
	summary.Blocks.PerHour.Day1 /= 24
	summary.Blocks.PerHour.Day14 /= (14 * 24)

	db.GetDB().Get(&(summary.Blocks.Hour1), "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-3600)
	db.GetDB().Get(&(summary.Blocks.Hour6), "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(6*3600))
	db.GetDB().Get(&(summary.Blocks.Day1), "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(24*3600))
	db.GetDB().Get(&(summary.Blocks.Day14), "SELECT COUNT(*) from blocks WHERE coin = ? AND timestamp >= ?", coin.Symbol, time.Now().Unix()-(14*24*3600))
}

func getPriceSummary(coin Coin, summary *CoinSummary) {
/*	db.GetDB().Get(&(summary.Price.Now), "SELECT price from prices WHERE coin = ? ORDER BY height DESC LIMIT 1", coin.Symbol)*/
	summary.Price.Now = ExchangeRates[coin.Symbol]

	db.GetDB().Get(&(summary.Price.Day1),
		"SELECT p.price from blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp < ? ORDER BY b.height DESC LIMIT 1",
		coin.Symbol, time.Now().Unix()-(24*3600))

	db.GetDB().Get(&(summary.Price.Day7),
		"SELECT p.price from blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp < ? ORDER BY b.height DESC LIMIT 1",
		coin.Symbol, time.Now().Unix()-(7*24*3600))

	db.GetDB().Get(&(summary.Price.Day30),
		"SELECT p.price from blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp < ? ORDER BY b.height DESC LIMIT 1",
		coin.Symbol, time.Now().Unix()-(30*24*3600))
}

func getHeight(coin Coin) int {
	var max int
	db.GetDB().Get(&max, "SELECT MAX(height) FROM blocks WHERE coin = ?", coin.Symbol)
	return max
}

func getPrice(coin Coin, height int) float64 {
	var price float64
	db.GetDB().Get(&price, "SELECT price FROM prices WHERE coin = ? AND height = ?", coin.Symbol, height)
	return price
}

func getAvgPrice(coin Coin, from, to int64) float64 {
	var price float64
	db.GetDB().Get(&price, "SELECT AVG(p.price) FROM blocks b JOIN prices p on b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp >= ? AND b.timestamp <= ?", coin.Symbol, from, to)

	return price
}
