import numeral from 'numeral'
import { sprintf } from 'sprintf-js'

export default {
    data() {
        return {
            chart: [],
            datasets: []
        }
    },
    computed: {
        poll() {
            return this.$store.getters['poll']
        }
    },
    watch: {
        poll() {
            if (this.redraw) {
                this.datasets = []
                this.redraw()
            }
        }
    },
    methods: {
        period_time() {
            return this.data.BTC.history.all[1].timestamp - this.data.BTC.history.all[0].timestamp
        },
        destroy_charts() {
            if (this.chart.length !== 0) {
                for (c in this.charts) {
                    this.chart[c].destroy()
                }
                this.chart = []
            }
        },
        add_dataset(set) {
            let dataset = {
                backgroundColor: '#fff'
            }
            Object.assign(dataset, set)

            dataset.lineTension = 0
            dataset.cubicInterpolationMode = 'default'

            this.datasets.push(dataset)
        },
        draw(id, config, labels) {
            this.destroy_charts()
            
            let p = $('#' +id).parent()
            p.children().remove()
            p.append(`<canvas id="${id}" style="height: 512px;"></canvas>`)

            let cfg = {
                data: {
                    labels: labels,
                    datasets: this.datasets
                },
            }
            Object.assign(cfg, config)


            if (!cfg.options.tooltips) {
                cfg.options.tooltips = {
                    mode: 'index',
                    intersect: false,
                }
            }

            if (!cfg.options.hover) {
                cfg.options.hover = {
                    mode: 'x',
                    intersect: false
                }
            }

            this.$nextTick(() => {
                var ctx = document.getElementById(id).getContext("2d")
                this.chart.push(new Chart(ctx, cfg))
            })
        },
        find_ts(coin, ts) {
            return this.data[coin].history['all'].find(h => h.timestamp === ts)
        }
    },
}