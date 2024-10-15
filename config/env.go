package config

import (
	"os"
	"runtime"
)

const (
	Windows = "windows"
	Linux   = "linux"
)

type Environment struct {
	RootDir  string
	Separate string
}

var Env *Environment

func LoadEnv() {

	var e Environment
	// 加载系统变量
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	e.RootDir = curDir

	if runtime.GOOS == Windows {
		e.Separate = "\\"
	} else {
		e.Separate = "/"
	}

	Env = &e
}
