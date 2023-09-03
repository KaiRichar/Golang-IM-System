package main

import(
	"net"
	"strings"
	"fmt"
)

type User struct{
	Name string
	Addr string
	conn net.Conn
	C chan string
	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User{
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		conn: conn,
		C: make(chan string),
		Server: server,
	}
	go user.ListenMsg()
	return user
}

func (this *User) ListenMsg(){
	// for msg := range this.C {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) OnLine(){
	this.Server.mapLock.Lock()
	this.Server.OnlineUserMap[this.Name] = this
	this.Server.mapLock.Unlock()
	this.Server.BroadCast(this, "上线了")
}

func (this *User) OffLine(){
	this.Server.mapLock.Lock()
	delete(this.Server.OnlineUserMap, this.Name)
	this.Server.mapLock.Unlock()
	this.Server.BroadCast(this, "已下线")
}

func (this *User) SendMsg(msg string){
	fmt.Println("向server发送消息" + msg)
	_, err := this.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("发送消息到server出错")
	}
}

func (this *User) DoMessage(msg string){
	// fmt.Println("调用user[" + this.Name + "]的DoMessage")
	if "who" == msg {
		this.Server.mapLock.Lock()
		sendMsg := ""
		for _, user := range this.Server.OnlineUserMap {
			sendMsg = "用户: " + user.Name + "["+ user.Addr +"]" + "在线...\n"  
			this.SendMsg(sendMsg)
		}
		this.Server.mapLock.Unlock()
	}else if len(msg) > 7 && "rename|" == msg[:7]{
		newName := strings.Split(msg, "|")[1]
		_, ok := this.Server.OnlineUserMap[newName]
		if ok {
			this.SendMsg("该用户名已被使用！\n")
		} else {
			this.Server.mapLock.Lock()
			delete(this.Server.OnlineUserMap, this.Name)
			this.Server.OnlineUserMap[newName] = this
			this.Server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("用户名修改成功，您的新用户名为：" + this.Name + "\n")
		}
	}else if len(msg) > 4 && "to" == strings.Split(msg, "|")[0]{
			content := strings.Split(msg, "|")[2]
			toUser := strings.Split(msg, "|")[1]
			remoteUser, ok := this.Server.OnlineUserMap[toUser]
			if !ok {
				this.SendMsg("不存在该用户\n")
			} else {
				remoteUser.SendMsg(this.Name + "对您说：" + content + "\n")
			}
	}else {
		fmt.Println("广播消息" + msg)
		this.Server.BroadCast(this, msg)
	}
}