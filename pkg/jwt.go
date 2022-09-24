package pkg

import (
	"github.com/costa92/errors"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/internal/middleware/auth"
	"github.com/costa92/go-web/pkg/cache"
)

func newCacheAuth() middleware.AuthStrategy {
	return auth.NewCacheStrategy(getSecretFunc())
}

func getSecretFunc() func(string) (auth.Secret, error) {
	return func(kid string) (auth.Secret, error) {
		cli, err := cache.GetCacheInsOr()
		if err != nil || cli == nil {
			return auth.Secret{}, errors.Wrap(err, "get cache instance failed")
		}
		err = cli.GetSecret(kid)
		if err != nil {
			return auth.Secret{}, err
		}

		return auth.Secret{
			Username: "",
			ID:       "1",
			Key:      "",
			Expires:  111,
		}, nil
	}
}
