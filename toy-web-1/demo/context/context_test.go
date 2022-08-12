package context_test

import (
	"context"
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	SomeBiz()
}

func SomeBiz() {
	ctx := context.Background()
	cctx, cancel := context.WithCancel(ctx)
	cancel()
	fmt.Println(cctx.Err())
}
