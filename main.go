package main

import (
	"fmt"
	"os"

	"github.com/HumXC/shiroko/tools/screencap"
)

func main() {
	screencap := screencap.NewScreencap()
	dis, _ := screencap.Displays()
	b, err := screencap.PngWithDisplay(dis[0])
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("/data/local/tmp/abb.png", b, 0755)
	if err != nil {
		fmt.Println(err)
	}
}
