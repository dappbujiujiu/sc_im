package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var _ = io.EOF
var _ = os.ErrClosed

const (
	CONNECT_NET_TYPE = "tcp"
)

var (
	ServerIp   string
	ServerPort int
)

type Client struct {
	ServerIp   string
	ServerPort int
	UserIp     int
	UserName   string
	Conn       net.Conn
	Flag       int
}

func init() {
	flag.StringVar(&ServerIp, "host", "127.0.0.1", "help-default value 127.0.0.1")
	flag.IntVar(&ServerPort, "port", 8888, "help-default value 8888")
}

func NewClient() (*Client, error) {
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		Flag:       999,
	}
	conn, err := net.Dial(CONNECT_NET_TYPE, client.ServerIp+":"+strconv.Itoa(client.ServerPort))
	if err != nil {
		return nil, errors.New("client dial err:" + err.Error())
	}
	client.Conn = conn

	return client, nil
}

//菜单方法
func (this *Client) Menu() bool {
	var flag int
	fmt.Println("1. 输入1选择公聊模式")
	fmt.Println("2. 输入2选择私聊模式")
	fmt.Println("3. 输入3选择修改用户名")
	fmt.Println("4. 输入4选择查询在线用户")
	fmt.Println("0. 输入0退出")

	fmt.Scanln(&flag)
	if flag < 0 || flag > 4 {
		fmt.Println("请输入合法数字 0-4!")
		return false
	}
	this.Flag = flag
	return true
}

//执行菜单
func (this *Client) Run() {
	for this.Flag != 0 {
		for this.Menu() != true {
		}
		switch this.Flag {
		case 1:
			this.RequestPublicChat()
			break
		case 2:
			this.RequestPrivateChat()
			break
		case 3:
			this.RequestRename()
			break
		case 4:
			this.RequestWho()
			break
		}
	}
}

//修改用户名
func (this *Client) RequestRename() bool {
	fmt.Println(">>>>>>>>请输入您的新用户名")
	fmt.Scanln(&this.UserName)
	msg := "rename|" + this.UserName + "\n"
	_, err := this.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("request rename err:", err)
		return false
	}
	return true
}

//who查询在线用户
func (this *Client) RequestWho() bool {
	msg := "who\n"
	_, err := this.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("request who err:", err)
		return false
	}
	return true
}

//公聊模式
func (this *Client) RequestPublicChat() {
	var sendMsg string
	fmt.Println(">>>>>>>>欢迎进入公聊模式")
	fmt.Scanln(&sendMsg)

	for sendMsg != "quit" {
		if len(sendMsg) != 0 {
			_, err := this.Conn.Write([]byte(sendMsg + "\n"))
			if err != nil {
				fmt.Println("request public chat err:", err)
			}
		}
		sendMsg = "" //重置发送内容
		fmt.Scanln(&sendMsg)
	}
}

//私聊模式
func (this *Client) RequestPrivateChat() {
	var selectUid int //选择的用户uid
	fmt.Println(">>>>>>>>欢迎进入私聊模式\n请输入要私聊的用户uid，进行选择")
	this.RequestWho()
	fmt.Scanln(&selectUid)

	fmt.Println(">>>>>>>>请输入聊天内容....,quit退出")
	var content string //私聊内容
	//注意fmt.scan相关函数有个问题，空格结束输入，会算成新的输入
	// fmt.Scanln(&content)
	var reader *bufio.Reader = bufio.NewReader(os.Stdin)
	content, _ = reader.ReadString('\n')
	content = strings.TrimSpace(content)
	for content != "quit" {
		if len(content) > 0 {
			sendMsg := "to|" + strconv.Itoa(selectUid) + "|" + content + "\n"
			_, err := this.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("request private chat err:", err)
			}
		}
		// content = ""
		// fmt.Scanln(&content)
		var reader *bufio.Reader = bufio.NewReader(os.Stdin)
		content, _ = reader.ReadString('\n')
		content = strings.TrimSpace(content)
	}
}

//处理返回数据
func (this *Client) DealResponse() {
	//阻塞监听 io从read中copy到标准输出流中
	io.Copy(os.Stdout, this.Conn)

	//下面作用同上
	// for {
	// 	buf := make([]byte, 4096)
	// 	_, err := this.Conn.Read(buf)
	// 	if err != nil {
	// 		fmt.Println("response read err:", err)
	// 	}
	// 	fmt.Println(string(buf))
	// }
}

func main() {
	flag.Parse()
	client, err := NewClient()
	if err != nil {
		fmt.Println(">>>>>>>>服务器链接失败 [", err, "]")
		return
	}
	fmt.Println(">>>>>>>>服务器连接成功")

	go client.DealResponse()
	//阻塞住
	// select {}
	client.Run()
}
