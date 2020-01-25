package main

import (
	"github.com/cobratbq/timelock"
)

func main() {
	// FIXME flags for specifying complexity, iterations, payload (text/stdin).
	timelock.Timelock([]byte("This is the world's \"most secure\" time-lock encryption mechanism. Don't tell anyone!"), 3, 2)
}
