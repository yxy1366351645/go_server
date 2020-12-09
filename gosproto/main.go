package main

import (
	"flag"
	"fmt"
	"gosproto/address"
	"gosproto/auth"
	"gosproto/header"
	"gosproto/room"
	"gosproto/sproto"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "10.10.210.79:5002", "http service address")

func main() {
	fmt.Println("hello world")
	c := connect()
	defer c.Close()
	testLogin(c)
	time.Sleep(time.Duration(2) * time.Second)
	testEnterRoom(c)
	time.Sleep(time.Duration(2) * time.Second)
	testGroupRequest(c)
	time.Sleep(time.Duration(2) * time.Second)
	testExitRoom(c)
	time.Sleep(time.Duration(5) * time.Second)
}

func login(c *websocket.Conn) {
	reqHeader := &header.Package{
		Protoid:  0x0101,
		Response: 0,
	}
	encodeHeader, err := sproto.Encode(reqHeader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化后的header数据：", encodeHeader)

	req := &auth.PlayerLoginReq{
		Token: "05708b3c37212f57a68216304d721171",
	}
	encodeReq, err := sproto.Encode(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化后的data数据：", encodeReq)
	var encodeHeaderReq = append(encodeHeader, encodeReq...)
	packEncodeHeaderReq := sproto.Pack(encodeHeaderReq)

	fmt.Println("压缩后的req数据：", packEncodeHeaderReq)

	send(c, packEncodeHeaderReq)
}

func testLogin(c *websocket.Conn) {
	login(c)
}

func enterRoom(c *websocket.Conn) {
	reqHeader := &header.Package{
		Protoid:   0x0d01,
		Roomproxy: "rummy2",
		Response:  0,
	}
	encodeHeader, err := sproto.Encode(reqHeader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化后的header数据：", encodeHeader)

	packReqHeader := sproto.Pack(encodeHeader)

	fmt.Println("压缩后的req数据：", packReqHeader)

	send(c, packReqHeader)
}

func testEnterRoom(c *websocket.Conn) {
	enterRoom(c)
}

func groupRequest(c *websocket.Conn) {
	reqHeader := &header.Package{
		Protoid:   0x0d03,
		Roomproxy: "rummy2",
		Response:  0,
	}
	encodeHeader, err := sproto.Encode(reqHeader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化后的header数据：", encodeHeader)

	req := &room.GroupReq{
		Game_id: 105,
	}
	encodeReq, err := sproto.Encode(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化后的data数据：", encodeReq)
	var encodeHeaderReq = append(encodeHeader, encodeReq...)
	packEncodeHeaderReq := sproto.Pack(encodeHeaderReq)

	fmt.Println("压缩后的req数据：", packEncodeHeaderReq)

	send(c, packEncodeHeaderReq)
}

func testGroupRequest(c *websocket.Conn) {
	groupRequest(c)
}

func exitRoom(c *websocket.Conn) {
	reqHeader := &header.Package{
		Protoid:  0x0d02,
		Response: 0,
	}
	encodeHeader, err := sproto.Encode(reqHeader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化后的header数据：", encodeHeader)

	packReqHeader := sproto.Pack(encodeHeader)

	fmt.Println("压缩后的req数据：", packReqHeader)

	send(c, packReqHeader)
}

func testExitRoom(c *websocket.Conn) {
	exitRoom(c)
}

func connect() *websocket.Conn {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/message"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func send(c *websocket.Conn, req []byte) {
	err := c.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func testAddressBook() {
	input := &address.AddressBook{
		Person: []*address.Person{
			&address.Person{
				Name: "Alice",
				Id:   10000,
				Phone: []*address.PhoneNumber{
					&address.PhoneNumber{
						Number: "123456789",
						Type:   1,
					},
					&address.PhoneNumber{
						Number: "87654321",
						Type:   2,
					},
				},
			},
			&address.Person{
				Name: "Bob",
				Id:   20000,
				Phone: []*address.PhoneNumber{
					&address.PhoneNumber{
						Number: "01234567890",
						Type:   3,
					},
				},
			},
		},
	}
	data, err := sproto.Encode(input)
	fmt.Println("序列化后的数据：", data)

	if err != nil {
		fmt.Println(err)
	}

	var sample address.AddressBook
	_, err = sproto.Decode(data, &sample)

	if err != nil {
		fmt.Println(err)
	}
}
