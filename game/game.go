package game

import (
	"servicegunbattle/awesomeProject/game/room"
)

type Game struct {
	RoomsMap *room.Rooms
}
// NewGame 创建一个room游戏服务
func NewGame() *Game {
	ret := &Game{}
	ret.RoomsMap = room.InitRooms()

	return ret
}
var GameMes *Game

func InitGameMes()  {
	GameMes = NewGame()
}