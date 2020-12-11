package main

import (
	"flag"
	"fmt"
	"gosproto/auth"
	"gosproto/header"
	"gosproto/room"
	"gosproto/sproto"
	"log"
	"net/url"
	"os"
	"os/signal"
	"reflect"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "10.10.210.79:5002", "http service address")

func main() {
	fmt.Println("hello world")
	// c := connect()
	// defer c.Close()
	// testLogin(c)
	// time.Sleep(time.Duration(2) * time.Second)
	// testEnterRoom(c)
	// time.Sleep(time.Duration(2) * time.Second)
	// testGroupRequest(c)
	// time.Sleep(time.Duration(2) * time.Second)
	// testExitRoom(c)
	// time.Sleep(time.Duration(5) * time.Second)
	testDecodeAuth()
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

func testDecodeAuth() {
	//压缩
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

	//解压
	unpackEncodeHeaderReq, err := sproto.Unpack(packEncodeHeaderReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("解压后的数据：", unpackEncodeHeaderReq)
	spHeader := reflect.New(header.SProtoStructs[0]).Interface()
	// spHeader := reflect.New(reflect.TypeOf(&header.Package{}).Elem()).Interface()

	decodeHeader, err := sproto.Decode(unpackEncodeHeaderReq, spHeader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("spHeader", spHeader)

	fmt.Println("decode header length:", decodeHeader)
	fmt.Println("decode后的header：", spHeader.(*header.Package).Protoid)

	spReq := reflect.New(auth.SProtoStructs[3]).Interface()
	decodeReq, err := sproto.Decode(unpackEncodeHeaderReq[decodeHeader:], spReq)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("spReq", spReq)

	fmt.Println("decode req length:", decodeReq)
	fmt.Println("decode后的req, Token：", spReq.(*auth.PlayerLoginReq).Token)
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
