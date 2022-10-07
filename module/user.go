package module

import (
	"math/rand"
	"fmt" //前面加 _ 只执行init函数不回真正的导入包
	"io/fs"
	"net"
	"reflect"
	"strconv"
	"strings"

	model "github.com/dappbujiujiu/sc_im/model"
)

var _ = fs.ErrExist //防止引入包未使用报错
var _ = reflect.Array

var NowUserId int                //全局uid
var CommandMap map[string]string //处理用户的消息命令映射 key => funcName

type User struct {
	Uid    int         //用户id
	Name   string      //客户端名称
	Phone  string      //手机号
	Addr   string      //客户端地址
	C      chan string //消息管道
	Conn   net.Conn    //socket链接
	Server *ServerObj  //客户端锁链接的server
}

func init() {
	CommandMap = make(map[string]string, 2)
	CommandMap["who"] = "commandWho"
	CommandMap["rename"] = "commandRename"
}

//注册一个新用户
func UserRegister() (uid int, uName string, uPhone string, uSex string){
	m := &model.User{}
	// phone := rand.Int63n(19999999999) + 10000000000
	phone := rand.Int63n(19999999999-10000000000) + 10000000000
	uName = strconv.Itoa(int(phone))
	uPhone = uName
	sexInt := rand.Intn(2) + 1
	sexMap := map[int]string {
		1: "m",
		2: "f",
	}
	uSex = sexMap[sexInt]
	uid, err, _ := m.Register(uName, uPhone, uSex)
	if err != nil {
		fmt.Println("register new user err:", err)
	}
	return 
}

//创建user对象api
func NewUser(conn net.Conn, server *ServerObj) *User {
	uid, uName, _, _ := UserRegister()
	NowUserId += 1
	user := &User{
		Uid:    uid,
		Name:   uName,
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		Conn:   conn,
		Server: server,
	}

	//启动goroutine，给客户端发送消息
	go user.ListenServerMessage()
	return user
}

//监听server发送给user的消息
func (this *User) ListenServerMessage() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg))
	}
}

//用户上线
func (this *User) Online() {
	server := this.Server
	//把user对象指针写入到server的mapping中
	server.mapLock.Lock()
	server.OnlineMap[this.Uid] = this
	server.mapLock.Unlock()

	//进行新用户上线广播通知
	server.BroadCast(this, "已上线")
}

//用户下线
func (this *User) Offline() {
	server := this.Server
	server.mapLock.Lock()
	delete(server.OnlineMap, this.Uid)
	server.mapLock.Unlock()
	server.BroadCast(this, "下线")
}

//向单个用户发送消息
func (this *User) SendMessage(msg string) {
	this.Conn.Write([]byte(msg))
}

//处理用户发送的消息
func (this *User) ExecMessage(msg string) {
	server := this.Server
	next := this.checkMessageCommand(msg)
	if next {
		server.BroadCast(this, msg) //用户的普通消息进行广播
	}
}

//命令逻辑-who
func (this *User) commandWho(msg string, next *bool) {
	if msg != "who" {
		return
	}
	*next = false
	server := this.Server
	var tmpSmg string
	for _, v := range server.OnlineMap {
		tmpSmg += "[" + strconv.Itoa(v.Uid) + "] " + v.Name + "在线\n" //拼接要发送的消息(要注意长度)
	}
	this.SendMessage(tmpSmg)
}

//命令逻辑-rename
func (this *User) commandRename(msg string, next *bool) {
	if len(msg) <= 7 || msg[:7] != "rename|" {
		return
	}
	*next = false
	server := this.Server
	// newName := msg[7:]
	newName := strings.Split(msg, "|")[1] //通过strings.Split也可以实现获取新用户名
	//验证用户名是否重复
	for uid, user := range server.OnlineMap {
		if uid != this.Uid && newName == user.Name {
			this.SendMessage("用户名:" + newName + " 已存在，修改失败...\n")
			return
		}
	}
	this.Name = newName
	server.mapLock.Lock()
	server.OnlineMap[this.Uid] = this
	server.mapLock.Unlock()
	this.SendMessage("您的用户名已经修改成:" + newName + "\n")
}

//命令逻辑-to (私聊)
func (this *User) commandTo(msg string, next *bool) {
	if len(msg) <= 3 || msg[:3] != "to|" {
		return
	}
	*next = false
	server := this.Server
	//格式   to|uid|content
	splitStrings := strings.Split(msg, "|")
	toUserId := splitStrings[1]
	content := this.Name + " 对您说: " + splitStrings[2] + "\n"
	uid, _ := strconv.ParseInt(toUserId, 10, 64)
	toUser := server.OnlineMap[int(uid)]
	toUser.SendMessage(content)
}

//处理非正常消息的命令
func (this *User) checkMessageCommand(msg string) (next bool) {
	next = true
	this.commandWho(msg, &next)    //查询当前有谁在线
	this.commandRename(msg, &next) //修改用户名
	this.commandTo(msg, &next)     //私聊
	return
}