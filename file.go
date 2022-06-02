package main

import (
	"ccat/log"
	"os"
	"syscall"
)

func fileOpen(path string, lock bool) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	log.Debugln("opened ", file.Name())
	if lock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
		if err != nil {
			file.Close()
			return nil, err
		}
		log.Println("locked ", file.Name())
	}
	return file, nil
}

func fileClose(file *os.File, unlock bool) {
	if unlock {
		log.Debugln("unlock ", file.Name())
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		if err != nil {
			log.Println(err)
		}
	}
	log.Debugln("close ", file.Name())
	err := file.Close()
	if err != nil {
		log.Println(err)
	}
}
