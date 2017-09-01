### fork.lol

More docs coming soon...

<!> Currently working on a rewrite of the backend using full nodes instead of an external API. </!>

See the `src/web` directory for the web frontend. Most of the calculations for the charts are being done there except for the hashrate averages, these are done (partly) in the `src/backend` project (see bitcoin/difficulty.go).

More info on the front and back-end in `/src/web` and `src/backend` README's.




#### TODO

##### High priority
- [ ] Peer/code review
- [ ] Describe the goal fork.lol is trying to achieve
- [ ] ...

##### General:

- [ ] Write more comprehensive docs to get contributors started
- [ ] ...

##### Code:
- [ ] Testing (unit/e2e) front and back-end
- [ ] Refactor front-end: get rid of WETness/dead code/etc
- [ ] Refactor back-end: simplify, remove old code
- [ ] Better (web) error handling (we have none atm)
- [ ] Better menu structure, evaluate "categories"
- [ ] Do not "double" average hashrates **
- [ ] ...

##### Infra:
- [ ] Use full nodes instead of api calls for block data
- [ ] Put front and back-end on same loadbalancer
- [ ] Use HTTPS (depends on point above)
- [ ] ...

##### Features:
- [ ] Allow user to select time range to display in chart
- [ ] Allow user to select rolling average time (3h, 12h, 1d, etc.)
- [ ] Add some basic info about individual blocks (? low prio)
- [ ] Allow users to create custom linkable charts (kinda like jsfiddle, very low prio/long term)
- [ ] Create a public api? (could be expensive, very low prio)
- [ ] ...

##### Charts:
- [ ] Switch to highcharts/plotly (chart.js lacks some features)
- [ ] Percentage of coins spent since fork (thanks to Felix Weis) (*)
- [ ] Add income rate (USD) per hour (total reward per hour)
- [ ] Add tx fee rate (USD) per hour
- [ ] Add reorg attack costs (USD) for 10, 50, 100 (?) blocks
- [ ] Add costs (USD) to fill up 100MB(?) of blockspace (perhaps multiply with node count?) *
- [ ] Add mempool size statistics (tx fee catergories) *
- [ ] Add transaction type statistics (p2pkh, p2sh, p2wph, etc.) *
- [ ] UTXO stats *
- [ ] Total block chain size *
- [ ] Median transaction fee sat/b and USD *
- [ ] Empty blocks(?)
- [ ] BCH MTP(?)
- [ ] ...

\* Full nodes are required for this feature.

\*\* An average is computed twice (and over a longer timeframe) than is advertised in some cases (mainly visible in short timespan charts that have to do with hashrates)

