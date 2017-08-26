
<template>
    <div>
        <ul class="nav nav-tabs">
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'security.fork'}">Fork rewrite</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'security.chain'}">Chain rewrite</router-link>
            </li>
        </ul>
        <div class="lol-window-chart">
        <transition name="fade">
            <router-view></router-view>
        </transition>
        <p>
            <br>The idea for these charts originate from Pieter Wuille. 
            See <a target="_blank" href="http://bitcoin.sipa.be/index.html">bitcoin.sipa.be</a> "Proof of Work Equivalent Days".
        </p>
        </div>

    </div>
</template>

<script>
export default {
    name: 'security'
}
</script>
<style>
</style>


