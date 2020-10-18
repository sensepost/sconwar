# sconwar

a bring your own client programming game

## building

run `make swagger install`

## docker

build the image with `make docker-image`

run with `docker run --rm -it -p 8080:8080 -e API_TOKEN=foo sconwar:local`

## test information

Example commands use [httpie](https://httpie.org/)

### register a new player

The UUID returned here is technically your "secret" to join games.

```text
$ post localhost:8080/api/player/register name=bob
HTTP/1.1 200 OK
Content-Length: 62
Content-Type: application/json; charset=utf-8
Date: Sat, 03 Oct 2020 15:40:10 GMT

{
    "Created": true,
    "UUID": "a3b7dee8-fa38-43dc-b635-1935cf0a4d6c"
}
```

### register a new game

```text
$ get localhost:8080/api/game/new
HTTP/1.1 200 OK
Content-Length: 62
Content-Type: application/json; charset=utf-8
Date: Sat, 03 Oct 2020 15:32:09 GMT

{
    "Created": true,
    "UUID": "09a997c0-e94d-41b0-97a4-d2abcbac0292"
}
```

### join the registered game

```text
$ post localhost:8080/api/game/join game_id=09a997c0-e94d-41b0-97a4-d2abcbac0292 player_id=a3b7dee8-fa38-43dc-b635-1935cf0a4d6c
HTTP/1.1 200 OK
Content-Length: 16
Content-Type: application/json; charset=utf-8
Date: Sat, 03 Oct 2020 15:32:28 GMT

{
    "Success": true
}
```

### start a game

```text
$ put localhost:8080/api/game/start/09a997c0-e94d-41b0-97a4-d2abcbac0292
HTTP/1.1 200 OK
Content-Length: 16
Content-Type: application/json; charset=utf-8
Date: Sat, 03 Oct 2020 15:34:54 GMT

{
    "Success": true
}
```

### get player status

```text
$ post localhost:8080/api/player/status game_id=09a997c0-e94d-41b0-97a4-d2abcbac0292 player_id=a3b7dee8-fa38-43dc-b635-1935cf0a4d6c
HTTP/1.1 200 OK
Content-Length: 116
Content-Type: application/json; charset=utf-8
Date: Sat, 03 Oct 2020 15:33:09 GMT

{
    "Player": {
        "Health": 100,
        "ID": "a3b7dee8-fa38-43dc-b635-1935cf0a4d6c",
        "Name": "rusty-magnet",
        "Position": {
            "X": 0,
            "Y": 9
        }
    }
}
```

### check surroundings (creep/people in range)

```text
$ post localhost:8080/api/player/surroundings game_id=09a997c0-e94d-41b0-97a4-d2abcbac0292  player_id=a3b7dee8-fa38-43dc-b635-1935cf0a4d6c
HTTP/1.1 200 OK
Content-Length: 29
Content-Type: application/json; charset=utf-8
Date: Sat, 03 Oct 2020 15:33:32 GMT

{
    "Creep": null,
    "Players": null
}
```

### action a move command

```text
$ http -pbB post localhost:8080/api/action/move game_player_id:='{"game_id" : "09a997c0-e94d-41b0-97a4-d2abcbac0292", "player_id" : "a3b7dee8-fa38-43dc-b635-1935cf0a4d6c"}' x:=14 y:=3
{
    "game_player_id": {
        "game_id": "09a997c0-e94d-41b0-97a4-d2abcbac0292",
        "player_id": "a3b7dee8-fa38-43dc-b635-1935cf0a4d6c"
    },
    "x": 14,
    "y": 3
}

{
    "Success": true
}
```
