package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

//用户上线业务
func (u *User) Online() {
	//用户上线,将用户加入到onlinemap中
	u.server.maplock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.maplock.Unlock()
	u.server.BroadCast(u, "已经上线")
}

//用户下线业务
func (u *User) Offline() {
	//用户下线,将用户从onlineMap中删除
	u.server.maplock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.BroadCast(u, "下线")
}
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

//用户处理消息业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户都有哪些
		u.server.maplock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.maplock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式:rename|张三
		newName:=strings.Split(msg,"|")[1]
		//判断name是否存在
		_,ok := u.server.OnlineMap[newName]
		if ok{
			fmt.Printf("\"当前用户名被使用\": %v\n", "当前用户名被使用")
		}else{
			u.server.maplock.Lock()
			delete(u.server.OnlineMap,u.Name)
			u.server.OnlineMap[newName]=u
			u.server.maplock.Unlock()
			u.Name=newName
			u.SendMsg("您已经更新用户名:"+u.Name+"\n")
		}
	} else if len(msg)>4&&msg[:3]=="to|"{
		//消息格式 to|张三|消息内容
		//获取对方用户名
		remoteName:=strings.Split(msg,"|")[1]
		if remoteName == ""{
			u.SendMsg("消息格式不正确,请使用\"to|张三|你好啊\"格式.\n")
			return 
		}
		//根据用户名,得到对方User对象
		remoteUser,ok:=u.server.OnlineMap[remoteName]
		if!ok{
			u.SendMsg("该用户不存在\n")
			return 
		}
		//获取消息内容，通过对方的user对象将消息发送出去
		content:=strings.Split(msg,"|")[2]
		if content == ""{
			u.SendMsg("无消息内容,请重新发送\n")
			return 
		}
		remoteUser.SendMsg(u.Name+"对您说:"+content)
	}else {
		u.server.BroadCast(u, msg)
	}

}

//监听当前user channel的方法，一点有消息，就直接发送给客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
