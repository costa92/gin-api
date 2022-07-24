package db

import (
	"context"
	"testing"
)

func Test_DisableRedis(t *testing.T) {
	DisableRedis(true)
}

func TestRedisCluster_Connect(t *testing.T) {
	opts := &RedisOpts{}
	simpleOpts := opts.simple()

	if simpleOpts.Addr != "127.0.0.1:6379" {
		t.Fatal("Wrong default single node address")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// keep redis connected
	ConnectToRedis(ctx, &Config{
		Addrs: []string{"127.0.0.1:6379"},
	})
}
