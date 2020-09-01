package bean

import (
	"net"
	"servicegunbattle/awesomeProject/tcpconn/baseim"
	"sync"
	"time"
)

type Clients struct {
	ClientsMap map[net.Conn]Client
	sync.RWMutex
}

func (cs *Clients) GetClient(key net.Conn) Client {
	cs.RLock()
	defer cs.RUnlock()
	return cs.ClientsMap[key]
}
func (cs *Clients) SetClient(key net.Conn, cli Client) {
	cs.Lock()
	defer cs.Unlock()
	cs.ClientsMap[key] = cli
}
func (cs *Clients) AlterClient(key net.Conn, cli Client) {
	cs.Lock()
	defer cs.Unlock()
	_, ok := cs.ClientsMap[key]
	if ok {
		cs.ClientsMap[key] = cli
	}

}
func (cs *Clients) DelClient(key net.Conn) {
	cs.Lock()
	defer cs.Unlock()
	delete(cs.ClientsMap, key)

}

func (cs *Clients) Range(f func(key net.Conn, value Client)) {
	//cs.Lock()
	for k, v := range cs.ClientsMap {
		f(k, v)
	}
	//cs.Unlock()
}
//每次连接的客户端
type Client struct {
	//ClientConn    		net.Conn          			//客户端连接
	Uid        int    //uid
	Name       string //name
	RemoteAddr string //
	C          chan []byte
	Heartbeat  int64 //上次心跳时间
}

func NewClient(remoteAddr string) *Client {
	cli := &Client{
		RemoteAddr: remoteAddr,
		C: make(chan []byte, 1000),
		Heartbeat: time.Now().Unix(),
	}
	return cli
}
func (c Client)WriteClient(res []byte) {
	c.C <- res
}
//这是更新时间函数
func (this Client) UpdateTime() {
	this.Heartbeat = time.Now().Unix()
}

//var ClientMap sync.Map
var ClientMaps Clients
var OnlineUIDs map[int]int

func Newinit() {
	ClientMaps.ClientsMap = make(map[net.Conn]Client)
	OnlineUIDs = make(map[int]int)
}

//全局消息
var Broadcast = make(chan *ReadMessage, 1000) // broadcast channel
//房间消息
var RoomMessage = make(chan *ReadMessage, 1000) // broadcast channel

//rpgame.read_message
type ReadMessage struct {
	Id           int    `xorm:"bigint" json:"-"`                 //
	Sockettype   int    `xorm:"int" json:"sockettype"`           //1是聊天 2是通知消息,3 是登陆，4是心跳
	Code         int    `xorm:"int" json:"code"`                 //1是房间聊天2.是私人聊天
	Currencyid   int    `xorm:"int" json:"currencyid,omitempty"` //根据code判断id是roomid还是friendid
	Uid          int    `xorm:"int" json:"uid,omitempty"`
	Accountsname string `xorm:"varchar(255)" json:"accountsname,omitempty"`
	Data         string `xorm:"varchar(255)" json:"data,omitempty"`
	Photo        string `xorm:"varchar(255)" json:"photo,omitempty"`
	Sendtime     string `xorm:"<- varchar(255)" json:"time,omitempty"`
	Issend       int    `xorm:"int" json:"-"`
	Token        string `xorm:"- varchar(255)" json:"appkey,omitempty"`
	OrderNumber  string `xorm:"- varchar(255)" json:"ordernumber,omitempty"`
	RemoteAddr   string `xorm:"- varchar(255)" json:"-"`
}

func NewHeartbeat() *ReadMessage {
	rm := new(ReadMessage)
	rm.Sockettype = 2
	rm.Code = baseim.Code_Heartbeat
	return rm
}
