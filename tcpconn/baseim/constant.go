package baseim

//Http返回 Code	说明
const (
	Code_ExitConn      = -1 //-1.退出连接
	Code_RoomChat      = 1  //Code=1.房间聊天
	Code_FriendChat    = 2  //Code=2.私聊
	Code_Login         = 3  //Code=3.登陆
	Code_FriendRes     = 4  //Code=4,有新增好友请求
	Code_FriendSuc     = 5  //Code=5,有新增好友
	Code_TaskSuc       = 6  //Code=6,有新增任务达成
	Code_GiftAdd       = 7  //Code=7,有新增礼物盒子
	Code_NoticeAdd     = 8  //Code=8,新增公告
	Code_Heartbeat     = 9  //Code=9,心跳包
	Code_DrawCards     = 10 //Code=10,抽卡消息
	Code_ChangeRoom    = 11 //Code=11,更换房间id
	Code_CloseConn     = 20 //Code=20,通知客户端下线
	Code_NoFriendQuest = 21 //Code=21,路人通知助战消息
	Code_FriendQuest   = 22 //Code=22好友通知助战消息

	Code_PaySucceed = 23 //Code=23,支付成功 sockettype=5
	Code_PayFaild   = 24 //Code=24支付失败 sockettype=5

	Code_LoginTwo   = 25 //Code=25  异地登陆
)
