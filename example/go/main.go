package main

import (
	"fmt"
	"log"
	"os"

	"github.com/HumXC/shiroko/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	target := "192.168.3.204:15600"
	client, err := client.New(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}

	ids, err := client.Screencap.Displays()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ids)
	b, err := client.Screencap.Png("")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("a.png", b, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
