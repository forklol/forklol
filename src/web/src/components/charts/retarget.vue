
<template>
    <div>
        <p>
            Estimated days for difficulty retarget.
        </p>
        <div style="height:100%; width: 100%;">
            <canvas id="canvas-retarget" style="height: 512px;"></canvas>
        </div>
        <br>
        <small>
        <p>
            Estimated next difficulty. The 12 hour hashrate average is used to make a prediction about upcoming blocks in the period. Please note that because of the BCH EDA exploit these number can be inaccurate right after difficulty changes.
        </p>

        <table class="table table-sm table-bordered">
            <thead class="thead-inverse">
                <tr>
                    <th>Coin</th>
                    <th>Date (UTC 24h)</th>
                    <th>Remaining</th>
                    <th>Change</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>BTC</td>
                    <td>{{ date('BTC') }}</td>
                    <td>
                        {{ days('BTC') }}
                        <small>({{ left('BTC') }} blk)</small>    
                    </td>
                    <td>
                        <i v-if="next_pct('BTC', true) > 1" class="fa fa-chevron-up up"></i>
                        <i v-if="next_pct('BTC', true) < 1" class="fa fa-chevron-down down"></i>
                        {{ next('BTC') }}
                    </td>
                </tr>
                <tr>
                    <td>BCH</td>
                    <td>{{ date('BCH') }}</td>
                    <td>
                        {{ days('BCH') }}
                        <small>({{ left('BCH') }} blk)</small>    
                    </td>
                    <td>
                        <i v-if="next_pct('BCH', true) > 1" class="fa fa-chevron-up up"></i>
                        <i v-if="next_pct('BCH', true) < 1" class="fa fa-chevron-down down"></i>
                        {{ next('BCH') }}
                    </td>
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
import humanize from 'humanize-duration'

export default {
    name: 'pow-retarget',
    mounted() {
        this.redraw()
    },
    methods: {
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                data: this.data.BTC.history['all'].map((obj) => this.estimate('BTC', obj)),
                spanGaps: true,
            })

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                data: this.data.BCH.history['all'].map((obj) => this.estimate('BCH', obj)),
                spanGaps: true,
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
                                labelString: 'Days'
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

            this.draw('canvas-retarget', config, labels)
        },
        estimate(coin, obj, m = false) {
            let prev = this.data[coin].retargets.slice(-1)[0]
            let last = obj

            if (prev === undefined) prev = this.preforktarget

            let mod = last.height % 2016
            if(mod === 0 || mod === 1) return null


            let left = 2016 - mod
            let blkhr = obj.rate12h / obj.avg_diff * 6

            let hr = left / blkhr

            let days = hr / 24
            let mins = hr * 60


            if(days > 28) return null

            return m === false ? this.sprintf('%.2f', days) : mins
        },
        days(coin) {
            let mins = this.estimate(coin, this.data[coin].averages['last'], true)
            if(mins === null) return '>20 days'
            return humanize(mins * 60 * 1000, { largest: 2, round: true })
        },
        date(coin) {
            let mins = this.estimate(coin, this.data[coin].averages['last'], true)
            if(mins === null) return 'T.B.D.'

            let when = moment().add(mins, 'm')

            return when.format('MMM Do, HH:mm')
        },
        next_dari(coin) {

        },
        next(coin) {
            let change = this.next_pct(coin, true)
            if (change === null) return 'T.B.D.'

            let pct = 0

            if (change > 1) {
                pct = (change - 1) / 1 * 100
            } else {
                pct = (1 - change) / 1 * 100
            }

            return change > 1 ? this.sprintf('+%.2f%%', pct) : this.sprintf('-%.2f%%', pct)
        },
        left(coin) {
            let h = this.data[coin].averages.last.height

            return (2016 - h % 2016)
        }
    }
}
</script>

<style scoped>
.up {
    color: green;
}

.down {
    color: red;
}
</style>
