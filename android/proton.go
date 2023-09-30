package android

import "os/exec"

var TestProp map[string]string

func Getprop(key string) string {
	if TestProp != nil {
		return TestProp[key]
	}

	cmd := exec.Command("getprop", key)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(output[:len(output)-1])
}
