package module

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)
var _ = strconv.ErrRange

type ServerObj struct {
	Ip        string
	Port      int
	OnlineMap map[int]*User //在线用户mapping
	mapLock   sync.RWMutex     //为在线用户映射加一个读写锁
	Message   chan string      //总消息管道
}

//进行初始化
func NewServer(ip string, port int) *ServerObj {
	server := &ServerObj{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[int]*User),
		Message:   make(chan string),
	}

	return server
}

//进行消息广播
func (this *ServerObj) BroadCast(user *User, msg string) {
	uid := strconv.Itoa(user.Uid)
	msg = "[" + uid + "]" + user.Name + " " + msg + "\n"
	// var builder strings.Builder
	// builder.Grow(50)
	// builder.WriteString("[")
	// builder.WriteRune(rune(user.Uid))
	// builder.WriteString("]" + user.Name + " " + msg + "\n")
	// msg = builder.String()
	this.Message <- msg
}

//监听server消息管道，放到全部用户管道中
func (this *ServerObj) ListenMessager() {
	for {
		msg := <-this.Message

		//别忘记加锁
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

//启动服务
func (this *ServerObj) Start() {
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen err :", err)
		return
	}

	fmt.Println("listening ip:", this.Ip, " and port:", this.Port)
	//关闭连接
	defer listener.Close()

	//监听广播消息
	go this.ListenMessager()

	//处理连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err", err)
			continue
		}

		//要并行处理链接
		go this.Handler(conn)
	}
}

//处理新用户链接
//user指针写入mapping中
//新加入用户消息金星广播
func (this *ServerObj) Handler(conn net.Conn) {
	user := NewUser(conn, this)
	//用户上线
	user.Online()

	isLive := make(chan bool)	//用户是否活跃管道
	//接受客户端消息进行广播（群聊）
	go func() {
		for{
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			//用户下线
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			//获取用户消息去除结尾\n
			msg := string(buf[:n-1])
			user.ExecMessage(msg)
			isLive <- true
		}
	}()

	//阻塞当前handler，注意一定要有真实的goroutine在跑，否则select{}会报panic
	for{
		select {
		case <-isLive:
			//什么都不用做，为了激活select，等待time.after到期后的逻辑执行
		case <-time.After(1000 * time.Second):
			//超时踢出用户
			user.SendMessage("您已超时，被踢出...\n")
			//销毁资源和关闭链接
			close(user.C)		
			conn.Close()
			return //或是 runtime.Goexit退出
		}
	}
}
