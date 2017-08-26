import numeral from 'numeral'
import {sprintf} from 'sprintf-js'

export default {
    computed: {
        data() {
            return this.$store.getters['data']
        },
        loaded() {
            return this.$store.getters['loaded']
        },
        error() {
            return this.$store.getters['error']
        },
        leader() {
            return this.coin_leader()
        },
        loser() {
            return this.coin_leader(true)
        },
        leader_ratio() {
            return this.sprintf('%.2f', this.data[this.leader].averages['last'].dari / this.data[this.loser].averages['last'].dari)
        },
        leader_ratio_float() {
            return this.data[this.leader].averages['last'].dari / this.data[this.loser].averages['last'].dari
        },
        btc_c() {
            return 'rgba(255, 136, 0, 1)'
        },
        btc_cbg() {
            return 'rgba(255, 136, 0, 0.1)'
        },
        btc_cbge() {
            return 'rgba(255, 136, 0, 0.2)'
        },

        bch_c() {
            return 'rgba(0, 0, 255, 1.0)'
        },
        bch_cbg() {
            return 'rgba(0, 0, 255, 0.1)'
        },
        bch_cbge() {
            return 'rgba(0, 0, 255, 0.2)'
        },
        preforktarget() {
            return {
                height: 477792,
                difficulty: 860221984436.22,
                timestamp: 1501153380
            }
        },
        btc_h() {
            return this.$store.getters['btc_h']
        },
        bch_h() {
            return this.$store.getters['bch_h']
        },
        split_height() {
            return 478558
        },
        all_timeframe() {
            return this.data.BTC.history.all[1].timestamp - this.data.BTC.history.all[0].timestamp
        }
    },
    methods: {
        sprintf(fmt, val) {
            return sprintf(fmt, val)
        },
        numeral(fmt, val) {
            return numeral(val).format(fmt)
        },
        money(val) {
            return this.numeral('$0,0.00', val)
        },
        coin_leader(r = false) {
            return this.data.BTC.averages.last.dari > this.data.BCH.averages.last.dari ? (r ? "BCH" : "BTC") : (r ? "BTC" : "BCH")
        },
        dari(coin) {
            return sprintf('%.2f', this.data[coin].averages.last.dari)
        },
        fees(obj){
            return obj.blocks > 0 && obj.price > 1 ? this.sprintf('%.2f', (obj.reward_avg / 100000000 - 12.5) * obj.price) : null
        },
        lastfee(coin) {
            return (this.data[coin].averages.last.reward_avg - 1250000000) / 100000000 * this.data[coin].averages.last.price
        },
        ts_coin(coin, ts) {
            return this.data[coin].history['all'].find(h => h.timestamp === ts)
        },
        next_pct(coin, factor = false) {
            let prev = this.data[coin].retargets.slice(-1)[0]
            let last = this.data[coin].averages.last

            let mod = last.height % 2016
            if(mod === 0 || mod === 1) return null

            let mins = this.estimate(coin, last, true)

            if(mins === null) return null

            let sec = (last.timestamp + mins * 60) - prev.timestamp
            let normal = 2016 * 600

            let exp = ((normal / sec)) * 100

            if(factor) {
                return exp > 300 ? 0.25 : (normal / sec)
            }

            return exp > 300 ? 300 : (exp < 100 ? exp - 100: exp)
        },
    }
}