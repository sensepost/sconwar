# 👾 sconwar - a bring your own client programming game

<img align="right" src="./images/logo.png" height="220" alt="objection">

`sconwar` is a "bring your own client" programming game.

- Primary interface is via a RESTful API
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

TODO

## license

`sconwar` is licensed under a [GNU General Public v3 License](https://www.gnu.org/licenses/gpl-3.0.en.html). Permissions beyond the scope of this license may be available at [http://sensepost.com/contact/](http://sensepost.com/contact/).

The sconwar logo is a derivative work of [Mini Mike's Metro Minis](https://github.com/mikelovesrobots/mmmm), and the license is available [here](https://github.com/mikelovesrobots/mmmm/blob/master/LICENSE).
