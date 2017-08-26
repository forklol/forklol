package bitcoin

import (
	"backend/db"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var api_pubkey, api_secret string

func SetBitcoinAverageApiCreds(pub, secret string) {
	api_pubkey = pub
	api_secret = secret
}

func fetchBtcAvg(url string) ([]byte, error) {
	payload := strconv.Itoa(int(time.Now().Unix())) + "." + api_pubkey

	mac := hmac.New(sha256.New, []byte(api_secret))
	mac.Write([]byte(payload))
	digest := hex.EncodeToString(mac.Sum(nil))
	sig := payload + "." + digest

	c := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-signature", sig)

	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, nil
}

func GetPriceDb(coin Coin, height uint32, timestamp int64) (float64, error) {

	if height <= CHAINSPLIT_HEIGHT {
		return 0.0, nil
	}

	p := struct {
		Price float64 `db:"price"`
	}{}
	if err := db.GetDB().Get(&p, "SELECT price FROM prices WHERE coin = ? AND height = ?", coin.Symbol, height); err != nil {
		log.Printf("No db price found, fetching from BitcoinAverage.com\n")

		a, err := fetchHistoricalPrice(coin, int64(timestamp))
		if err != nil {
			log.Printf("No price could be fetched for %s @ %d\n", coin.Symbol, height)
			return 0.0, err
		}

		setPriceDb(coin, height, a)

		p.Price = a
	}

	return p.Price, nil
}

func setPriceDb(coin Coin, height uint32, price float64) {
	db.GetDB().Exec("INSERT INTO prices (coin, height, price) VALUES(?,?,?)", coin.Symbol, height, price)
}

func fetchPrice(coin Coin) (float64, error) {
	url := fmt.Sprintf("https://apiv2.bitcoinaverage.com/indices/global/ticker/%sUSD", coin.Symbol)
	body, err := fetchBtcAvg(url)

	data := struct {
		Last float64 `json:"last"`
	}{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0.0, errors.New("Could not decode price response")
	}

	return data.Last, nil
}

func fetchHistoricalPrice(coin Coin, timestamp int64) (float64, error) {
	url := fmt.Sprintf("https://apiv2.bitcoinaverage.com/indices/global/history/%sUSD?at=%d", coin.Symbol, timestamp)
	body, err := fetchBtcAvg(url)

	data := struct {
		Average float64 `json:"average"`
	}{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0.0, errors.New("Could not decode price response")
	}

	return data.Average, nil
}
