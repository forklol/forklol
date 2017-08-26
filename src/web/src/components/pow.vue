
<template>
    <div>
        <ul class="nav nav-tabs">
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'pow.hashrate'}">Hashrate</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'pow.work'}">Total work</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'pow.retarget'}">Retarget</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'pow.speed'}">Speed</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'pow.difficulty'}">Difficulty</router-link>
            </li>
        </ul>
        <div class="lol-window-chart">
        <transition name="fade">
            <router-view></router-view>
        </transition>
        </div>
    </div>
</template>

<script>
export default {
    name: 'pow'
}
</script>
<style>
</style>

