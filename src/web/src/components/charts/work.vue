
<template>
    <div>
        <p>
            Total work done, or in other words, the sum of all block difficulties.
            <div style="height:100%; width: 100%;">
                <canvas id="canvas-work" style="height: 512px;"></canvas>
            </div>
            <br>
            <p>
                Numbers are in trillions (1,000,000,000,000x).
            </p>
            <small>
                Just for fun, the approximate total number of SHA256 hashes needed to produce the work: {{ hashes }}.
                <br>
                Protip: memorize this number to impress your friends and family, they love to hear about it. No joke! ;)            
            </small>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'
import words from 'written-number'
import big from 'big-integer'

export default {
    name: 'pow-work',
    mounted() {
        this.redraw()
    },
    computed: {
        hashes() {
            return big(this.data.BTC.averages.last.work).multiply(Math.pow(2, 32)).toString()
        },
    },
    methods: {
        work(obj) {
            if(obj.work === 0) return null
            return this.sprintf('%.2f', obj.work / 1000000000000)
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.work(obj))
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.work(obj))
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
                                labelString: 'Total work (trillions)'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return this.numeral('0,0', value)
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-work', config, labels)
        }
    }
}
</script>

<style>

</style>
