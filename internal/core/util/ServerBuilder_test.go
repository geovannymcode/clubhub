package domain

import (
	"context"
	coreDomain "github.com/Geovanny0401/clubhub/internal/core/domain"
	"testing"
)

func TestMissingEnrichmentDoesNotPanic(t *testing.T) {
	data := coreDomain.SSL{Endpoints: []coreDomain.Endpoint{{IpAddress: "127.0.0.1", Grade: ""}}}
	result := BuildServer(context.Background(), data, nil, false, "Example", "")
	if len(result.Serves) != 1 || result.Serves[0].Country != "" || result.Serves[0].Owner != "" {
		t.Fatalf("unexpected enrichment: %+v", result)
	}
}
