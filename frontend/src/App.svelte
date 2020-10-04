<script>
  export let name = "Joe";
  let baseURL = "http://localhost:8080";
  let promise = getGames();
  let currentGameUUID = "";
  let creeps = [];
  let x = 20;
  let y = 20;

  async function newGame() {
    const res = await fetch(`${baseURL}/api/game/new`);
    const data = await res.json();
    if (res.ok && data.Created) {
      return data;
    } else {
      throw new Error(text);
    }
  }

  function selectGame(gameID) {
    currentGameUUID = gameID;
  }

  $: if (currentGameUUID != "") {
    getGameDetails();
  }

  let cells = [];
  async function getGameDetails() {
    const res = await fetch(`${baseURL}/api/game/get/${currentGameUUID}`);
    const data = await res.json();

    if (res.ok) {
      //   creeps = data.Game.Creeps;
      x = data.Game.SizeX;
      y = data.Game.SizeY;

      let c = [];
      for (var i = 0; i < x; i++) {
        var a = Array.apply(null, Array(y)).map(function (
          currentValue,
          mapIndex
        ) {
          return 0;
        });
        c.push(a);
      }
      if (data.Game.Creeps) {
        data.Game.Creeps.forEach(function (cc) {
          c[cc.Position.X][cc.Position.Y] = 1;
        });
      }
      if (data.Game.Players) {
        data.Game.Players.forEach(function (cc) {
          c[cc.Position.X][cc.Position.Y] = 2;
        });
      }
      cells = c;

      return data;
    } else {
      throw new Error(text);
    }
  }

  function getCellClass(currentCellValue) {
    if (currentCellValue == 1) {
      return "creep";
    } else if (currentCellValue == 2) {
      return "human";
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
    intervalPointer = setInterval(() => getGameDetails(), 1000);
  } else {
    if (intervalPointer) {
      clearInterval(intervalPointer);
    }
  }
  //   let creepMap = {};
  //   //   $: creeps.forEach(function (c) {
  //   //     creepMap[c.Position.x + "," + c.Position.y] = 1;
  //   //   });
  //   let humans = [{ x: 11, y: 4 }];
  //   let humanMap = {};
  //   humans.forEach(function (c) {
  //     humanMap[c.x + "," + c.y] = 1;
  //   });
  //   console.log(creepMap);

  //   let final = '<div class="board">';

  //   for (var i = 0; i < x; i++) {
  //     final += '<div class="cell">';
  //     for (var k = 0; k < y; k++) {
  //       final += '<div class="cell1">';
  //       if (creepMap[i + "," + k]) {
  //         final += "C";
  //       } else if (humanMap[i + "," + k]) {
  //         final += "H";
  //       } else {
  //         final += "";
  //       }
  //       final += "</div>";
  //     }
  //     final += "</div>";
  //   }

  //   final += "</div>";
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

  .cell {
    margin: 1px;
    border: 1px black;
    display: flex;
  }

  .cell1 {
    width: 10px;
    height: 10px;
    background: slategray;
    padding: 10px;
    margin: 1px;
    border: 1px black;
  }

  .creep {
    color: red;
  }

  .human {
    color: greenyellow;
  }

  .human {
    color: pink;
  }

  .board {
    display: flex;
    display: block;
  }
</style>

<main>
  <h1>Hello {name}!</h1>
  {#await promise}
    <p>...waiting</p>
  {:then number}
    {#if number.Games}
      <p>Running Games</p>
      {#each number.Games as gameUUID}
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
    <p style="color: red">{error.message}</p>
  {/await}

  <div class="board">
    {#each cells as r}
      <div class="cell">
        {#each r as c}
          <div class="cell1 {getCellClass(c)}">{c}</div>
        {/each}
      </div>
    {/each}
  </div>
</main>
