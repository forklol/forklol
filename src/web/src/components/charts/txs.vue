
<template>
    <div>
        <p>
            Average number of transactions per block.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-txs" style="height: 512px;"></canvas>
        </div>
        <br>
        <stat-table :stats="stats"></stat-table>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'block-txs',
    mounted() {
        this.redraw()
    }, 
    computed: {
        stats() {
            return {
                header: ['3h', '6h', '12h', '1d', '7d'],
                f: (coin, avg) => {
                    return this.txs(this.data[coin].averages[avg])
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
        txs(obj) {
            if(obj.blocks < 1) return null
            return obj.blocks > 0 ? this.sprintf('%.2f', obj.txs) : null
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.txs(obj)),
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.txs(obj)),
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
                                labelString: 'Transactions'
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-txs', config, labels)
        }
    }
}
</script>

<style>

</style>
