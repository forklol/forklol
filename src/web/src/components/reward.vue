<template>
    <div>
        <ul class="nav nav-tabs">
          <li class="nav-item">
            <router-link class="nav-link" :to="{name: 'reward.dari'}">DARI</router-link>
          </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'reward.blocks'}">Blocks</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'reward.fees'}">Fees</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'reward.feepct'}">Fee %</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'reward.inflation'}">Inflation</router-link>
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
    name: 'inflation'
}
</script>
<style>
</style>
