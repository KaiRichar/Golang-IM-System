package main

import (
	"net"
	"fmt"
	"io"
	"sync"
	 "time"
)


type Server struct{
	Ip string
	Port int
	OnlineUserMap map[string]*User
	MessageChan chan string
	mapLock sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineUserMap: make(map[string]*User),
		MessageChan: make(chan string),
	}
	return server
}

//收到消息就通知所有客户端的goroutine
func (this *Server) ListenMessage(){
	// for msg := range this.MessageChan {
	for {
		msg := <-this.MessageChan	
		this.mapLock.Lock()
		for _, user := range this.OnlineUserMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}

}

//广播消息
func (this *Server) BroadCast(user *User, msg string){
	sendMsg := "用户-" + user.Name +"[" + user.Addr+ "]: "+ msg
	this.MessageChan <- sendMsg
}

func (this *Server) Hanlder(conn net.Conn){
	user := NewUser(conn, this)
	user.OnLine()
	isLive := make(chan bool)
	//接收用户发送的消息
	go func(){
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.OffLine();
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户发送的消息并去除换行
			msg := string(buf[:n-1])
			fmt.Println(">>>>>>>>收到消息" + msg)
			user.DoMessage(msg)

			isLive <- true

		}
	}()
	for {
		select {
			case <-isLive:
			case <-time.After(time.Second * 300):
				user.SendMsg("您已被踢出")
				//关闭用户管道
				close(user.C)
				//关闭连接
				conn.Close()
				return
		}
	}
	
}

func (this *Server) Start(){
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port));
	if err != nil {
		fmt.Println("链接出错！", err)
		return
	}
	//启动监听
	go this.ListenMessage()
	//最后关闭连接
	defer listener.Close();
	for {
		conn, err := listener.Accept();
		if err != nil{
			fmt.Println("监听出错", err)
			continue
		}
		go this.Hanlder(conn)
	}
}