
<template>
    <div>
        <p>
            Combined (sum of chain difficulties) expected difficulty based on hashrate versus actual difficulty.
        </p>
        <div style="height:512px; width: 100%;">
            <canvas id="canvas-difficulty-expected" style="height: 512px;"></canvas>
        </div>
        <br>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'pow-difficulty-expected',
    mounted() {
        this.redraw()
    },
    methods: {
        expected(obj){
            let a = this.ts_coin('BTC', obj.timestamp)
            let b = this.ts_coin('BCH', obj.timestamp)


            if(!a || !b) return null

            if(b.height < 478600) return null

            return a.rate7d + b.rate7d
        },
        actual(obj) {
            let a = this.ts_coin('BTC', obj.timestamp)
            let b = this.ts_coin('BCH', obj.timestamp)

            if(!a || !b) return null

            if(b.height < 478600) return null

            return a.avg_diff + b.avg_diff
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'Expected',
                borderColor: '#000000',
                spanGaps: true,
                data: this.data.BTC.history['all'].map((obj) => this.expected(obj)),
            })

            this.add_dataset({
                label: 'Actual',
                borderColor: '#ff0000',
                spanGaps: true,
                data: this.data.BTC.history['all'].map((obj) => this.actual(obj)),
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
                                labelString: 'Difficulty'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return this.numeral('0,0', value)
                                },
                                min: 800000000000
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-difficulty-expected', config, labels)
        }
    }
}
</script>

<style>

</style>
