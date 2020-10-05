<script>
  let baseURL = "http://localhost:8080";
  let promise = getGames();
  let currentGameUUID = "";
  let creeps = [];
  let players = [];

  async function newGame() {
    const res = await fetch(`${baseURL}/api/game/new`);
    const data = await res.json();
    if (res.ok && data.created) {
      return data;
    } else {
      throw new Error(text);
    }
  }

  function selectGame(gameID) {
    currentGameUUID = gameID;
  }

  let cells = [];

  function updateGameBoard(data) {
    let x = data.game.size_x;
    let y = data.game.size_y;

    let c = [];
    for (var i = 0; i < x; i++) {
      c.push(
        Array.apply(null, Array(y)).map(function () {
          return 0;
        })
      );
    }
    
    if (data.game.creeps) {
      creeps = data.game.creeps;
      data.game.creeps.forEach(function (cc) {
        c[cc.position.x][cc.position.y] = 1;
      });
    }else {
      creeps = [];
    }

    if (data.game.players) {
      players = data.game.players;
      data.game.players.forEach(function (cc) {
        c[cc.position.x][cc.position.y] = 2;
      });
    } else {
      players = [];
    }

    if (data.game.powerups) {
      data.game.powerups.forEach(function (cc) {
        c[cc.position.x][cc.position.y] = "P" + cc.type;
      });
    }
    cells = c;
  }

  async function getGameDetails() {
    const res = await fetch(`${baseURL}/api/game/get/${currentGameUUID}`);
    const data = await res.json();

    if (res.ok) {
      return data;
    } else {
      throw new Error(text);
    }
  }

  function redrawGameBoard() {
    getGameDetails().then(function (data) {
      updateGameBoard(data);
    });
  }

  function getCellClass(currentCellValue) {
    if (currentCellValue == 1) {
      return "creep";
    } else if (currentCellValue == 2) {
      return "human";
    } else if (typeof currentCellValue == "string") {
      if (currentCellValue[0] == "P") {
        return "powerup";
      }
    }
    return "open";
  }

  async function getGames() {
    const res = await fetch(`${baseURL}/api/game/`);
    const data = await res.json();

    if (res.ok) {
      return data;
    } else {
      throw new Error(text);
    }
  }

  var intervalPointer;
  $: if (currentGameUUID != "") {
    intervalPointer = setInterval(() => redrawGameBoard(), 1000);
  } else {
    if (intervalPointer) {
      clearInterval(intervalPointer);
    }
  }

  $: if (currentGameUUID != "") {
    redrawGameBoard();
  }
</script>

<style>
  main {
    text-align: center;
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
  }

  h1 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 4em;
    font-weight: 100;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }

  .row {
    margin: 1px;
    border: 1px black;
    display: flex;
  }

  .cell {
    width: 10px;
    height: 10px;
    padding: 10px;
    margin: 1px;
    border: 1px black;
  }

  .creep {
    background-color: red;
  }

  .human {
    background-color: greenyellow;
  }

  .open {
    background-color: pink;
  }

  .powerup {
    background-color: blueviolet;
  }

  .board {
    display: block;
  }
</style>

<main>
  <h1>Welcome to SCONWAR</h1>
  {#await promise}
    <p>...waiting</p>
  {:then number}
    {#if number.games}
      <p>Running Games</p>
      {#each number.games as gameUUID}
        <li>
          <button on:click={() => selectGame(gameUUID)}> {gameUUID} </button>
        </li>
      {/each}
      <p>Selected Game : {currentGameUUID}</p>
    {:else}
      <p>There are no games running</p>
      <button on:click={newGame}>Start New Game</button>
    {/if}
  {:catch error}
    <p style="color: red">
      Failed to get game list (check that the server is running)
    </p>
  {/await}

  <div style="display:flex">
    <div class="board">
      {#each cells as r}
        <div class="row">
          {#each r as c}
            <div class="cell {getCellClass(c)}">{c}</div>
          {/each}
        </div>
      {/each}
    </div>
    <div class="sidebar">
      [Creeps]
      {#each creeps as cre}
      <div class="row">
        {cre.health} - {cre.id}
      </div>
    {/each}
    [Players]
      {#each players as ply}
      <div class="row">
        {ply.health} - {ply.name} - {ply.id}
      </div>
    {/each}
    </div>
  </div>
  
</main>
