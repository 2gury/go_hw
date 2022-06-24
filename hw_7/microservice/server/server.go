package main

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	StartMyMicroservice(ctx, "127.0.0.1:8082", "")
	time.Sleep(120 * time.Second)
	cancel()
}
