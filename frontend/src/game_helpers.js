function updateBoardStateForEntity(board, entity) {
    board[board.length - entity.y][entity.x] = entity;
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

  export let game_helpers = {
    calculateHealthText: calculateHealthText,
    updateBoardStateForEntity : updateBoardStateForEntity
  };