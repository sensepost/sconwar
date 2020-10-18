# 👾 sconwar - a bring your own client programming game

<img align="right" src="./images/logo.png" height="220" alt="sconwar">

`sconwar` is a "bring your own client" programming game where the only interface with the game is via a RESTful api. Some `sconwar` features are:

- Primary interface is via a RESTful API
- Turn-based game logic
- Written in Go, clients can be in any language

## installation

The game server can be compiled with:

```bash
make clean swagger-install deps swagger install
```

This will download dependencies and compile the `sconwar` executable.

If you prefer docker, simply run `make docker` and then run the built image with:

```bash
docker run --rm -it -p 8080:8080 -e API_TOKEN=foo sconwar:local
```

## how to play

<img align="right" src="./images/api.png" height="250" alt="sconwar">

The most important resource you need to know about the is the API documentation. Once the server is running you can find the documentation by browsing to it. Unless you have a custom hosting setup, you can find this at <http://localhost:8080/>.

### getting started overview

To start a `sconwar` game, you need two things:

- A player ID, obtainable by registering to the server. This ID is a secret, and you should treat it that way.
- A game ID, obtainable by starting a new game.

Depending on the server setup, an administrator could either share the key configured to setup new users, or you could ask for a player token.

### game rules

`sconwar` itself is really simple. A game board that is typically 20 by 20 tiles big is populated with a number of creep, powerups and other players. Players are given 30 seconds to issue commands in their turn, after which the next player will be granted a chance to issue an action. A player may queue up to two commands before their turn which will be executed as soon as its their turn.

Issue enough attack, move & pickup commands to be the last person standing, and win the round!

### starting a new game

Games are identified by a UUID which should be used as the `gameid` whenever an API call requires that. Before playing a game of `sconwar`, you need to start and join a game.

To start a new game, call the `game/new` endpoint, recording the returned UUID. Next, join your player to that game with the `game/join` endpoint. Once all of the players in the game have joined, the `game/start/{uuid}` endpoint should be called to start the game.

This game & player id combination is used in all `action/*` endpoints to issue commands.

## license

`sconwar` is licensed under a [GNU General Public v3 License](https://www.gnu.org/licenses/gpl-3.0.en.html). Permissions beyond the scope of this license may be available at [http://sensepost.com/contact/](http://sensepost.com/contact/).

The sconwar logo is a derivative work of [Mini Mike's Metro Minis](https://github.com/mikelovesrobots/mmmm), and the license is available [here](https://github.com/mikelovesrobots/mmmm/blob/master/LICENSE).
