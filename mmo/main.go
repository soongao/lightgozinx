package main

import (
	"fmt"
	"gozinx/ziface"
	"gozinx/znet"
	"mmo/api"
	"mmo/core"
)

// 当客户端建立连接的时候的hook函数
func OnConnecionAdd(conn ziface.IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)
	//同步当前的PlayerID给客户端， 走MsgID:1 消息
	player.SyncPid()
	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()
	// 上线玩家添加到world中
	core.WorldMgrObj.AddPlayer(player)

	conn.SetProperty("pid", player.Pid)
	// 同步上线玩家位置
	player.SyncSurrounding()
	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}

// 当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	//获取当前连接的Pid属性
	pid, _ := conn.GetProperty("pid")

	//根据pid获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//触发玩家下线业务
	if pid != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player ", pid, " left =====")

}

func main() {
	//创建服务器句柄
	s := znet.NewServer()

	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnecionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// api
	s.AddRouter(2, &api.WorldChatApi{})
	s.AddRouter(3, &api.MoveApi{}) //移动
	//启动服务
	s.Serve()
}
