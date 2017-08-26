### fork.lol backend

More docs coming, just the basics for now..


Before you start: If you're only want to work on the web stuff you don't need to use a local backend. You can just use `http://api.fork.lol` as the backend in `/web/config/dev.env.js`.

The project is written in go. Since this is my first go project expect some funky business in the code.. beware! I love to hear how things can be improved.

What you need before running this:

1. Mysql database with the forklol data (see /assets/forklol_db.sql.gz to get started)
2. A bitcoinaverage.com pubkey/secret. Currently I'm on a paid plan to be able to get historical prices, will have to find a solution to this (more on this later when I have time).


Getting it running:

1. Make sure you have go dep installed: https://github.com/golang/dep
2. Set your `$GOPATH` to the repository root `GOPATH=<path to repo>`
3. Run `dep ensure -update` to fetch the dependencies.
4. Run `go run main.go --help` to see what arguments you need to set.
5. Run `go run main.go -debug true ...`. Debug option must be true when not making use of a nginx frontend (CORS header issue).