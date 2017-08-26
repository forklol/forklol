
<template>
    <div>
        <p>
            Fee percentage contributing to total reward (calculated in USD).
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-feepct" style="height: 512px;"></canvas>
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
    name: 'blocks-feepct',
    mounted() {
        this.redraw()
    },
    computed: {
        stats() {
            return {
                header: ['3h', '6h', '12h', '1d', '7d'],
                f: (coin, avg) => {
                    let p = this.feepct(this.data[coin].averages[avg])
                    if(p === null) return null
                    return  p + '%'
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
        feepct(obj) {
            if (obj.blocks < 1 || obj.price < 1) {
                return null
            }

            let reward = obj.reward_avg / 100000000 * obj.price
            let fee = (obj.reward_avg / 100000000 - 12.5) * obj.price

            return this.sprintf('%.2f', fee / reward * 100)
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                backgroundColor: this.btc_cbg,
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.feepct(obj)),
                fill: true,
            })

            this.add_dataset({
                label: 'BCH',
                backgroundColor: this.bch_cbg,
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.feepct(obj)),
                fill: true,
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
                                labelString: 'Fee %'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return value + '%'
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-feepct', config, labels)
        },
    }
}
</script>

<style>

</style>
