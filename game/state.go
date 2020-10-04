package game

// Games represent all of the games currently initialised
var Games map[string]*Board

// Setup initialises the Games variable
func Setup() {
	Games = make(map[string]*Board)
}
