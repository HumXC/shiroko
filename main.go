package main

import (
	"fmt"

	"github.com/HumXC/shiroko/android"
	"github.com/HumXC/shiroko/tools/minicap"
)

func main() {
	minicap := minicap.NewMinicap()
	err := minicap.Base.Health()
	if err != nil {
		fmt.Println(err)
		err := minicap.Base.Install()
		fmt.Println(err)
	}
	args := append(minicap.Base.Args())
	cmd := android.Command(minicap.Base.Exe(), args...)
	cmd.SetEnv(minicap.Base.Env())
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
