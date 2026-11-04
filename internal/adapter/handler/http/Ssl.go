package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
	"github.com/Geovanny0401/clubhub/internal/security"
)

func GetDataSSl(ctx context.Context, address string) (coreDomain.SSL, error) {
	var result coreDomain.SSL
	body, err := security.Fetch(ctx, "https://api.ssllabs.com/api/v3/analyze?"+url.Values{"host": {address}}.Encode())
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	if len(result.Endpoints) > 100 {
		return coreDomain.SSL{}, errors.New("too many upstream endpoints")
	}
	return result, nil
}
