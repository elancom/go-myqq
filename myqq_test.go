package myqq

import (
	"fmt"
	"github.com/elancom/go-util/str"
	"log"
	"testing"
	"time"
)

var api = NewApi("http://localhost:10008/MyQQHTTPAPI")

func TestGetQQList(t *testing.T) {
	list, err := api.GetQQList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(list)
}

func TestGetOnlineQQList(t *testing.T) {
	list, err := api.GetOnlineQQList()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(list)
}

func TestSendMsg(t *testing.T) {
	err := api.SendMsg("372666003", MsgTypeFriend, "", "272926206", "你好呀"+str.String(time.Now().UnixMilli()))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("发送成功")
}

func TestSendMsg1(t *testing.T) {
	err := api.SendMsg("372666003", MsgTypeGroup, "608084817", "", "你好呀"+str.String(time.Now().UnixMilli()))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("发送群消息成功")
}

func TestSearchGroup(t *testing.T) {
	list, err := api.SearchGroup("372666003", "健身", 6)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("长度:", len(list))
	for _, qs := range list {
		fmt.Printf("%+v\n", *qs)
	}
	fmt.Println("完成")
	time.Sleep(time.Second)
}
