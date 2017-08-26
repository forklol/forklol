
<template>
    <div>
        <p>
            Total inflation since fork and inflation rate.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-inflation" style="height: 512px;"></canvas>
        </div>
        <br>
        <p>
            The chart is composed of 72 chunks of data since the moment the fork happened. The inflation rate is therefore calculated over a 1/72 timeframe. The total inflation is calculated over the complete dataset.
        </p>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'blocks-inflation',
    data() {
        return {
                previ: {
                    BTC: null,
                    BCH: null
                }
        }
    },
    mounted() {
        this.redraw()
    },
    methods: {
        inflation_rate(obj) {
            
            let coins = 16406250;
            coins += (obj.height - 472500) * 12.5

            let newcoins = (coins + (obj.blocks * 12.5))
            let inflation = (newcoins - coins) / newcoins * 100

            return this.sprintf('%.4f', inflation)
        },
        inflation_rate_fork(coin, obj) {

            if (obj.height == 0) {
                return this.previ[coin]
            }
 
            let coins = 16406250;
            coins += (obj.height - 472500) * 12.5
            let sincefork = (obj.height - 478559) * 12.5
            let prefork = coins-sincefork

            let diff = (coins - prefork) / prefork * 100

            this.previ[coin] = diff


            return this.sprintf('%.4f', diff)
        },
        redraw() {
            this.previ.BTC = null
            this.previ.BCH = null


            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))
            labels.pop()

            let btc = this.data.BTC.history['all'].map((obj) => this.inflation_rate_fork('BTC', obj))
            btc.pop()
            let bch = this.data.BCH.history['all'].map((obj) => this.inflation_rate_fork('BCH', obj))
            bch.pop()

            let btcr = this.data.BTC.history['all'].map((obj) => this.inflation_rate(obj))
            btcr.pop()
            let bchr = this.data.BCH.history['all'].map((obj) => this.inflation_rate(obj))
            bchr.pop()

            this.add_dataset({
                label: 'BTC total',
                backgroundColor: this.btc_cbg,
                borderColor: this.btc_c,
                data: btc,
                yAxisID: 'inflation'
            })
            this.add_dataset({
                label: 'BTC rate',
                backgroundColor: this.btc_cbg,
                borderColor: this.btc_c,
                data: btcr,
                borderDash: [3, 3],
                yAxisID: 'rate'
            })

            this.add_dataset({
                label: 'BCH total',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                data: bch,
                yAxisID: 'inflation'
            })

            this.add_dataset({
                label: 'BCH rate',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                data: bchr,
                borderDash: [3, 3],
                yAxisID: 'rate'
            })

            var config = {
                type: 'line',
                options: {
                    scales: {
                        xAxes: [{
                            display: true,
                            scaleLabel: {
                                display: true,
                                labelString: 'Date/Time'
                            }
                        }],
                        yAxes: [{
                            display: true,
                            id: 'inflation',
                            position: 'right',
                            scaleLabel: {
                                display: true,
                                labelString: 'Total inflation'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return this.sprintf('%.2f', value) + '%'
                                }
                            }
                        },{
                            display: true,
                            id: 'rate',
                            position: 'right',
                            scaleLabel: {
                                display: true,
                                labelString: 'Inflation rate'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return this.sprintf('%.3f', value) + '%'
                                }
                            }
                        }
                        ]
                    }
                }
            };

            this.draw('canvas-inflation', config, labels)
        },
    }
}
</script>

<style>

</style>
