<template>
    <div>
        <small>
            <table class="table table-sm table-bordered">
                <thead class="thead-inverse">
                    <tr>
                        <th>Coin</th>
                        <th v-for="h in stats.header" v-bind:key="h">{{ h }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(g, idx) in stats.coins" v-bind:key="idx">
                        <td>{{ g.coin }}</td>
                        <td v-for="(v, idxx) in g.values" v-bind:key="idx + '_' + idxx">{{ stats.f(g.coin, v) }}</td>
                    </tr>
                </tbody>
            </table>
        </small>
    </div>
</template>

<script>
export default {
    name: 'block-txfee-sat',
    props: {
        stats: Object
    },
    data() {
        return {
        }
    }
}
</script>
<style>

</style>