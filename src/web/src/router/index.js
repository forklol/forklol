import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

import Home from '../components/home.vue'
import Summary from '../components/summary.vue'

import Reward from '../components/reward.vue'
import RewardBlocks from '../components/charts/blockreward.vue'
import RewardFees from '../components/charts/fees.vue'
import RewardOpportunity from '../components/charts/opportunity.vue'
import RewardDARI from '../components/charts/dari.vue'
import RewardAlgo from '../components/charts/algo.vue'
import RewardFeepct from '../components/charts/feepct.vue'
import RewardInflation from '../components/charts/inflation.vue'

import Pow from '../components/pow.vue'
import PowHashrate from '../components/charts/hashrate.vue'
import PowWork from '../components/charts/work.vue'
import PowDifficulty from '../components/charts/difficulty.vue'
import PowSpeed from '../components/charts/speed.vue'
import PowRetarget from '../components/charts/retarget.vue'

import Blocks from '../components/blocks.vue'
import BlocksTime from '../components/charts/blocktime.vue'
import BlocksSize from '../components/charts/blocksize.vue'
import BlocksCDD from '../components/charts/cdd.vue'
import BlocksHeight from '../components/charts/height.vue'

import Tx from '../components/tx.vue'
import TxTxs from '../components/charts/txs.vue'
import TxFee from '../components/charts/txfee.vue'

import Security from '../components/security.vue'
import SecurityFork from '../components/charts/rewrite-fork.vue'
import SecurityChain from '../components/charts/rewrite-chain.vue'

import HashDari from '../components/pages/hashdari.vue'
import BCH from '../components/charts/coin/daa.vue'
import BTC from '../components/charts/coin/btc.vue'

import Contact from '../components/contact.vue'

const router = new Router({
  mode: 'history',
  linkActiveClass: 'active',
  linkExactActiveClass: 'active',
  routes: [{
    path: '/',
    name: 'home',
    component: Summary,
  }, {
    path: '/hashdari',
    name: 'hashdari',
    component: HashDari
  }, {
    path: '/coin/bch',
    name: 'coin.bch',
    component: BCH
  }, {
    path: '/coin/btc',
    name: 'coin.btc',
    component: BTC
  }, {
    path: '/reward',
    name: 'reward',
    component: Reward,
    redirect: {name: 'reward.dari.btc'},
    children: [{
      path: '/reward/dari',
      name: 'reward.dari',
      component: RewardDARI,
      redirect: {name: 'reward.dari.btc'},
      children: [{
        path: '/reward/dari/btc',
        name: 'reward.dari.btc',
        component: RewardAlgo,
        props: {coin: 'BTC'}
      }, {
        path: '/reward/dari/bch',
        name: 'reward.dari.bch',
        component: RewardAlgo,
        props: {coin: 'BCH'}
      }]
    }, {
      path: '/reward/blocks',
      name: 'reward.blocks',
      component: RewardBlocks
    }, {
      path: '/reward/fees',
      name: 'reward.fees',
      component: RewardFees
    }, {
      path: '/reward/opportunity',
      name: 'reward.opportunity',
      redirect: {name: 'home'}
    }, {
      path: '/reward/lossgain',
      name: 'reward.lossgain',
      redirect: {name: 'reward.opportunity'}
    }, {
      path: '/reward/inflation',
      name: 'reward.inflation',
      component: RewardInflation
    }, {
      path: '/reward/feepct',
      name: 'reward.feepct',
      component: RewardFeepct
    }]
  }, {
    path: '/pow',
    name: 'pow',
    component: Pow,
    redirect: {name: 'pow.hashrateabs'},
    children: [{
      path: '/pow/hashrate',
      name: 'pow.hashrate',
      component: PowHashrate
    }, {
      path: '/pow/work',
      name: 'pow.work',
      component: PowWork
    }, {
      path: '/pow/hashrateabs',
      name: 'pow.hashrateabs',
      redirect: {name: 'pow.hashrate'}
    }, {
      path: '/pow/difficulty',
      name: 'pow.difficulty',
      component: PowDifficulty
    }, {
      path: '/pow/speed',
      name: 'pow.speed',
      component: PowSpeed
    }, {
      path: '/pow/retarget',
      name: 'pow.retarget',
      component: PowRetarget
    }]
  }, {
    path: '/blocks',
    name: 'blocks',
    component: Blocks,
    redirect: {name: 'blocks.time'},
    children: [{
      path: '/blocks/time',
      name: 'blocks.time',
      component: BlocksTime
    }, {
      path: '/blocks/size',
      name: 'blocks.size',
      component: BlocksSize
    }, {
      path: '/blocks/height',
      name: 'blocks.height',
      component: BlocksHeight
    }, {
      path: '/blocks/cdd',
      name: 'blocks.cdd',
      component: BlocksCDD
    }]
  }, {
    path: '/tx',
    name: 'tx',
    component: Tx,
    redirect: {name: 'tx.txs'},
    children: [{
      path: '/tx/txs',
      name: 'tx.txs',
      component: TxTxs
    }, {
      path: '/tx/fee',
      name: 'tx.fee',
      component: TxFee
    }]
  }, {
    path: '/security',
    name: 'security',
    component: Security,
    redirect: {name: 'security.fork'},
    children: [{
      path: '/security/fork',
      name: 'security.fork',
      component: SecurityFork
    }, {
      path: '/security/chain',
      name: 'security.chain',
      component: SecurityChain
    }]
  }, {
    path: '/contact',
    name: 'contact',
    component: Contact
  }
  ]
})
export default router
