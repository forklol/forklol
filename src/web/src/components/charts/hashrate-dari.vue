
<template>
    <div>
        <div style="height:512px; width: 100%;">
            <canvas :id="'canvas-hashrate-dari-' + avg" style="height: 512px;"></canvas>
        </div>
    </div>
</template>

<script>
import moment from 'moment'

export default {
    props: {
        avg: String
    },
    data() {
        return {}
    },
    computed: {

    },
    mounted() {
        this.redraw()
    },
    methods: {
        hashrate_dari(coin, obj) {
            let o = this.ts_coin(coin === 'BTC' ? 'BCH' : 'BTC', obj.timestamp)

            if(!o) return null

            if(o.blocks === 0 || obj.blocks === 0) {
                return null
            }

            if(o.dari > obj.dari) return null

            let rate = o[this.avg] + obj[this.avg]
            return (rate * Math.pow(2,32)) / 600 / Math.pow(10, 18)
        },
        redraw() {
            let labels = this.data.BTC.history['all'].map((obj) => moment.unix(obj.timestamp).format("MMM, DD HH:mm"))

            this.add_dataset({
                label: 'BCH',
                borderColor: this.bch_c,
                backgroundColor: this.bch_cbg,
                data: this.data.BCH.history['all'].map((obj) => this.hashrate_dari('BCH', obj)),
                pointRadius: 0
            })

            this.add_dataset({
                label: 'BTC',
                borderColor: this.btc_c,
                backgroundColor: this.btc_cbg,
                data: this.data.BTC.history['all'].map((obj) => this.hashrate_dari('BTC', obj)),
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

            this.draw('canvas-hashrate-dari-' + this.avg, config, labels)
        },
    }
}
</script>

<style>

</style>
