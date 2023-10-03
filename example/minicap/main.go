package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/HumXC/shiroko/client"
	"github.com/HumXC/shiroko/tools/minicap"
	"github.com/gorilla/websocket"
)

// GlobalHeader represents the global header from minicap.
type GlobalHeader struct {
	Version       uint8
	HeaderSize    uint8
	PID           uint32
	RealWidth     uint32
	RealHeight    uint32
	VirtualWidth  uint32
	VirtualHeight uint32
	Orientation   uint8
	QuirkBitflags uint8
}

// FrameHeader represents the header of a frame from minicap.
type FrameHeader struct {
	FrameSize uint32
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var minicapOutput io.Reader

func handleConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("处理连接")
	*&flagConnected = true
	defer func() {
		*&flagConnected = false
	}()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		frame := <-chFrame
		err := conn.WriteMessage(websocket.BinaryMessage, frame)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

var flagConnected bool = false
var chFrame chan []byte = make(chan []byte, 1)

func main() {
	target := "192.168.3.204:15600"
	client, err := client.New(target)
	if err != nil {
		panic(err)
	}
	minicapOutput = GetMinicap(client.Minicap)
	go func() {
		ParseMinicap(minicapOutput, chFrame)
	}()
	go func() {
		HandleMinicap(chFrame, &flagConnected)
	}()
	fmt.Println("获取到minicap")
	http.HandleFunc("/ws", handleConnection)
	http.ListenAndServe(":8080", nil)
}
func HandleMinicap(ch <-chan []byte, flagConnected *bool) {
	for {
		if !*flagConnected {
			select {
			case _ = <-ch:
				fmt.Print(".")
			default:
			}
		}
	}
}
func ParseMinicap(reader io.Reader, ch chan<- []byte) {
	var globalHeader GlobalHeader
	if err := binary.Read(reader, binary.LittleEndian, &globalHeader); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Global Header: %+v\n", globalHeader)
	for {
		var frameHeader FrameHeader
		if err := binary.Read(reader, binary.LittleEndian, &frameHeader); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		frameData := make([]byte, frameHeader.FrameSize)
		if _, err := io.ReadFull(reader, frameData); err != nil {
			panic(err)
		}
		ch <- frameData
	}
}
func GetMinicap(m minicap.IMinicap) io.Reader {
	err := m.Stop()
	if err != nil {
		panic(err)
	}
	info, err := m.Info()
	if err != nil {
		log.Fatal(err)
	}
	err = m.Start(info.Width, info.Height, info.Width, info.Height, 0, 60)
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	reader, err := m.Cat()
	if err != nil {
		panic(err)
	}
	return reader
}
