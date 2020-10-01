package game

var Games map[string]*Board

func Setup() {
	Games = make(map[string]*Board)
}
