package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	maplock sync.RWMutex
	//消息广播channel
	Message chan string
}

//创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}
//监听Message广播消息channel的goroutine,一旦有消息就发送给全部的user
func(s *Server)LitenMessageer(){
	for{
		msg:=<-s.Message
		//将msg发送给全部的在线user
		s.maplock.Lock()
		for _,cli:=range s.OnlineMap{
			cli.C<-msg
		}
		s.maplock.Unlock()
	}
}
//广播消息方法
func (s *Server)BroadCast(user *User,msg string){
	sendMsg:="["+user.Addr+"]"+user.Name+":"+msg
	s.Message<-sendMsg
}
func (s *Server)Handler(conn net.Conn){
	// fmt.Printf("\"连接建立成功...\": %v\n", "连接建立成功...")
	//用户上线,将用户加入到onlinemap中
	user:=NewUser(conn,s)
	user.Online()
	

	//监听用户是否活跃的channel
	isLive:=make(chan bool)

	//接受客户端发送的消息
	go func(){
		buf:=make([]byte,4096)
		for{
			n, err := conn.Read(buf)
			if n == 0{
				user.Offline()
				return 
			}
			if err!=nil&&err!=io.EOF{
				fmt.Printf("\"err\": %v\n", "err")
				return 
			}
			msg:=string(buf[:n-1])
			//用户针对msg进行消息处理
			user.DoMessage(msg)
			//用户的任意消息,代表当前用户是活跃的
			isLive<-true
		}
	}()
	for{
		select{
		case<-isLive:
			//当前用户是活跃的，应该充值定时器
		case <-time.After(time.Second*10):
			//已经超时
			//将当前的user强制关闭
			user.SendMsg("你被踢了")
			//销毁用的资源
			close(user.C) 
			//关闭连接
			conn.Close()
			//退出当前Handler
			return 
		}
	}
}
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Printf("\"net.listen\": %v\n", "net.listen")
		return 
	}
	defer listener.Close()
	//启动监听Message的goroutine
	go s.LitenMessageer()
	//accept

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:",err)
			continue
		}
		go s.Handler(conn)
	}
	//do handler

	//close socket listen
}