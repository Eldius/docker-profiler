package docker

import (
	"context"
	"testing"
)

func TestDebug(t *testing.T) {
	c, _ := NewClient()
	c.GetRuntimeStatistcs(context.Background())
}
