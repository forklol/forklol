package bitcoin

import (
	"backend/db"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BlockData struct {
	Height     uint32
	Date       time.Time
	Reward     uint64
	Size       float64
	Difficulty float64
	Txs        uint64
	CDD        float64
}

type Retarget struct {
	Height     uint64  `json:"height"`
	Difficulty float64 `json:"difficulty"`
	Timestamp  uint64  `json:"timestamp"`
}

type ChainAverage struct {
	Coin        string  `json:"-" db:"coin"`
	Height      uint64  `json:"height" db:"height"`
	Blocks      float64 `json:"blocks" db:"blocks"`
	AvgDiff     float64 `json:"avg_diff" db:"avg_diff"`
	Work        float64 `json:"work" db:"work"`
	Rate        float64 `json:"rate" db:"rate"`
	Rate1d      float64 `json:"rate1d" db:"rate1d"`
	Rate3d      float64 `json:"rate3d" db:"rate3d"`
	Rate7d      float64 `json:"rate7d" db:"rate7d"`
	Rate3h      float64 `json:"rate3h" db:"rate3h"`
	Rate12h     float64 `json:"rate12h" db:"rate12h"`
	TotalReward float64 `json:"reward_total" db:"reward_total"`
	AvgReward   float64 `json:"reward_avg" db:"reward_avg"`
	CoinPrice   float64 `json:"price" db:"price"`
	DARI        float64 `json:"dari" db:"dari"`
	Size        float64 `json:"size" db:"size"`
	Txs         float64 `json:"txs" db:"txs"`
	CDD         float64 `json:"cdd" db:"cdd"`
	Timestamp   uint64  `json:"timestamp" db:"timestamp"`
}

const CHAINSPLIT_TIMESTAMP = 1501593374

const QRY_INSERT = `
    INSERT INTO blocks (
        coin, 
        height, 
        difficulty, 
        work,
        rate1d,
        rate3d,
        rate7d,
        rate3h,
        rate12h,
        reward, 
        dari,
        size,
        txs,
        cdd,
        timestamp
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

const QRY_SELECT = `
    SELECT
        b.coin as coin, 
        max(b.height) as height,
        count(*) as blocks, 
        avg(b.difficulty) as avg_diff, 
        avg(b.size) as size, 
        max(b.work) as work,
        avg(b.rate1d) as rate,
        avg(b.rate1d) as rate1d,
        avg(b.rate3d) as rate3d,
        avg(b.rate7d) as rate7d,
        avg(b.rate3h) as rate3h,
        avg(b.rate12h) as rate12h,
        sum(b.reward) as reward_total, 
        avg(b.reward) as reward_avg, 
        avg(p.price) as price, 
        avg(b.dari) as dari, 
        avg(b.txs) as txs, 
        avg(b.cdd) as cdd, 
        round(avg(b.timestamp)) as timestamp
    FROM
        blocks b
    JOIN 
        prices p 
    ON
        b.coin = p.coin AND b.height = p.height
    WHERE 
        b.coin = ?
    AND 
        b.timestamp > ? 
    AND 
        b.height > 478558
    GROUP BY 
        b.coin
`

const QRY_RETARGET = `
    SELECT 
        height,
        difficulty,
        timestamp
    FROM 
        blocks 
    WHERE 
        height % 2016 = 0 
    AND
        coin = ?
    ORDER BY 
        height ASC;
`

func GetLastBlockHeight(coin Coin) uint64 {
	max := struct {
		Max uint64 `db:"max"`
	}{}

	if err := db.GetDB().Get(&max, "SELECT max(height) as max FROM blocks WHERE coin = ?", coin.Symbol); err != nil {
		return 478000
		log.Fatalln("Could not get max block height for " + coin.Symbol)
	}

	return max.Max
}

func GetLastBlockAverages(coin Coin) *ChainAverage {
	ts := struct {
		Timestamp uint64 `db:"timestamp"`
	}{}

	if err := db.GetDB().Get(&ts, "SELECT timestamp FROM blocks WHERE coin = ? ORDER BY height DESC LIMIT 1", coin.Symbol); err != nil {
		log.Fatalln("Could not get latest block for " + coin.Symbol)
	}

	return GetBlockAverages(coin, ts.Timestamp-1)
}

func SyncBlocks(coin Coin) (bool, error) {
	return importBlocks(GetLastBlockHeight(coin)+1, coin)
}

func importBlocks(from uint64, coin Coin) (bool, error) {
	url := fmt.Sprintf("https://api.blockchair.com/%s/blocks?q=id(%d..)&s=id(asc)&fields=id,time,reward,size,difficulty,transaction_count,cdd_total&export=csv", coin.BlockchairURL, from)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var dummy map[string]interface{}
	if err := json.Unmarshal(data, &dummy); err == nil {
		log.Println("No new blocks found for " + coin.Symbol)
		return false, nil
	}

	csv := csv.NewReader(strings.NewReader(string(data)))

	blocks := make([]BlockData, 0, 256)
	start := true
	for {
		record, err := csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}

		if start {
			start = false
			continue
		}

		if len(record) < 2 {
			return false, errors.New("No records found")
		}

		if blocktime, err := time.Parse("2006-01-02 15:04:05", record[1]); err == nil {
			height, _ := strconv.ParseUint(record[0], 10, 32)
			reward, _ := strconv.ParseUint(record[2], 10, 64)
			size, _ := strconv.ParseFloat(record[3], 64)
			diff, _ := strconv.ParseFloat(record[4], 64)
			txs, _ := strconv.ParseUint(record[5], 10, 64)
			cdd, _ := strconv.ParseFloat(record[6], 64)

			blocks = append(blocks, BlockData{
				Height:     uint32(height),
				Date:       blocktime,
				Reward:     reward,
				Size:       size,
				Difficulty: diff,
				Txs:        txs,
				CDD:        cdd,
			})
		} else {
			return false, err
		}
	}

	work, err := GetWorkTotal(coin)
	if err != nil {
		log.Println("Could not get total work for " + coin.Symbol)
		return false, err
	}

	for _, block := range blocks {

		if block.Height > CHAINSPLIT_HEIGHT {
			work.Work += block.Difficulty
		}

		then := uint64(block.Date.Unix())

		price, err := GetPriceDb(coin, block.Height, block.Date.Unix())
		if err != nil {
			log.Printf("No price known for %s block %d @ %s\n", coin.Symbol, block.Height, block.Date)
			return false, err
		}

		log.Printf("Imported %s price of %.2f for block %d @ %s\n", coin.Symbol, price, block.Height, block.Date)
		index := float64(block.Reward) / block.Difficulty * price
		_, err = db.GetDB().Exec(QRY_INSERT,
			coin.Symbol,
			block.Height,
			block.Difficulty,
			work.Work,
			getHashRatePeriod(coin, then-(1*24*3600), then),
			getHashRatePeriod(coin, then-(3*24*3600), then),
			getHashRatePeriod(coin, then-(7*24*3600), then),
			getHashRatePeriod(coin, then-(3*3600), then),
			getHashRatePeriod(coin, then-(12*3600), then),
			block.Reward,
			index,
			block.Size,
			block.Txs,
			block.CDD,
			block.Date.Unix())

		if err != nil {
			log.Printf("Could not insert block: %s\n", err.Error())
			return false, err
		}
	}

	return true, nil
}

func getHashRatePeriod(coin Coin, from, to uint64) float64 {
	if coin.Symbol == "BCH" && from < 1501593374 {
		from = 1501593374
	} else if coin.Symbol == "BTC" && from < 1501268136 {
		from = 1501268136
	}

	f, err := GetHashFactorInPeriod(coin, from, to)

	if err != nil {
		return 0.0
	}

	return f.ActualWork
}

func GetHistoricalData(coin Coin, timestamp, interval uint64) *[]*ChainAverage {
	data, err := fetchBlocksGrouped(coin, timestamp, interval)

	if err != nil {
		log.Printf("Could not get historical data %s\n", err.Error())
	}

	return data
}

func GetRetargets(coin Coin) (*[]Retarget, error) {
	retargets := []Retarget{}
	if err := db.GetDB().Select(&retargets, QRY_RETARGET, coin.Symbol); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &retargets, nil
}

func fetchBlocksGrouped(coin Coin, timestamp, interval uint64) (*[]*ChainAverage, error) {
	blocks := []Block{}
	if err := db.GetDB().Select(&blocks, "SELECT b.*, p.price as price FROM blocks b JOIN prices p ON b.coin = p.coin AND b.height = p.height WHERE b.coin = ? AND b.timestamp > ?", coin.Symbol, timestamp); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return groupBlocks(coin, &blocks, timestamp, interval)
}

func groupBlocks(coin Coin, values *[]Block, start uint64, interval uint64) (*[]*ChainAverage, error) {
	ts := start - (start % interval) + interval

	groups := make([]*ChainAverage, 0, 128)
	group := make([]*Block, 0, 32)

	idx := 0
	now := uint64(time.Now().Unix())

	for n := ts; n < now; n += interval {
		for i := idx; i < len(*values); i++ {
			if (*values)[i].Timestamp < n {
				group = append(group, &(*values)[i])
				idx++
			} else {
				break
			}
		}
		groups = append(groups, calculateGroup(&group, n))
		group = make([]*Block, 0, 32)
	}

	if idx < len(*values) {
		for i := idx; i < len(*values); i++ {
			group = append(group, &(*values)[i])
			idx++
		}
		groups = append(groups, calculateGroup(&group, now))
	}

	linearFill(&groups)

	return &groups, nil
}

func linearFill(data *[]*ChainAverage) {
	prev := 0

	for n := 0; n < len(*data); n++ {
		if (*data)[n].Blocks > 0 {
			if prev < n-1 && prev > 0 {
				// we've got a gap ho ho ho. fill 'er up boys!
				for c := 0; c < n-prev; c++ {
					diff := float64(n - prev)
					p := (*data)[prev]

					(*data)[prev+c].Work = p.Work - (p.Work-(*data)[n].Work)/diff*(float64(c)+1)
					(*data)[prev+c].Rate1d = p.Rate1d - (p.Rate1d-(*data)[n].Rate1d)/diff*(float64(c)+1)
					(*data)[prev+c].Rate3d = p.Rate3d - (p.Rate3d-(*data)[n].Rate3d)/diff*(float64(c)+1)
					(*data)[prev+c].Rate7d = p.Rate7d - (p.Rate7d-(*data)[n].Rate7d)/diff*(float64(c)+1)
					(*data)[prev+c].Rate3h = p.Rate3h - (p.Rate3h-(*data)[n].Rate3h)/diff*(float64(c)+1)
					(*data)[prev+c].Rate12h = p.Rate12h - (p.Rate12h-(*data)[n].Rate12h)/diff*(float64(c)+1)
					(*data)[prev+c].AvgDiff = p.AvgDiff - (p.AvgDiff-(*data)[n].AvgDiff)/diff*(float64(c)+1)
					(*data)[prev+c].AvgReward = p.AvgReward - (p.AvgReward-(*data)[n].AvgReward)/diff*(float64(c)+1)
					(*data)[prev+c].Blocks = p.Blocks - (p.Blocks-(*data)[n].Blocks)/diff*(float64(c)+1)
					(*data)[prev+c].CoinPrice = p.CoinPrice - (p.CoinPrice-(*data)[n].CoinPrice)/diff*(float64(c)+1)
					(*data)[prev+c].DARI = p.DARI - (p.DARI-(*data)[n].DARI)/diff*(float64(c)+1)
					(*data)[prev+c].TotalReward = p.TotalReward - (p.TotalReward-(*data)[n].TotalReward)/diff*(float64(c)+1)
					(*data)[prev+c].Size = p.Size - (p.Size-(*data)[n].Size)/diff*(float64(c)+1)
					(*data)[prev+c].Txs = p.Txs - (p.Txs-(*data)[n].Txs)/diff*(float64(c)+1)
					(*data)[prev+c].CDD = p.CDD - (p.CDD-(*data)[n].CDD)/diff*(float64(c)+1)
				}
			}
			prev = n
		}
	}
}

func calculateGroup(group *[]*Block, timestamp uint64) *ChainAverage {
	blocks := len(*group)

	if blocks == 0 {
		return &ChainAverage{Timestamp: timestamp}
	}

	var difficulty, reward, price, dari, work, rate, rate3, rate7, rate3h, rate12h, size, txs, cdd float64
	var height uint64

	for _, blk := range *group {
		difficulty += blk.Difficulty
		reward += blk.Reward
		price += blk.Price
		dari += blk.DARI
		work += blk.Work
		rate += blk.Rate1d
		rate3 += blk.Rate3d
		rate7 += blk.Rate7d
		rate3h += blk.Rate3h
		rate12h += blk.Rate12h
		size += blk.Size
		txs += blk.Txs
		cdd += blk.CDD

		if blk.Height > height {
			height = blk.Height
		}
	}

	f := float64(blocks)

	return &ChainAverage{
		Height:      height,
		Work:        work / f,
		Rate:        rate / f,
		Rate1d:      rate / f,
		Rate3d:      rate3 / f,
		Rate7d:      rate7 / f,
		Rate3h:      rate3h / f,
		Rate12h:     rate12h / f,
		AvgDiff:     difficulty / f,
		AvgReward:   reward / f,
		Blocks:      f,
		Coin:        "",
		CoinPrice:   price / f,
		DARI:        dari / f,
		Timestamp:   timestamp,
		TotalReward: reward,
		Size:        size / f,
		Txs:         txs / f,
		CDD:         cdd / f,
	}
}

func GetBlockAverages(coin Coin, timestamp uint64) *ChainAverage {
	data := ChainAverage{}
	if err := db.GetDB().Get(&data, QRY_SELECT, coin.Symbol, timestamp); err != nil {
		data.Blocks = 0
	}

	return &data
}
