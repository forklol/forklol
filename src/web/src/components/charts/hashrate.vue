
<template>
    <div>
        <hra></hra>
        <br>
        <hr>
        <br>
        <p>
            Relative hashrate in percentage of total (stacked, <u>3h averages</u>). 
        </p>
        <div style="height:512px; width: 100%;">
            <canvas id="canvas-hashrate" style="height: 512px;"></canvas>
        </div>
        <br>
        <p>
            Disclaimer: Please note that using a <u>3 hour average</u> is not the most reliable way of measuring this data. This data should be interpreted as an <u>estimate</u>. The real number can differ by several percent.
        </p>
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
                        <td>{{ perct('BTC', '3h') }}</td>
                        <td>{{ perct('BTC', '12h') }}</td>
                        <td>{{ perct('BTC', '1d') }}</td>
                        <td>{{ perct('BTC', '3d') }}</td>
                        <td>{{ perct('BTC', '7d') }}</td>
                    </tr>
                    <tr>
                        <td>BCH</td>
                        <td>{{ perct('BCH', '3h') }}</td>
                        <td>{{ perct('BCH', '12h') }}</td>
                        <td>{{ perct('BCH', '1d') }}</td>
                        <td>{{ perct('BCH', '3d') }}</td>
                        <td>{{ perct('BCH', '7d') }}</td>
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
import hra from './hashrateabs.vue'

export default {
    name: 'pow-hashrate',
    components: {
        hra: hra
    },
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

        if(this.$route.hash === '#rel') {
            location.hash = ''
            this.$nextTick(() => {
                location.hash = '#rel'
            })
        }
    },
    methods: {
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_cbg,
                backgroundColor: this.bch_cbge,
                data: this.data.BCH.history['all'].map((obj) => {
                    let o = this.tsrate('BTC', obj.timestamp)
                    if(!o) o = 0
                    return o ? sprintf('%.2f', obj.rate3h / (obj.rate3h + o) * 100) : null
                }),
                fill: true,
                spanGaps: true,
                pointRadius: 0
            })

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_cbg,
                backgroundColor: this.btc_cbge,
                data: this.data.BTC.history['all'].map((obj) => {
                    let o = this.tsrate('BCH', obj.timestamp)
                    if(!o) o = 0
                    return o ? sprintf('%.2f', obj.rate3h / (obj.rate3h + o) * 100) : null
                }),
                fill: '-1',
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
                            stacked: true,
                            display: true,
                            position: 'right',
                            scaleLabel: {
                                display: true,
                                labelString: 'Hashrate % of total'
                            },
                            ticks: {
                                callback: (value, index, values) => {
                                    return value + ''
                                }
                            }
                        }]
                    }
                }
            };

            this.draw('canvas-hashrate', config, labels)
        },
        perct(coin, avg = false) {
            if(this.data[coin].averages.last.blocks < 1) return null
            return sprintf('%.2f%%', this.data[coin].averages.last['rate' + avg] / (this.data.BTC.averages.last['rate' + avg] + this.data.BCH.averages.last['rate' + avg]) * 100)
        },
        tsrate(coin, ts) {
            let o = this.data[coin].history['all'].find(h => h.timestamp === ts)
            return o ? o.rate3h : null
        },
    }
}
</script>

<style>

</style>
