package utils

import (
	`fmt`
	`runtime/debug`
)

func RunGoroutine(f func(...interface{}), param ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(fmt.Sprintf("[Panic]: %v\n%s", err, string(debug.Stack())))
			}
		}()
		f(param...)
	}()
}
