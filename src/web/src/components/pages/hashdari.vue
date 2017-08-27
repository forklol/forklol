<template>
    <div>
        <p>
            Total hashrate during times when DARI is highest. Orange line means BTC DARI was highest, blue means BCH DARI was highest. 
        </p>

        7 day rolling average
        <hashrate-dari avg="rate7d"></hashrate-dari>
        <br>
        3 day rolling average
        <hashrate-dari avg="rate3d"></hashrate-dari>
        <br>
        1 day rolling average
        <hashrate-dari avg="rate1d"></hashrate-dari>
        <br>
        12 hour rolling average
        <hashrate-dari avg="rate12h"></hashrate-dari>
        <br>
        3 hour rolling average
        <hashrate-dari avg="rate3h"></hashrate-dari>
    </div>
</template>

<script>
import Chart from 'chart.js'
import Vue from 'vue'
import HashDari from '../charts/hashrate-dari.vue'

export default {
    name: 'hashdari',
    data() {
        return {}
    },
    components: {
        'hashrate-dari': HashDari
    }
}
</script>

<style>
</style>
