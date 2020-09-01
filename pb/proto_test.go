package pb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"servicegunbattle/awesomeProject/lib/logs"

	"reflect"
	"testing"
)

func Test(t *testing.T) {

	logs.Info("pb.CSID_TCP_LOGIN:%v",(CSID_TCP_LOGIN.Number()))

	logs.Info(" 测试==================")
	//TODO  测试protobuf消息 序列化 和 反序列化
	logs.Info("test protobuf ...")
	ack := CommandCode_ACK
	_c := CommandHeader{
		CommandCode:  &ack,
		Flag: proto.Int32(2),
	}
	data, err := proto.Marshal(&_c)
	logs.Info("Marshal ...",&_c,"Marshal:",data)
	//proto解析[]byte{} 与protobuf结构体
	_c2 := &CommandHeader{
	}
	err = proto.Unmarshal(data,_c2)
	fmt.Println("test err ...",err)

	logs.Info("test Unmarshal(data) ...",_c2)

	var a int = 50
	v := reflect.ValueOf(a) // 返回Value类型对象，值为50
	tt := reflect.TypeOf(a)  // 返回Type类型对象，值为int
	logs.Info("ValueOf(a) %v",v, "TypeOf(a)%v",tt)
	logs.Info("ValueOf(a).Type() %v", v.Type(),"ValueOf(a).Kind() %v", tt.Kind())


	var b [5]int = [5]int{5, 6, 7, 8}
	logs.Info(reflect.TypeOf(b).String(),"|", reflect.TypeOf(b).Kind(),"|",reflect.TypeOf(b).Elem()) // [5]int array int

	var Pupil CommandHeader
	p := reflect.ValueOf(Pupil) // 使用ValueOf()获取到结构体的Value对象

	fmt.Println(p.Type()) // 输出:pb.CommandHeader
	fmt.Println(p.Kind()) // 输出:struct

	//int 类型转换成 若干长度的[]byte
	var inta int8 = 13
	pp := reflect.ValueOf(inta)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &inta)

	logs.Info("var inta %v",pp.Type(),"=>%v",bytesBuffer.Bytes())
	aaa,err := BytesToIntU(bytesBuffer.Bytes())
	if err != nil {
		logs.Info("bytesToIntU--err:%v",err.Error())
	}
	logs.Info("bytesToIntUERR:%v",aaa)

}

//字节数(大端)组转成int(无符号的)
func BytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0},b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 8:
		var tmp uint64
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0,fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}
