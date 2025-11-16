package logging

import "log"

func Err(op string, err error) {
	log.Printf("ERROR %s: %v", op, err)
}

func Info(op string, msg string) {
	log.Printf("INFO %s: %s", op, msg)
}

func Debug() {}
