package main

import (
	"fmt"
	"log"
	"time"

	"github.com/HumXC/shiroko/client"
)

func main() {
	target := "192.168.3.204:15600"
	client, err := client.New(target)
	if err != nil {
		log.Fatal(err)
	}
	// info, _ := client.Minicap.Info()
	screencap := func(i int) error {
		// return client.Minicap.Jpg(info.Width, info.Height, info.Width, info.Height, 0, 100)
		_, err := client.Screencap.Png("")
		// b, _ := io.ReadAll(r)
		// os.WriteFile(fmt.Sprintf("./%d.png", i), r, 0644)
		return err
	}
	befor := time.Now()
	// wg := &sync.WaitGroup{}
	count := 10
	for i := 0; i < count; i++ {
		// wg.Add(1)
		// go func(wg *sync.WaitGroup) {

		// 	wg.Done()

		// }(wg)
		err := screencap(i)
		if err != nil {
			log.Fatal(err)
		}
	}
	// wg.Wait()
	// 输出运行时间

	fmt.Println(time.Since(befor))

}
