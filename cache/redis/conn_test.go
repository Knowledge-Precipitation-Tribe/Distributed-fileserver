package redis

import (
	"fmt"
	"testing"
)

func TestRedisPool(t *testing.T) {
	pool := RedisPool()
	fmt.Println(pool)
}