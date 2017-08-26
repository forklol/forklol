<template>
    <div>
        <ul class="nav nav-tabs">
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'blocks.time'}">Blocks/hr</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'blocks.size'}">Blocksize</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'blocks.height'}">Height</router-link>
            </li>
            <li class="nav-item">
                <router-link class="nav-link" :to="{name: 'blocks.cdd'}">Coindays</router-link>
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
    name: 'blocks'
}
</script>
<style>
</style>

