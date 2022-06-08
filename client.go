package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/qiniu/go-sdk/v7/sms/client"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverip string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverip,
		ServerPort: serverPort,
		flag: 999,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverip, serverPort))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	client.conn = conn
	//返回对象
	return client
}

func (c *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		return true
	} else {
		fmt.Println(">>>请输入合法范围内的数字<<<<")
		return false
	}
}
func(c *Client)Run(){
	for c.flag!=0{
		for c.menu()!=true{
			//根据不同的模式处理不同的业务
			switch  c.flag{
			case 1:
				fmt.Println("公聊模式选择...")
				break
			case 2:
				fmt.Println("私聊模式选择...")
				break
			case 3:
				fmt.Println("更新用户名选择...")
				break
			}
		}
	}
}
var serverip string
var serverPort int

func init() {
	flag.StringVar(&serverip, "ip", "127.0.0.1", "设置服务器ip地址(默认127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认8888)")
}
func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverip, serverPort)
	if client == nil {
		fmt.Printf("\">>>>>链接服务器失败...\": %v\n", ">>>>>链接服务器失败...")
		return
	}
	fmt.Printf("\">>>>>>链接服务器成功\": %v\n", ">>>>>>链接服务器成功")
	//启动客户端的业务
	client.
}
