package cache

import "testing"

func Test_GetSecret(t *testing.T) {
	cache, err := GetCacheInsOr()
	if err != nil {
		t.Fatal(err)
	}
	cache.SetSecret("eqweqwe", "122")

	_ = cache.GetSecret("eqweqwe")
	// t.Fatal(err)
}
