package room

import (
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/lib/util"
	"sort"
)

type players struct {
	usermap *util.Map
}

func NewPlayers() players {
	ret := players{
		usermap: &util.Map{},
	}
	return ret
}

func (us players) ForEach(f func(key int, s *Player)) {
	//logs.Info("666---ForEach")
	ss := us.GetPlayers()
	keys := make([]int, 0, len(ss))
	for k, _ := range ss {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, key := range keys {
		val, ok := ss[key]
		if ok {
			f(key, val)
		}
	}
}

func (us players) GetSeatBySeatId(sid int) *Player {
	//logs.Info("666---GetSeatBySeatId")
	_seat := us.usermap.Get(sid).(*Player)
	return _seat
}
func (us players) GetPlayerByUID(uid int) *Player {
	//logs.Info("666---GetSeatByUID")

	ss := us.GetPlayers()
	for _, v := range ss {
		if v.Userid == uid {
			return v
		}
	}
	return nil
}
func (us players) SetPlayer(id int, s *Player) {
	//logs.Info("666---SetSeat")
	us.usermap.Set(id, s)
}
func (us players) Length() int {
	//logs.Info("666---Length")

	return us.usermap.Len()
}
func (us players) GetPlayers() map[int]*Player {
	//logs.Info("666---GetSeats")

	ret := make(map[int]*Player)
	us.usermap.RLockRange(func(k interface{}, s interface{}) {
		tid := k.(int)
		switch s.(type) {
		case *Player:
			ret[k.(int)] = s.(*Player)
		default:
			logs.Info(" ---------------> ReadRange GetSeats seat not found seatID:%v", tid)
		}
	})
	return ret
}
func (us players) RemovePlayer(sid int) {
	us.usermap.Del(sid)
	//logs.Info("666---RemoveSeat----%v", sid)
}

