
<template>
    <div>
        <p>
            Absolute hashrate in exahashes per second (<u>12h averages</u>). 
        </p>
        <div style="height:512px; width: 100%;">
            <canvas id="canvas-hashrateabs" style="height: 512px;"></canvas>
        </div>
        <br>
        <small>
            <table class="table table-sm table-bordered">
                <thead class="thead-inverse">
                    <tr>
                        <th>Coin</th>
                        <th>3h</th>
                        <th>12h</th>
                        <th>1d</th>
                        <th>3d</th>
                        <th>7d</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>BTC</td>
                        <td>{{ abs('BTC', '3h') }}</td>
                        <td>{{ abs('BTC', '12h') }}</td>
                        <td>{{ abs('BTC', '1d') }}</td>
                        <td>{{ abs('BTC', '3d') }}</td>
                        <td>{{ abs('BTC', '7d') }}</td>
                    </tr>
                    <tr>
                        <td>BCH</td>
                        <td>{{ abs('BCH', '3h') }}</td>
                        <td>{{ abs('BCH', '12h') }}</td>
                        <td>{{ abs('BCH', '1d') }}</td>
                        <td>{{ abs('BCH', '3d') }}</td>
                        <td>{{ abs('BCH', '7d') }}</td>
                    </tr>
                </tbody>
            </table>
        </small>
    </div>
</template>

<script>
import moment from 'moment'
import Chart from 'chart.js'
import Vue from 'vue'

export default {
    name: 'pow-hashrateabs',
    data() {
        return {
            replacechart: null,
            hashchart: null
        }
    },
    computed: {
        replace_chain() {
            return moment.duration({
                minutes: this.data.BCH.pow.split.work / this.data.BTC.pow.current.work_actual * 10
            }).humanize()
        }
    },
    mounted() {
        this.redraw()
    },
    methods: {
        tsrate(coin, ts) {
            let o = this.data[coin].history['all'].find(h => h.timestamp === ts)
            return o ? o.rate : null
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                backgroundColor: this.btc_cbg,
                data: this.data.BTC.history['all'].map((obj) => {
                    return sprintf('%.2f', (obj.rate12h * Math.pow(2,32)) / 600 / Math.pow(10, 18))
                }),
                spanGaps: true,
                pointRadius: 0
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                backgroundColor: this.bch_cbg,
                data: this.data.BCH.history['all'].map((obj) => {
                    if (obj.rate12h < 1) return null
                    return sprintf('%.2f', (obj.rate12h * Math.pow(2,32)) / 600 / Math.pow(10, 18))
                }),
                spanGaps: true,
                pointRadius: 0
            })

            this.add_dataset({
                label: 'Total',
                borderColor: 'rgba(0,0,0,0.5)',
                borderDash: [3, 3],
                data: this.data.BCH.history['all'].map((obj) => {
                    if (obj.rate12h < 1) return null
                    let o = this.tsrate('BTC', obj.timestamp)
                    return o ? sprintf('%.2f', ((obj.rate12h + o) * Math.pow(2,32)) / 600 / Math.pow(10, 18)): null
                }),
                spanGaps: true,
                pointRadius: 0
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
                                labelString: 'Exahashes per second'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return this.sprintf('%.2f', value)
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-hashrateabs', config, labels)
        },
        abs(coin, avg = false) {
            if(this.data[coin].averages.last.blocks < 1) return null
            return sprintf('%.2f', this.data[coin].averages.last['rate' + avg] * Math.pow(2,32) / 600 / Math.pow(10, 18))
        },
        tsrate(coin, ts) {
            let o = this.data[coin].history['all'].find(h => h.timestamp === ts)
            return o ? o.rate12h : null
        },
    }
}
</script>

<style>

</style>
