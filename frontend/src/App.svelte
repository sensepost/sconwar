<script>
  const apiUrl = __myapp.env.API_URL
  const debugFlag = __myapp.env.DEBUG

  let baseURL = apiUrl;
  let debug = debugFlag;
  let promise = getGames();

  let currentGameUUID = "";
  let currentGame = {};
  let playerStatus;

  let fow = 0;
  let player;
  let error;

  let creeps = [];
  let players = [];
  let scores = [];
  let leaderboard = [];
  let player_id;
  let board_x = 20;
  let board_y = 20;

  let power = false;

  let select = false;
  let start = false;
  let a = false;
  let b = false;
  let up = false;
  let down = false;
  let left = false;
  let right = false;
  let extra = false;
  let events = [];

  function hitButton(button) {
    if (button === "select") {
      select = true;
    } else if (button === "start") {
      start = true;
    } else if (button === "a") {
      a = true;
    } else if (button === "b") {
      b = true;
    } else if (button === "up") {
      up = true;
    } else if (button === "down") {
      down = true;
    } else if (button === "left") {
      left = true;
    } else if (button === "right") {
      right = true;
    } else if (button === "extra") {
      extra = true;
    }
    if (select && start && a && b && up && down && left && right && extra) {
      new Audio("./mp3/pwned.mp3").play();
      select = false;
      start = false;
      a = false;
      b = false;
      up = false;
      down = false;
      left = false;
      right = false;
      extra = false;
    }
  }

  function tooglePower() {
    power = !power;
    currentGameUUID = "";
    if (!power) {
      location.reload();
    }
  }

  async function getGameInfo() {
    const res = await fetch(`${baseURL}/api/game/info/${currentGameUUID}`);
    const data = await res.json();

    if (res.ok) {
      fow = Math.round(data.fow);
      return data;
    } else {
      throw new Error(text);
    }
  }

  function selectGame(game) {
    currentGame = game;
    currentGameUUID = game.id;

    getPlayerStatus().then(function(data){
        if(!data){
          error = "Your player doesnt exist in that game!";
          currentGame = null;
          currentGameUUID = null;
        }else{
          error = null;
          gameover = false;
          gameoverDisplay = false;

          playerStatus = data;

          getGameInfo().then(function (data) {
            board_x = data.size_x;
            board_y = data.size_y;
          });
          getScores();
          new Audio("./mp3/soundtrack.mp3").play();
        }
    });
  }

  let cells = [];
  let creepPositions = new Map();
  let gameover = false;
  let gameoverDisplay = false;

  function updateBoardStateForEntity(board, entity) {
    let boardLength = board.length;
    // have to -1 on the X since we are diffing against the  board length
    board[(boardLength - entity.x) - 1][entity.y] = entity;
    return board;
  }

  function calculateHealthText(cc) {
    let ch = "";
    if (cc.health > 75) {
      ch = "hi";
    } else if (cc.health < 76 && cc.health > 40) {
      ch = "med";
    } else if (cc.health < 41 && cc.health > 0) {
      ch = "low";
    } else if (cc.health === 0) {
      ch = "dead";
    } else {
      ch = "";
    }
    return ch;
  }

  function updateGameBoard(data) {
    getGameEvents();

    // get all the player info
    getPlayerStatus().then(function(data){
        if(data){
          playerStatus = data;
        }
    });

    let x = board_x;
    let y = board_y;

    let c = [];
    let creepNewPos = new Map();

    for (var i = 0; i < x; i++) {
      c.push(
        Array.apply(null, Array(y)).map(function () {
          return 0;
        })
      );
    }

    if (data.creep) {
      creeps = data.creep;
      data.creep.forEach(function (cc) {
        let pieceObj = {};
        pieceObj.type = 1;
        pieceObj.health = cc.health;
        pieceObj.healthString = calculateHealthText(cc);
        pieceObj.id = cc.id;
        pieceObj.x = cc.position.x - 1;
        pieceObj.y = cc.position.y - 1;

        c = updateBoardStateForEntity(c, pieceObj);
        creepNewPos.set(cc.id, pieceObj);
      });
    } else {
      creeps = [];
    }

    let alldead = 0;
    let playerx;
    let playery;

    //This is the current player not the enemy players
    if (data.player) {
      //TODO : refactor this out
      let cc = data.player;

      let pieceObj = {};
      pieceObj.type = 2; //Might want to change to be more identifiable
      pieceObj.health = cc.health;
      pieceObj.healthString = calculateHealthText(cc);
      pieceObj.id = cc.name;
      pieceObj.x = cc.position.x - 1;
      pieceObj.y = cc.position.y - 1;

      playerx = pieceObj.x;
      playery = pieceObj.y;

      c = updateBoardStateForEntity(c, pieceObj);
      alldead++;
    }

    if (data.players) {
      players = data.players;
      data.players.forEach(function (cc) {
        let pieceObj = {};
        pieceObj.type = 2;
        pieceObj.health = cc.health;
        pieceObj.healthString = calculateHealthText(cc);
        pieceObj.id = cc.name;
        pieceObj.x = cc.position.x - 1;
        pieceObj.y = cc.position.y - 1;

        c = updateBoardStateForEntity(c, pieceObj);
        alldead++;
      });
    } else {
      players = [];
    }

    if (data.powerups) {
      data.powerups.forEach(function (cc) {
        let pieceObj = {};
        pieceObj.type = "P" + cc.type;
        pieceObj.x = cc.position.x - 1;
        pieceObj.y = cc.position.y - 1;

        c = updateBoardStateForEntity(c, pieceObj);
      });
    }


    if(data.players && alldead === data.players.length+1){
      gameover = true;
      gameoverDisplay = true;
      setTimeout(() => {
        gameoverDisplay = false;
      }, 2000);
    }

    // fow calculations
    if(!gameover){
      c.forEach(function(item, index) {
        let rowIndex = index;
        item.forEach(function(item, index){
          let rev = (board_x - playerx) - 1;

          if(rowIndex > rev+fow
            || rowIndex < rev-fow 
            || index > playery+fow || index < playery-fow){
              let pieceObj = {};
              pieceObj.type = "FOG";
              c[rowIndex][index] = pieceObj;
          }
        });
      });
    }

    // find out what creeps have moved where by diffing to previous state stored in cells
    // create map of creep id and position moved
    if(creepPositions.size > 0){
      for (const [key, oldPos] of creepPositions.entries()) {   
        let newPos = creepNewPos.get(oldPos.id);

        if(newPos){
          let xdiff = newPos.x - oldPos.x;
          let ydiff = newPos.y - oldPos.y;

          // for animation purposes only let it animate 1 position
          if(xdiff > 1){ xdiff = 1;}
          if(ydiff > 1){ ydiff = 1;}

          if(xdiff === 1 && ydiff === 0 && document.getElementById(oldPos.id)){
            document.getElementById(oldPos.id).classList.add("moveup");
          } else if (xdiff === -1 && ydiff === 0 && document.getElementById(oldPos.id)) {
            document.getElementById(oldPos.id).classList.add("movedown");
          } else if (ydiff === -1 && xdiff === 0 && document.getElementById(oldPos.id)) {
            document.getElementById(oldPos.id).classList.add("moveright");
          } else if (ydiff === 1 && xdiff === 0 && document.getElementById(oldPos.id)) {
            document.getElementById(oldPos.id).classList.add("moveleft");
          } else if (xdiff === 1 && ydiff === -1 && document.getElementById(oldPos.id)) {
            document.getElementById(oldPos.id).classList.add("moveupleft");
          } else if (xdiff === -1 && ydiff === -1 && document.getElementById(oldPos.id)) {
            document.getElementById(oldPos.id).classList.add("movedownleft");
          } else if (xdiff === 1 && ydiff === 1 && document.getElementById(oldPos.id)) {
            document.getElementById(oldPos.id).classList.add("moveupright");
          } else if (xdiff === -1 && ydiff === 1 && document.getElementById(oldPos.id)) {
            document.getElementById(oldPos.id).classList.add("movedownright");
          }
        }
      }
    }

    setTimeout(() => {
      cells = c;
      creepPositions = creepNewPos;
    }, 1000);
  }

  async function getScores() {
    const res = await fetch(`${baseURL}/api/game/scores/${currentGameUUID}`, {
      method: "GET",
      cache: "no-cache",
      headers: {
        "Content-Type": "application/json",
      },
    });

    const data = await res.json();

    if (res.ok) {
      scores = data.scores;
    } else {
      throw new Error(text);
    }
  }

  async function getLeaderboard() {
    const res = await fetch(`${baseURL}/api/meta/leaderboard`, {
      method: "GET",
      cache: "no-cache",
      headers: {
        "Content-Type": "application/json",
      },
    });
    const data = await res.json();
    if (res.ok) {
      leaderboard = data.scores;
    } else {
      throw new Error(text);
    }
  }

  async function getPlayer() {
    const res = await fetch(`${baseURL}/api/player/`, {
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ player_id: player_id })
    });

    const data = await res.json();
    if (res.ok) { 
      return data;
    } else {
      return res.ok;
    }
  }

  async function getPlayerStatus() {
    const res = await fetch(`${baseURL}/api/player/status`, {
      method: "POST",
      cache: "no-cache", 
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ game_id: currentGameUUID, player_id: player_id })
    });

    const data = await res.json();

    if (res.ok) {
      return data;
    } else {
      return res.ok;
    }
  }

  async function getPlayerSurroundings() {
    const res = await fetch(`${baseURL}/api/player/surroundings`, {
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ game_id: currentGameUUID, player_id: player_id }), // body data type must match "Content-Type" header
    });

    const data = await res.json();

    if (res.ok) {
      return data;
    } else {
      throw new Error(text);
    }
  }

  function redrawGameBoard() {
    if (!gameover) {
      getPlayerStatus().then(function (statusInfo) {
        getPlayerSurroundings().then(function (data) {
          data.player = statusInfo.player;
          updateGameBoard(data);
        });
      });
    }
  }

  function getCellImageClass(currentCellValue) {
    if (String(currentCellValue.type).startsWith("1")) {
      return "creep" + String(currentCellValue.healthString);
    } else if (String(currentCellValue.type).startsWith("2")) {
      return "human";
    } else if (typeof currentCellValue.type == "string") {
      if (currentCellValue.type[0] == "P") {
        return "powerup";
      } else if (currentCellValue.type == "FOG"){
        return "fog";
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

  async function getGameEvents() {
    const res = await fetch(`${baseURL}/api/game/events/${currentGameUUID}`);
    const data = await res.json();

    const dtFormat = new Intl.DateTimeFormat("en-GB", {
      timeStyle: "medium",
      timeZone: "UTC",
    });

    if(data.events){
      data.events.forEach((e) => {
        e.CreatedAt = dtFormat.format(Date.parse(e.CreatedAt));
      });
    }

    if (res.ok) {
      //TODO Fix this , find out if a lib is missing
      events = data.events.reverse();
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

  const truncate = (text, length) => {
    length = length - 3; // make space for the ...
    if (text.length > length) {
      return `${text.substring(0, length)}...`;
    }
    return text;
  }

  let scoresActive = false;
  let eventsActive = false;
  let gamestarted = false;
  let leaderboardActive = false;

  const toggleScoreboard = () => {
    scoresActive = !scoresActive;
    if (scoresActive) {
      eventsActive = false;
      leaderboardActive = false;
    }
  };
  const toggleEvents = () => {
    eventsActive = !eventsActive;
    if (eventsActive) {
      scoresActive = false;
      leaderboardActive = false;
    }
  };
  const toggleLeaderboard = () => {
    leaderboardActive = !leaderboardActive;
    if (leaderboardActive) {
      // refresh leaderboard with game scores
      getLeaderboard();
      scoresActive = false;
      eventsActive = false;
    }
  };
  const toggleStartGame = () => {
    getPlayer().then(function (data) {
      if(!data){
        player_id = "INVALID - TRY AGAIN"
      }else{
        player = data;
        gamestarted = !gamestarted;
      }
    });
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
    margin: 0px;
    border: 1px black;
    display: flex;
  }

  .cell {
    width: 10px;
    height: 10px;
    padding: 10px;
    margin: 0px;
    border: 1px solid darkgrey;
    box-sizing: unset;
  }
  .headercell {
    border: unset;
  }

  .board {
    display: block;
  }
</style>

<link href="./css/hacker.css" rel="stylesheet" />
<main>
  <div style="position: absolute; left: 850px; top: 10px;">
   
    {#if debug == true }
    <div class="sidebar">
      [Creeps]
      {#each creeps as cre}
        <div class="row">{cre.health} - {cre.id}</div>
      {/each}
      [Players]
      {#each players as ply}
        <div class="row">{ply.health} - {ply.name}</div>
      {/each}
    </div>
    <div class="sidebar">
      [Events]
      {#if events}
      {#each events as eve}
        <div class="row">{eve.CreatedAt} - {eve.msg}</div>
      {/each}
      {/if}
    </div>
    {/if}
  </div>

  <div class="gameboy">
    <div
      class={power ? 'power-button on' : 'power-button'}
      on:click={() => tooglePower()} />
    <div class="gb-body">
      <div class="screen-box">
        <div class="power-box">
          <div class={power ? 'power-light on' : 'power-light'} />
          <div class="power-status">))</div>
          <div class="power-text">POWER</div>
        </div>
        <div class="screen-display">
          {#if currentGameUUID && power && !gameoverDisplay && gamestarted && !scoresActive && !eventsActive && !leaderboardActive}
            <div style="display:flex; box-sizing: unset;">
              <div class="board">
                <div class="row">
                  <div class="cell headercell central info" style="margin-left:30px">
                    1
                  </div>
                  <div class="cell headercell central info">2</div>
                  <div class="cell headercell central info">3</div>
                  <div class="cell headercell central info">4</div>
                  <div class="cell headercell central info">5</div>
                  <div class="cell headercell central info">6</div>
                  <div class="cell headercell central info">7</div>
                  <div class="cell headercell central info">8</div>
                  <div class="cell headercell central info">9</div>
                  <div class="cell headercell central info">10</div>
                  <div class="cell headercell central info">11</div>
                  <div class="cell headercell central info">12</div>
                  <div class="cell headercell central info">13</div>
                  <div class="cell headercell central info">14</div>
                  <div class="cell headercell central info">15</div>
                  <div class="cell headercell central info">16</div>
                  <div class="cell headercell central info">17</div>
                  <div class="cell headercell central info">18</div>
                  <div class="cell headercell central info">19</div>
                  <div class="cell headercell central info">20</div>
                </div>

                {#each cells as r, i}
                  <div class="row">
                    <div class="cell headercell central info">{ (board_x - i)}</div>
                    {#each r as c}
                      <div
                        id={c.id}
                        class="cell {getCellImageClass(c)}" />
                    {/each}
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          {#if power && !currentGameUUID && !gameoverDisplay && !gamestarted && !scoresActive && !eventsActive && !leaderboardActive}
            <div id="movetxt">
              <h3>Welcome to SCONWAR</h3>
              <br/>
              <img style="height:150px"alt="Space invader header" src="./images/invader.gif" />
              <br/>
              <div>Enter Player ID</div>
              <input id="player" type="text"  bind:value={player_id}/>
              {#if player_id}
              <div on:click={() => toggleStartGame()} >Click to start</div>
              {/if}
            </div>
          {/if}

          {#if power && !currentGameUUID && !gameoverDisplay && gamestarted && !scoresActive && !eventsActive && !leaderboardActive}
            <div class="info">
              <div class="infotext">WELCOME {player.Name}</div>
              {#if error}
                <p style="color: red">
                  {error}
                </p>
              {/if}
              <div class="infotextmed">GAMES IN PLAY</div>
              <div class="infoscrollhalf">
                {#await promise}
                  <p>Loading games</p>
                {:then number}
                  {#if number.games}
                    {#each number.games as gameUUID}
                      {#if gameUUID.status !== 2}
                        <div class="clickable" on:click={() => selectGame(gameUUID)}> {gameUUID.name} </div>
                      {/if}
                    {/each}
                  {:else}
                    <p>No games currently running</p>
                  {/if}
                {:catch error}
                  <p style="color: red">
                    Failed to get current game list (check that the server is running)
                  </p>
                {/await}
              </div>
              <div class="infotextmed">GAMES FINISHED</div>
              <div class="infoscrollhalf">
                {#await promise}
                  <p>Loading games</p>
                {:then number}
                  {#if number.games}
                    {#each number.games as gameUUID}
                      {#if gameUUID.status === 2}
                        <div
                          class="clickable"
                          on:click={() => selectGame(gameUUID)}>
                          {gameUUID.name}
                        </div>
                      {/if}
                    {/each}
                  {:else}
                    <p>No games currently running</p>
                  {/if}
                {:catch error}
                  <p style="color: red">
                    Failed to get game finished list (check that the server is running)
                  </p>
                {/await}
              </div>
            </div>
          {/if}

          {#if power && scoresActive && !eventsActive && !leaderboardActive}
            <div class="info">
              <div class="infotext">SCORES</div>
              <div class="infoscroll" style="padding-left:10px;">
                <table style="width:100%">
                  <tr>
                    <th class="central">Name</th>
                    <th class="central">D+</th>
                    <th class="central">D-</th>
                    <th class="central">KC</th>
                    <th class="central">KP</th>
                    <th class="central">Score</th>
                  </tr>

                  {#if scores}
                    {#each scores as score}
                      <tr>
                        <td>{truncate(score.name, 10)}</td>
                        <td>{score.damage_dealt}</td>
                        <td>{score.damage_taken}</td>
                        <td>{score.killed_creep}</td>
                        <td>{score.killed_players}</td>
                        <td>{score.score}</td>
                      </tr>
                    {/each}
                  {/if}
                </table>
              </div>
            </div>
          {/if}

          {#if power && leaderboardActive && !scoresActive && !eventsActive}
            <div class="info">
              <div class="infotext">LEADERBOARD</div>
              <div class="infoscroll" style="padding-left:10px;">
                <table style="width:100%">
                  <tr>
                    <th class="central">Player</th>
                    <th class="central">Game</th>
                    <th class="central">Pos</th>
                    <th class="central">D+/D-</th>
                    <th class="central">CK/PK</th>
                    <th class="central">Score</th>
                  </tr>

                  {#if leaderboard}
                    {#each leaderboard as score}
                      <tr>
                        <td>{truncate(score.name, 10)}</td>
                        <td>{truncate(score.game_name, 15)}</td>
                        <td>{score.position}</td>
                        <td>{score.damage_dealt}/{score.damage_taken}</td>
                        <td>{score.creep_kills}/{score.player_kills}</td>
                        <td>{score.score}</td>
                      </tr>
                    {/each}
                  {/if}
                </table>
              </div>
            </div>
          {/if}

          {#if power && currentGameUUID && gamestarted && !scoresActive && eventsActive && !leaderboardActive}
            <div class="info">
              <div class="infotext">EVENTS</div>
              <div class="infotextsmall">Game: {currentGame.name}</div>
              <div class="infotextsmall">Game ID: {currentGameUUID}</div>
              <div class="infoscroll">
                {#each events as eve}
                  <div
                    class="row"
                    style="position:relative; top:10px; left:0px;">
                    {eve.CreatedAt}:
                    {eve.msg}
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          {#if power && currentGameUUID && gameoverDisplay && gamestarted && !scoresActive && !eventsActive && !leaderboardActive}
            <div class="background">
              <div class="gameovertext">GAME OVER</div>
            </div>
          {/if}
        </div>

        <div class="playerstatus" > 
          {#if playerStatus}
            <span>
              Health: {playerStatus.player.health} 
              <br/>
              Power Ups: 
              {#if playerStatus.player.powerups}
              {#each playerStatus.player.powerups as pu}
                {#if pu.type === 0}
                  Health  
                {:else if pu.type === 1}
                  Teleport
                {:else}
                  Double Damage
                {/if}
                &nbsp;
              {/each}
              {/if}
          </span>
          {/if}
        </div>
        <div class="gameboy-color-logo"> 
          <div class="logo-gb">SCONWAR</div>
        </div>
      </div>
      <div class="button-box">
        <div class="arrow-group">
          <div class="up-box" on:click={() => hitButton('up')}>
            <span class="arrow up" />
          </div>
          <div class="right-box" on:click={() => hitButton('right')}>
            <span class="arrow right" />
          </div>
          <div class="down-box" on:click={() => hitButton('down')}>
            <span class="arrow down" />
          </div>
          <div class="center-box">
            <span class="dent"><span class="dent-highlight" /></span>
          </div>
          <div class="left-box" on:click={() => hitButton('left')}>
            <span class="arrow left" />
          </div>
        </div>
        <div class="ab-button a" on:click={() => hitButton('a')}>
          <span class="button-text-height">B</span>
        </div>
        <div class="ab-button b" on:click={() => hitButton('b')}>
          <span class="button-text-height">A</span>
        </div>
      </div>
      <div
        class="pill-button button-select"
        on:click={() => hitButton('select')}
        on:click={() => toggleScoreboard()}>
        <span class="select">Scores</span>
      </div>
      <div
        class="pill-button button-start"
        on:click={() => hitButton('start')}
        on:click={() => toggleEvents()}>
        <span class="start">Events</span>
      </div>
      <div
        class="pill-button button-extra"
        on:click={() => hitButton('extra')}
        on:click={() => toggleLeaderboard()}>
        <span class="extra">Leaderboard</span>
      </div>
      <div class="speaker">
        <div class="row1">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row2">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row3">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row4">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row5">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row6">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row7">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row8">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row9">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
        <div class="row10">
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
          <div class="dot-hole" />
        </div>
      </div>
    </div>
  </div>
</main>

