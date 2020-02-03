.PHONY: all
all: timelock unlock

timelock: library
	go build ./cmd/timelock

unlock: library
	go build ./cmd/unlock

.PHONY: library
library:
	go build ./...

.PHONY: clean
clean:
	rm timelock unlock
