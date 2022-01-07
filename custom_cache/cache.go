package custom_cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var Ca *cache.Cache

func Init() {
	Ca = cache.New(5*time.Minute, 10*time.Minute)

	fmt.Println("cache init completed.")
}
