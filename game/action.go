package game

import (
	"fmt"
	"time"

	"github.com/sensepost/sconwar/storage"
)

// ActionType is an Action that can be invoked
type ActionType int

// Supported Actions
const (
	Move ActionType = iota
	Attack
	Pickup
	Nothing
)

// Action is the action that can be taken by a player
type Action struct {
	Action ActionType
	X      int
	Y      int
}

// NewAction starts a new action instance
func NewAction(action ActionType) *Action {
	return &Action{
		Action: action,
	}
}

// SetXY sets the x, t of an action
func (a *Action) SetXY(x int, y int) {
	a.X = x
	a.Y = y
}

// Execute executes an action
func (a *Action) Execute(player *Player, board *Board) {

	e := &storage.Event{
		Date:        time.Now(),
		SrcEntity:   int(PlayerEntity),
		SrcEntityID: player.ID,
		SrcPos:      player.Position.ToSingleValue(),
	}

	switch a.Action {
	case Move:

		e.DstEntity = int(PlayerEntity)
		e.DstEntityID = player.ID
		e.DstPos = player.Position.ToSingleValue()
		e.Action = int(Move)

		//This is a slight hack to let the player move in a random direction
		//If the requested X/Y co-ords are set to -1 (which is not possible as it would be out of the board)
		//The player.Move is executed which moves the player is moved in a random direction.
		if a.X == -1 && a.Y == -1 {
			player.Move()
			e.Msg = fmt.Sprintf(`player %s moved to a new position`, player.Name)
			break
		}

		distance := distanceBetween(player, &Position{X: a.X, Y: a.Y})

		if player.HasAvailableBuf(Teleport) {
			// for moves > MaxMoveDistance, use the teleport buf
			if distance > MaxPlayerMoveDistance {
				player.RemoveBuf(Teleport)
			}
		} else {
			if distance > MaxPlayerMoveDistance {
				e.Msg = fmt.Sprintf(`player %s tried to move to a position that is out of range`,
					player.Name)
				break
			}
		}

		player.MoveTo(a.X, a.Y)
		e.Msg = fmt.Sprintf(`player %s moved to a new position`, player.Name)

		break
	case Attack:
		// find entities on the x, y and if there is something, take damage
		for _, c := range board.AliveCreep() {

			if !player.IsInAttackRangeOf(c) {
				continue
			}

			cx, cy := c.GetPosition()
			if cx == a.X && cy == a.Y {

				multiplier := 1
				if player.HasAvailableBuf(DoubleDamage) {
					multiplier = 2
					player.RemoveBuf(DoubleDamage)
				}

				dmg, h := c.TakeDamage(-1, multiplier)
				player.AddDamageDealt(dmg)

				e.DstEntity = int(CreepEntity)
				e.DstEntityID = c.ID
				e.DstPos = c.Position.ToSingleValue()
				e.Action = int(Move)
				e.Msg = fmt.Sprintf(`player %s attacked creep %s for %d damage`,
					player.Name, c.Name, dmg)

				if h == 0 {
					board.LogEvent(&storage.Event{
						Date:        time.Now(),
						SrcEntity:   int(PlayerEntity),
						SrcEntityID: player.ID,
						SrcPos:      player.Position.ToSingleValue(),
						Action:      int(Attack),
						Msg:         fmt.Sprintf(`creep %s has been killed`, c.Name),
					})

					player.RecordCreepKilled()
				}

				break
			}
		}

		for _, p := range board.AlivePlayers() {

			if !player.IsInAttackRangeOf(p) {
				continue
			}

			px, py := p.GetPosition()
			if px == a.X && py == a.Y {

				multiplier := 1
				if player.HasAvailableBuf(DoubleDamage) {
					multiplier = 2
					player.RemoveBuf(DoubleDamage)
				}

				dmg, h := p.TakeDamage(-1, multiplier)
				player.AddDamageDealt(dmg)

				e.DstEntity = int(PlayerEntity)
				e.DstEntityID = p.ID
				e.DstPos = p.Position.ToSingleValue()
				e.Action = int(Move)
				e.Msg = fmt.Sprintf(`player %s attacked player %s for %d damage`,
					player.Name, p.Name, dmg)

				if h == 0 {
					board.LogEvent(&storage.Event{
						Date:        time.Now(),
						SrcEntity:   int(PlayerEntity),
						SrcEntityID: player.ID,
						SrcPos:      player.Position.ToSingleValue(),
						Action:      int(Attack),
						Msg:         fmt.Sprintf(`player %s has been killed`, p.Name),
					})

					player.RecordPlayerKilled()
					p.SaveFinalScore(board.ID, board.CurrentDeathPosition())
				}

				break
			}
		}

		break
	case Pickup:
		for _, u := range board.PowerUps {

			// use attack range for now. maybe change this to its own
			// unique value? maybe an argument to IsInRangeOf
			if !player.IsInAttackRangeOf(u) {
				continue
			}

			ux, uy := u.GetPosition()
			if ux == a.X && uy == a.Y {
				player.GivePowerUp(*u)
				player.AddScore(PickedUpPowerup)
				board.RemovePowerUp(u)

				e.DstEntity = int(PowerupEntity)
				e.DstEntityID = u.ID
				e.DstPos = u.Position.ToSingleValue()
				e.Action = int(Pickup)
				e.Msg = fmt.Sprintf(`player %s picked up a powerup of type %d`, player.Name, u.Type)

				break
			}
		}

		break
	default:
		board.LogEvent(&storage.Event{
			Date:        time.Now(),
			SrcEntity:   int(PlayerEntity),
			SrcEntityID: player.ID,
			SrcPos:      player.Position.ToSingleValue(),
			Action:      int(Attack),
			Msg:         fmt.Sprintf(`player %s managed to queue an invalid action`, player.Name),
		})
		break
	}

	if len(e.Msg) == 0 {
		e.Msg = `no action was executed`
	}

	board.LogEvent(e)
}
