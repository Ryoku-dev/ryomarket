package catalog

import (
	"runtime"
	"strings"
	"testing"

	"github.com/Ryoku-dev/ryomarket/internal/model"
)

func TestEditionSelection(t *testing.T) {
	for value, want := range map[string]string{"": EditionDefault, "default": EditionDefault, "hyperion-v1": EditionDefault, "hyperion": EditionHyperion, " Hyperion ": EditionHyperion} {
		t.Setenv("RYOMARKET_EDITION", value)
		if got := Edition(); got != want {
			t.Errorf("RYOMARKET_EDITION=%q: got %s, want %s", value, got, want)
		}
	}
}

func TestDefaultEditionIgnoresHyperion(t *testing.T) {
	t.Setenv("RYOMARKET_EDITION", "")
	t.Setenv("RYOMARKET_APPS_FLAKE", "path:/nix/store/packaged-source?narHash=sha256-test")
	for _, app := range RyokuApplications() {
		for _, v := range app.Variants[model.DistroNixOS] {
			if v.FlakeRef == HyperionFlakeRef {
				t.Fatalf("default edition offered Hyperion source for %s", app.ID)
			}
		}
	}
}

func TestHyperionApplications(t *testing.T) {
	t.Setenv("RYOMARKET_EDITION", "hyperion")
	t.Setenv("RYOMARKET_APPS_FLAKE", "path:/nix/store/packaged-source?narHash=sha256-test")
	apps := RyokuApplications()
	want := map[string]bool{"ryoku:ryotunes": true, "ryoku:ryostore": true, "ryoku:ryoport": true, "ryoku:ryomanager": true, "ryoku:ryomotion": true, "ryoku:rashin": true}
	if len(apps) != len(want) {
		t.Fatalf("expected %d Hyperion applications, got %d", len(want), len(apps))
	}
	seen := map[string]bool{}
	for _, app := range apps {
		if seen[app.ID] {
			t.Fatalf("duplicate id %s", app.ID)
		}
		seen[app.ID] = true
		if !want[app.ID] {
			t.Fatalf("unexpected application %s", app.ID)
		}
		if app.ID == "ryoku:ryogami" {
			t.Fatal("essential Ryogami must not become optional")
		}
		if err := validateApp(app); err != nil {
			t.Fatalf("%s: %v", app.ID, err)
		}
		for _, text := range []string{app.Homepage, app.Description, app.Ecosystem.Reason} {
			if strings.Contains(text, "RyotunesNix") || strings.Contains(text, "1.2.3") {
				t.Fatalf("%s references legacy Ryotunes: %q", app.ID, text)
			}
		}
		if !strings.Contains(app.Homepage, "Ryoku-on-NixOS/tree/feature/hyperion-v2-nixos") {
			t.Fatalf("%s homepage is not the Hyperion branch: %s", app.ID, app.Homepage)
		}
		if runtime.GOARCH != "amd64" {
			if len(app.Variants) != 0 {
				t.Fatalf("%s offered installation on unsupported architecture", app.ID)
			}
			continue
		}
		variants := app.Variants[model.DistroNixOS]
		if len(variants) != 1 {
			t.Fatalf("%s: expected one installation choice, got %d", app.ID, len(variants))
		}
		v := variants[0]
		name := strings.TrimPrefix(app.ID, "ryoku:")
		if v.Kind != model.ProviderNixProfile || v.FlakeRef != HyperionFlakeRef || v.Package != name || v.ProfileName != name || !v.Official {
			t.Fatalf("%s: unexpected variant %#v", app.ID, v)
		}
		if app.Ecosystem.PackageOutput != "packages.x86_64-linux."+name {
			t.Fatalf("%s: package output %s", app.ID, app.Ecosystem.PackageOutput)
		}
	}
}

func TestHyperionFlakeRefIsBranch(t *testing.T) {
	if HyperionFlakeRef != "github:Ryoku-dev/Ryoku-on-NixOS/feature/hyperion-v2-nixos" {
		t.Fatalf("unexpected Hyperion flake ref %s", HyperionFlakeRef)
	}
	if strings.Contains(HyperionFlakeRef, "RyotunesNix") {
		t.Fatal("Hyperion must not use the legacy Ryotunes flake")
	}
}

func TestHyperionCatalogLoads(t *testing.T) {
	t.Setenv("RYOMARKET_EDITION", "hyperion")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range hyperionOutputs {
		if _, ok := c.App("ryoku:" + id); !ok {
			t.Fatalf("missing detail page for %s", id)
		}
	}
}
