
<template>
    <div>
        <p>
            Average tx fee in satoshis per byte.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-txfee-sat" style="height: 512px;"></canvas>
        </div>
        <br>
        <p>
            <strong>Note:</strong> these statistics show fees for the
            <u>average</u> tx size. For regular transactions (with few inputs/outputs) the median tx size is a more useful statistic but that data is currently not available.
        </p>
        <br>
        <stat-table :stats="stats"></stat-table>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'block-txfee-sat',
    mounted() {
        this.redraw()
    },
    computed: {
        stats() {
            return {
                header: ['3h', '6h', '12h', '1d', '7d'],
                f: (coin, avg) => {
                    return this.txfee(this.data[coin].averages[avg])
                },
                coins: [{
                    coin: 'BTC',
                    values: ['3h', '6h', '12h', '1d', '7d']
                }, {
                    coin: 'BCH',
                    values: ['3h', '6h', '12h', '1d', '7d']
                }]
            }

        }
    },
    methods: {
        txfee(obj) {
            if (obj.blocks < 1 || obj.price < 1) return null

            return this.sprintf('%.2f', (obj.reward_avg - (12.5 * 100000000)) / obj.txs / 1000)
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.txfee(obj)),
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.txfee(obj)),
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
                            position: 'right',
                            scaleLabel: {
                                display: true,
                                labelString: 'Tx fee/byte (satoshis)'
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-txfee-sat', config, labels)
        }
    }
}
</script>

<style>

</style>
