package main

import (
	"net"
	"fmt"
	"flag"
	"io"
	"os"
)

type Client struct{
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int
}

var serverIp string
var serverPort int

func init(){
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置需要连接的服务器地址ip(默认值为127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置需要连接的服务器的端口(默认值为8888)")
}

func (this *Client)Run(){
	for this.flag != 0 {
		for this.menu() != true{

		}
		switch this.flag {
		case 1:
			this.PublicChat()
			break
		case 2:
			this.PrivateChat()
			break
		case 3:
			this.updateName()
			break
		}
	}
}

func (this *Client)updateName() bool{
	fmt.Println(">>>>>>>请输入用户名：")
	fmt.Scanln(&this.Name)

	sendMsg := "rename|" + this.Name + "\n"

	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Conn write err :", err)
		return false
	}
	return true

}

func (this *Client) DealResponse(){
	//阻塞方法，conn一有消息就会输出到标准输出
	io.Copy(os.Stdout, this.conn)
}

func (client *Client) menu() bool{
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	
	fmt.Scanln(&flag)

	if flag >=0 && flag <=3{
		client.flag = flag
		return true
	}else{
		fmt.Println("输入的菜单有误")
		return false
	}
}

func main(){
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)


	if client == nil {
		fmt.Println("创建客户端连接失败")
		return
	}else{
		fmt.Println("创建客户端连接成功")
	}

	

	//单独使用一个goroutine去处理server回执的消息
	go client.DealResponse()

	client.Run()
}


func NewClient(serverIp string, serverPort int) *Client{
	client := &Client{
		ServerIp : serverIp,
		ServerPort : serverPort,
		flag : 999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	
	if err != nil {
		fmt.Println("连接出错: ", err)
		return nil
	}else{
		client.conn = conn
		return client
	} 
}

func (this *Client) SelectUser (){
	fmt.Println("当前系统用户为")
	queryUser := "who\n"
	_, err := this.conn.Write([]byte(queryUser))
	if err != nil {
		fmt.Println("查询在线用户失败：", err)
	}
}

func (this *Client) PrivateChat(){
	var remoteName string
	var sendMsg string

	this.SelectUser()

	fmt.Println(">>>>>>>>>>>>请输入需要对话的用户名，输入'exit'退出：")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>>>>>请输入聊天内容，输入'exit'退出")
		fmt.Scanln(&sendMsg)
		for sendMsg != "exit" {
			sendMsg  = "to|" + remoteName + "|" + sendMsg
			_, err := this.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("连接出错")
				break
			}
			sendMsg = ""
			fmt.Println(">>>>>>>>>>请输入聊天内容，输入'exit'退出")
			fmt.Scanln(&sendMsg)
		}
	}


	

}

func (this *Client) PublicChat(){
	var msg string

	fmt.Println(">>>>>>>>>>请输入聊天内容，输入'exit'退出")

	fmt.Scanln(&msg)

	for msg != "exit" {
		if len("") != 0 {
			sendMsg := msg + "\n"
			_, err := this.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("连接出错")
				break
			}
		}

		msg = ""
		fmt.Println(">>>>>>>>>>请输入聊天内容，输入'exit'退出")
		fmt.Scanln(&msg)
	}
}