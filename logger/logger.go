package logger

import (
	"fmt"
	"os"
	"time"
)

var tags = map[string]string{
	"info":    "INFO",
	"error":   "ERR ",
	"message": "MSG ",
	"image":   "IMG ",
	"fatal":   "FATL",
}

// Info (a ...any): logs with INFO tag
func Info(a ...any) {
	log(tags["info"], a...)
}

// Error (a ...any): logs with ERR tag
func Error(a ...any) {
	log(tags["error"], a...)
}

// Message (a ...any): logs with MSG tag
func Message(a ...any) {
	log(tags["message"], a...)
}

// Image (a ...any): logs with IMG tag
func Image(a ...any) {
	log(tags["image"], a...)
}

// Fatal (a ...any): logs the error with the FATL tag, before immediately exiting the program with error code 1
func Fatal(a ...any) {
	log(tags["fatal"], a...)
	os.Exit(1)
}

func log(tag string, a ...any) {
	prefix := fmt.Sprintf("%v [%v]", time.Now().Format(time.Stamp), tag)
	args := fmt.Sprint(a...)
	fmt.Println(prefix, args)
}
