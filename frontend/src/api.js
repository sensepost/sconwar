async function getGameInfo(baseURL,currentGameUUID) {
  const res = await fetch(`${baseURL}/api/game/info/${currentGameUUID}`);
  const data = await res.json();

  if (res.ok) {
    let fow = Math.round(data.fow);
    return {fow: fow , data: data};
  } else {
    throw new Error(text);
  }
}

async function getScores(baseURL,currentGameUUID) {
  const res = await fetch(`${baseURL}/api/game/scores/${currentGameUUID}`);
  const data = await res.json();

  if (res.ok) {
    return data.scores;
  } else {
    throw new Error(text);
  }
}

async function getLeaderboard(baseURL) {
  const res = await fetch(`${baseURL}/api/meta/leaderboard`);
  const data = await res.json();
  if (res.ok) {
    return data.scores;
  } else {
    throw new Error(res.text());
  }
}

async function getPlayerSurroundings(baseURL,currentGameUUID,player_id) {
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

async function getPlayerStatus(baseURL,currentGameUUID,player_id) {
  const res = await fetch(`${baseURL}/api/player/status`, {
    method: "POST",
    cache: "no-cache",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ game_id: currentGameUUID, player_id: player_id }),
  });

  const data = await res.json();

  if (res.ok) {
    return data;
  } else {
    return res.ok;
  }
}



export let api = {
    getGameInfo : getGameInfo,
    getScores : getScores,
    getLeaderboard : getLeaderboard,
    getPlayerSurroundings : getPlayerSurroundings,
    getPlayerStatus: getPlayerStatus
};