package base

import (
	json "encoding/json"
	"fmt"
	io "io/ioutil"
)
var KeyMap = &KeyMapData{}

var Gamejson = "GameCfg"
func InitReadKeyMapData() {
	JsonParse := NewJsonStruct()
	JsonParse.Load("conf/gamecfg.json", KeyMap)
	fmt.Println(Gamejson,"--------------------------------------> ReadJson keymap success !!! ")
	PrintKeyMapData()
}
func PrintKeyMapData() {
	fmt.Println(Gamejson,"-------------------------------------->json数据")
	fmt.Println(Gamejson,"------------------->ProjectName:",KeyMap.ProjectName)
	fmt.Println(Gamejson,"------------------->TCP端口:",KeyMap.Tcpport)
	fmt.Println(Gamejson,"------------------->TCP地址:",KeyMap.TcpAddress)
	//fmt.Println(Gamejson,"------------------->redis地址:",KeyMap.RedisAddress)
	//fmt.Println(Gamejson,"------------------->redis密码:",KeyMap.RedisPassWord)
	//fmt.Println(Gamejson,"------------------->mysql数据库地址:",KeyMap.DataSourceName)
}
type KeyMapData struct {
	ProjectName string
	Tcpport     string
	TcpAddress     string
	RedisAddress     string
	RedisPassWord string						 //redis密码
	DataSourceName string                        //mysql数据库地址



}
type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}

}

func (self *JsonStruct) Load(filename string, v interface{}) {
	
	data, err := io.ReadFile(filename)
	if err != nil {
		fmt.Println(Gamejson,"-------------------------------------->io.ReadFile(filename)err",err.Error())
		return
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, v)
	
	if err != nil {
		return
	}
	
}