package main

import "github.com/cgit6/go_slot/game"

func runner() {
	game.NewConfig(game.REELSTRIPS, game.SYMBOLS, game.LINES, game.PAYTABLE, game.ROWS, game.COLS, game.ModeLine)
}

func main() {
	runner() // 執行模擬
}
