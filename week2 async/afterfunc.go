package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello World")
}

func main() {
	timer := time.AfterFunc(1*time.Second, sayHello)

	_, _ = fmt.Scanln()
	timer.Stop()

	_, _ = fmt.Scanln()
}
