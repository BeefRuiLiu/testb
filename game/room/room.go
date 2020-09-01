package room

import (
	"fmt"
	"runtime"
	"servicegunbattle/awesomeProject/lib/logs"
	"servicegunbattle/awesomeProject/lib/qo"
	"servicegunbattle/awesomeProject/lib/util"
	"servicegunbattle/awesomeProject/pb"
)

type Room struct {
	Roomid 					int 					//房间id 随机生成
	Peoplecount           	int     				//房间人数
	state            		int    					//房间状态 0人数未满 1,人数已满，2已经开始
	RandSeed            	int    					//随机数种子
	FrameID            		int    					//FrameID
	BattleID            	int    					//战斗ID
	AllUIDs            		[]int    				//所有玩家ID
	ReadyUIDs            	[]int    				//已准备的玩家ID
	ql   					*qo.Qo              	// 同步处理模块
	CliConnMap				*util.Map				//用户 Clien
	Playermap            	players            		//玩家位置
	//AllPlayerOTP           	[]*pb.PlayerOperation  //玩家位置
	Userinfos            	[]*pb.BattleUserInfo            	//wanjia xinxi
}

func NewRoom(roomid int) *Room {
	r := &Room{
		Roomid: roomid,
		state: 0,
		FrameID: 1,
		Peoplecount: 0,
		ReadyUIDs: make([]int,0),
		AllUIDs: make([]int,0),
		CliConnMap:        &util.Map{},
		Playermap:       NewPlayers(),
		//AllPlayerOTP:      make([]*pb.PlayerOperation,0),
		Userinfos:      make([]*pb.BattleUserInfo,0),
	}

	return r
}
func (r *Room)GetRoomState() int {
	return r.state
}
func (r *Room)SetRoomState(state int) {
	r.state = state
}

func (r *Room) Go(fn func()) {
	if r == nil {
		return
	}
	if r.ql == nil {
		r.ql = qo.New()
	}
	r.ql.Go(fn)
}
type Rooms struct {
	rooms   *util.Map
}

//初始化 Rooms
func InitRooms() *Rooms {
	ret := &Rooms{
		rooms:   &util.Map{},
	}
	return ret
}

// ReadRange 遍历房间
func (rs *Rooms) ReadRange(fn func(roomid int, room *Room)) {
	tbs := make([]*Room, 0, rs.rooms.Len())
	rs.rooms.RLockRange(func(k interface{}, t interface{}) {
		rid := k.(int)
		switch t.(type) {
		case *Room:
			tbs = append(tbs, t.(*Room))
		default:
			logs.Info(" ---------------> ReadRange GetRoom room not found _RoomID:%v", rid)
		}
	})
	for _, t := range tbs {
		if t != nil {
			fn(t.Roomid, t)
		}
	}
}

func (rs *Rooms) GetRoom(_roomid int) *Room {
	r := rs.rooms.Get(_roomid)
	switch r.(type) {
	case *Room:
		return r.(*Room)
	default:
		logs.Info("  ---------------> GetRoom not found _roomid:%v", _roomid)
	}
	return nil
}

//根据配置创建牌桌
func (rs *Rooms) CreateRoom(roomid int) *Room {

	_room := NewRoom(roomid)

	logs.Info("tableId: %v--------------->创建新的桌子  ", roomid)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logs.Info("roomid:%v----------->%v", roomid, fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
					err, n, trace[:n]))
				rs.rooms.Del(roomid)
			}
		}()
		//room.Serve()
	}()

	//添加到tables中
	rs.rooms.Set(roomid, _room)

	return _room

}

// 检测是否删除room
func (rs *Rooms) RemoveRoom(_roomId int) {

	_room := rs.GetRoom(_roomId)
	if _room == nil {
		logs.Info("RoomId:%v-------------------------房间已经被销毁!",_room)
		return
	}
	len1 := rs.rooms.Len()
	//_room.Destroy()

	rs.rooms.Del(_roomId)
	len2 := rs.rooms.Len()
	logs.Info("tableId:%v---------------RemoveRoom------------------------------>len1:%v,len2:%v", _roomId, len1, len2)

}