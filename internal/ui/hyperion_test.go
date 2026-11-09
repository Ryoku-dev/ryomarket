package ui

import (
	"strings"
	"testing"

	"github.com/Ryoku-dev/ryomarket/internal/catalog"
)

func TestRyokuCollectionFollowsEdition(t *testing.T) {
	t.Setenv("RYOMARKET_EDITION", "")
	title, ids := ryokuCollection()
	if title != "From Ryoku" || strings.Join(ids, ",") != "ryoku:ryotunes,ryoku:ryomanager,ryoku:ryostore,ryoku:ryomotion,ryoku:rashin,ryoku:ryoport" {
		t.Fatalf("default collection changed: %s %v", title, ids)
	}

	t.Setenv("RYOMARKET_EDITION", "hyperion")
	title, ids = ryokuCollection()
	if title != "From Ryoku Hyperion" || len(ids) != 6 {
		t.Fatalf("Hyperion collection: %s %v", title, ids)
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] || id == "ryoku:ryogami" {
			t.Fatalf("invalid Hyperion collection entry %s", id)
		}
		seen[id] = true
	}
}

func TestHyperionSearchUsesHyperionSource(t *testing.T) {
	t.Setenv("RYOMARKET_EDITION", "hyperion")
	a := storefrontFixture(t)
	a.query.SetText("Ryotunes")
	a.ensureDiscoverCache()
	if len(a.discoverApps) == 0 || a.discoverApps[0].ID != "ryoku:ryotunes" {
		t.Fatal("Hyperion Ryotunes was not surfaced")
	}
	for _, variants := range a.discoverApps[0].Variants {
		for _, v := range variants {
			if v.FlakeRef != catalog.HyperionFlakeRef {
				t.Fatalf("Hyperion search offered another source: %#v", v)
			}
		}
	}
}
