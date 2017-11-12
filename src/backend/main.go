package main

import (
	"backend/api"
	"backend/bitcoin"
	"backend/config"
	"backend/db"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type ChainApiResponse struct {
	History   map[string]*[]*bitcoin.ChainAverage `json:"history"`
	Averages  map[string]*bitcoin.ChainAverage    `json:"averages"`
	Pow       map[string]interface{}              `json:"pow"`
	Retargets *[]bitcoin.Retarget                 `json:"retargets"`
	Summary   *bitcoin.CoinSummary                `json:"summary"`
}

type BtcAvgResp struct {
	Last float64
}

var port, connstr string

type LastData struct {
	BTC, BCH uint64
}

func main() {
	Init()

	retries := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, v := range retries {
		if err := db.InitDB(connstr); err != nil {
			log.Printf("#%d) Could not connect to db: %+v\n", v, err)
			if v == 10 {
				log.Panicf("EXIT: %+v\n", err)
			} else {
				time.Sleep(time.Duration(v) * time.Second)
				continue
			}
		}

		break
	}

	coins := []bitcoin.Coin{bitcoin.Coin{"BCH", "bitcoin-cash"}, bitcoin.Coin{"BTC", "bitcoin"}}

	startApi(coins)
}

func Init() {
	env_pubkey, _ := os.LookupEnv("BTCAVG_PUBKEY")
	env_secret, _ := os.LookupEnv("BTCAVG_SECRET")

	env_dbuser, _ := os.LookupEnv("DB_USER")
	env_dbpass, _ := os.LookupEnv("DB_PASS")
	env_dbhost, _ := os.LookupEnv("DB_HOST")
	env_dbport, _ := os.LookupEnv("DB_PORT")
	env_dbscheme, _ := os.LookupEnv("DB_SCHEME")

	pub := flag.String("pubkey", env_pubkey, "bitcoinaverage.com api public key, defaults to env var BTCAVG_PUBKEY")
	sec := flag.String("secret", env_secret, "bitcoinaverage.com api secret, defaults to env var BTCAVG_SECRET")
	prt := flag.String("port", "8888", "port the api webserver listens on")
	dbg := flag.Bool("debug", false, "enable debugging")

	dbuser := flag.String("dbuser", env_dbuser, "mysql user")
	dbpass := flag.String("dbpass", env_dbpass, "mysql password")
	dbhost := flag.String("dbhost", env_dbhost, "mysql host/address")
	dbport := flag.String("dbport", env_dbport, "mysql port")
	dbscheme := flag.String("dbscheme", env_dbscheme, "mysql database name/scheme")

	flag.Parse()

	port = *prt
	config.OPTIONS.Debug = *dbg

	bitcoin.SetBitcoinAverageApiCreds(*pub, *sec)
	bitcoin.ExchangeRates = make(map[string]float64)

	connstr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbuser, *dbpass, *dbhost, *dbport, *dbscheme)
}

func startApi(coins []bitcoin.Coin) {
	c := make(chan bool)

	go fetchRoutine(coins, c)
	api.Serve(port, c)
}

func fetchRoutine(coins []bitcoin.Coin, c chan bool) {
	var interval time.Duration = 5
	if config.OPTIONS.Debug == true {
		interval = 15
	}

	ticker := time.NewTicker(time.Second * interval)
	tickerPrice := time.NewTicker(time.Second * interval)

	log.Println("Prefetching data and saving stats..")
	go bitcoin.UpdatePrices(coins)
	saveStats(coins, true)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			log.Println("Fetching new blocks..")
			if saveStats(coins, false) {
				c <- true
			}
			ticker = time.NewTicker(interval * time.Second)
		case <-tickerPrice.C:
			go bitcoin.UpdatePrices(coins)
		}
	}
}

