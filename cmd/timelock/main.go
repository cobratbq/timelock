package main

import (
	"github.com/cobratbq/timelock"
)

func main() {
	timelock.Timelock([]byte("This is the world's most secure time-lock encryption mechanism. Don't tell anyone!"), 3, 2)
}
