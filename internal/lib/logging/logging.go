package logging

import "log"

func Err(op string, err error) {
	log.Printf("ERROR %s: %v", op, err)
}

func Info() {}

func Debug() {}