func saveStats(coins []bitcoin.Coin, force bool) bool {
	start := time.Now()

	newBlocks := false

	for _, coin := range coins {
		nw, err := bitcoin.SyncBlocks(coin)
		if err != nil {
			fmt.Printf("Error syncing blocks: %+v\n", err)
			return false
		}
		newBlocks = newBlocks || nw
	}

	if !newBlocks && !force {
		log.Printf("No new blocks found.\n")
		return false
	}

	apiChains := make(map[string]ChainApiResponse)

	now := uint64(time.Now().Unix())
	day := uint64(24 * 3600)

	for _, coin := range coins {
		history := make(map[string]*[]*bitcoin.ChainAverage)
		averages := make(map[string]*bitcoin.ChainAverage)
		pow := make(map[string]interface{})

		retargets, err := bitcoin.GetRetargets(coin)

		if err != nil {
			log.Printf("Could not determine retargets, %s\n", err.Error())
			return false
		}

		apiChain := ChainApiResponse{history, averages, pow, retargets, bitcoin.MakeCoinSummary(coin)}

		factor, err := bitcoin.GetHashFactorInPeriod(coin, now-3*day, now)
		if err != nil {
			log.Printf("Could not determine HashFactorInPeriod 3d, %s\n", err.Error())
			return false
		}

		split, err := bitcoin.GetWorkSinceSplit(coin)
		if err != nil {
			log.Printf("Could not determine WorkSinceSplit, %s\n", err.Error())
			return false
		}

		total, err := bitcoin.GetWorkTotal(coin)
		if err != nil {
			log.Printf("Could not determine TotalWork, %s\n", err.Error())
			return false
		}

		pow["current"] = factor
		pow["split"] = split
		pow["total"] = total

		/*		averages["7d"] = bitcoin.GetBlockAverages(coin, now-(3600*24*7))
				averages["1d"] = bitcoin.GetBlockAverages(coin, now-(3600*24))
				averages["12h"] = bitcoin.GetBlockAverages(coin, now-(3600*12))
				averages["6h"] = bitcoin.GetBlockAverages(coin, now-(3600*6))
				averages["3h"] = bitcoin.GetBlockAverages(coin, now-(3600*3))
				averages["last"] = bitcoin.GetLastBlockAverages(coin)*/

		averages["7d"] = bitcoin.GetLastAverages(coin, now, 3600*24*7)
		averages["1d"] = bitcoin.GetLastAverages(coin, now, 3600*24)
		averages["12h"] = bitcoin.GetLastAverages(coin, now, 3600*12)
		averages["6h"] = bitcoin.GetLastAverages(coin, now, 3600*6)
		averages["3h"] = bitcoin.GetLastAverages(coin, now, 3600*3)
		averages["last"] = bitcoin.GetLastBlockAverages(coin)

		from := now - (30 * 24 * 3600)

		//history["all"] = bitcoin.GetHistoricalData(coin, bitcoin.CHAINSPLIT_TIMESTAMP, (now-bitcoin.CHAINSPLIT_TIMESTAMP)/256)
		//history["all"] = bitcoin.GetHistoricalData(coin, from, (now-from)/360)
		history["all"] = bitcoin.GetHistoricalDataDescending(coin, now, from, (now-from)/360)

		apiChains[coin.Symbol] = apiChain
	}

	data, err := json.Marshal(apiChains)
	if err != nil {
		log.Printf("Could not format json %s\n", err.Error())
		return false
	}

	/*	var pretty bytes.Buffer
		json.Indent(&pretty, data, "", "\t")*/
	ioutil.WriteFile(api.JSON_PATH, data, 0644)

	last := LastData{
		BTC: bitcoin.GetLastBlockHeight(bitcoin.Coin{"BTC", ""}),
		BCH: bitcoin.GetLastBlockHeight(bitcoin.Coin{"BCH", ""}),
	}

	bts, _ := json.Marshal(last)
	ioutil.WriteFile(api.JSON_PATH+".last", bts, 0644)

	log.Printf("Done saving stats in %s\n", time.Now().Sub(start))

	return newBlocks
}
