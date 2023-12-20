package main

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/HumXC/shiroko/client"
	"github.com/HumXC/shiroko/example"
	"github.com/HumXC/shiroko/tools/minicap"
	"github.com/gorilla/websocket"
)

//go:embed index.html
var Html []byte

const ServeAddr = "localhost:8080"

func main() {
	ss, err := example.FindServer()
	if err != nil {
		log.Fatal(err)
	}
	client, err := client.New(ss[0].Addr)
	if err != nil {
		panic(err)
	}
	// 检查 minicap 是否安装
	err = client.Manager.Health("minicap")
	if err != nil {
		err = client.Manager.Install("minicap")
		if err != nil {
			panic(err)
		}
	}
	minicapOutput, err := GetMinicap(client.Minicap)
	if err != nil {
		panic(err)
	}
	fmt.Println("获取到minicap")
	var chFrame chan []byte = make(chan []byte, 1)
	go func() {
		ParseMinicap(minicapOutput, chFrame)
	}()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/ws", handleConnection(chFrame))
	fmt.Println("访问: http://" + ServeAddr)
	http.ListenAndServe(ServeAddr, nil)
}

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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("网页连接")
	w.Header().Set("Content-Type", "text/html")
	w.Write(Html)
}
func handleConnection(chFrame <-chan []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer w.Header().Clone()
		fmt.Println("处理连接")
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

func GetMinicap(m minicap.IMinicap) (io.Reader, error) {
	// 关闭 minicap
	err := m.Stop()
	if err != nil {
		return nil, err
	}
	// 获取屏幕信息
	info, err := m.Info()
	if err != nil {
		return nil, err
	}

	var orientation int32 = 0
	switch info.Rotation {
	case 1:
		orientation = 90
	case 2:
		orientation = 180
	case 3:
		orientation = 270
	}
	err = m.Start(info.Width, info.Height, info.Width, info.Height, orientation, 60)
	if err != nil {
		return nil, err
	}
	// 等待 minicap 启动
	time.Sleep(1 * time.Second)
	reader, err := m.Cat()
	if err != nil {
		return nil, err
	}
	return reader, nil
}
