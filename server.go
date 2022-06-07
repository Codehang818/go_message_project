package main

import (
	"fmt"
	"net"
	"sync"
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
	user:=NewUser(conn)
	s.maplock.Lock()
	s.OnlineMap[user.Name]=user
	s.maplock.Unlock()
	s.BroadCast(user,"已上线")
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