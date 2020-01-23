package main

import (
	"github.com/cobratbq/timelock"
)

func main() {
	timelock.Timelock([]byte("Hello world!"), 3, 2)
}
