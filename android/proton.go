package android

import "os/exec"

func Getprop(key string) (string, error) {
	cmd := exec.Command("getprop", key)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
