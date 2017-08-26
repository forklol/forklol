import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

// root state object.
// each Vuex instance is just a single state tree.
const state = {
    data: {},
    loaded: false,
    error: false,
    poll: 0,
    last_poll: 0,
    btc_h: 0,
    bch_h: 0,

    polling: false,
    token: null
}

// mutations are operations that actually mutates the state.
// each mutation handler gets the entire state tree as the
// first argument, followed by additional payload arguments.
// mutations must be synchronous and can be recorded by plugins
// for debugging purposes.
const mutations = {
    data(state, data) {
        state.loaded = true
        state.error = false
        state.data = data


        state.btc_h = data.BTC.averages.last.height
        state.bch_h = data.BCH.averages.last.height
        state.poll++
    },
    polling(state, p) {
        state.polling = p.isPolling
        state.token = p.token
    },
    last_poll(state, last) {
        state.last_poll = last
    }
}

// actions are functions that cause side effects and can involve
// asynchronous operations.
const actions = {
    fix_data: ({ commit, dispatch }, data) => {
        if (data.BTC.history.all.length > data.BCH.history.all.length) {
            data.BTC.history.all.pop()
        } else if (data.BTC.history.all.length < data.BCH.history.all.length) {
            data.BCH.history.all.pop()
        }

        commit('data', data)
    },
    reload_data: ({ commit, dispatch }) => {

        let t = +Date.now() / 1000
        let antidos = t - t % 60

        axios.get(process.env.API + '?t=' + antidos).then((resp) => {
            dispatch('fix_data', resp.data)
        }).catch((err) => {
            this.loaded = false
            this.error = err && err.message ? err.message : err
        })
    },
    poll_data: ({ commit, dispatch, state }) => {
        if (state.polling === true) return

        let after = 0

        let between = (+Date.now() / 1000) - state.last_poll
        if (between < 30) {
            after = (30 - between) * 1000
        }

        let CancelToken = axios.CancelToken
        let token = CancelToken.source()
        commit('polling', { isPolling: true, token: token })

        setTimeout(() => {

            state.last_poll = +Date.now() / 1000
            axios.get(`${process.env.API}poll/${state.poll}?t=${state.last_poll}`, {
                timeout: 900 * 1000,
                cancelToken: state.token.token
            }).then((resp) => {
                dispatch('fix_data', resp.data)
                commit('polling', { isPolling: false, token: null })
                commit('last_poll', 0)
                dispatch('poll_data')
            }).catch((err) => {
                if (!axios.isCancel(err)) {
                    commit('polling', { isPolling: false, token: null })
                    dispatch('poll_data')
                }
            })
        }, after)
    },
    poll_stop: ({ commit, dispatch, state }) => {
        if (state.polling && state.token) {
            state.token.cancel('window inactive')
            commit('last_poll', 0)
        }

        commit('polling', { isPolling: false, token: null })
    }
}

// getters are functions
const getters = {
    data: state => state.data,
    loaded: state => state.loaded,
    error: state => state.error,
    poll: state => state.poll,
    btc_h: state => state.btc_h,
    bch_h: state => state.bch_h,
}

// A Vuex instance is created by combining the state, mutations, actions,
// and getters.
export default new Vuex.Store({
    state,
    getters,
    actions,
    mutations
})