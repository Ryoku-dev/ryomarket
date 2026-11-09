package catalog

import (
	"github.com/Ryoku-dev/ryomarket/internal/model"
	"os"
	"runtime"
	"strings"
)

const (
	EditionDefault   = "default"
	EditionHyperion  = "hyperion"
	HyperionFlakeRef = "github:Ryoku-dev/Ryoku-on-NixOS/feature/hyperion-v2-nixos"
	hyperionHomepage = "https://github.com/Ryoku-dev/Ryoku-on-NixOS/tree/feature/hyperion-v2-nixos"
)

// Edition reports the runtime edition selected by RYOMARKET_EDITION.
func Edition() string {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("RYOMARKET_EDITION")), EditionHyperion) {
		return EditionHyperion
	}
	return EditionDefault
}

type ryokuEntry struct{ id, name, summary, description, homepage string }

var ryokuEntries = []ryokuEntry{
	{"ryotunes", "Ryotunes", "Your music, at home in Ryoku.", "Bring your local library and Spotify listening together, settle into Focus Mode, and shape your sound with the built-in equaliser.", "https://github.com/Ryoku-dev/RyotunesNix"},
	{"ryomanager", "RyoManager", "See what your system is doing.", "Monitor running processes and resource usage, find demanding applications, and keep an eye on your machine through a task manager designed for Ryoku.", "https://github.com/aethctl/ryomanager"},
	{"ryostore", "RyoStore", "Make your Ryoku desktop your own.", "Explore subshells, desktop components and customisations, then manage supported additions from one dedicated storefront.", "https://github.com/Ryoku-dev/ryoku"},
	{"ryomotion", "RyoMotion", "Turn screen recordings into polished demos.", "Record and refine walkthroughs, feature showcases and quick explanations with recording and editing tools in one place.", "https://github.com/neur0map/ryomotion"},
	{"rashin", "Rashin", "A local assistant for your Ryoku desktop.", "Ask questions, inspect activity and work with agent tools through a companion window and dashboard connected to your configured model provider.", "https://github.com/Ryoku-dev/ryoku"},
	{"ryoport", "Ryoport", "Your virtual machines and remote hosts, together.", "Monitor local VMs, manage remote servers and access SSH connections from a single fleet dashboard.", "https://github.com/Ryoku-dev/ryoku"},
}

// hyperionOutputs lists the optional packages.x86_64-linux outputs of HyperionFlakeRef.
// Ryogami is part of the Hyperion core and is intentionally absent.
var hyperionOutputs = []string{"ryotunes", "ryostore", "ryoport", "ryomanager", "ryomotion", "rashin"}

func RyokuApplications() []model.Application {
	if Edition() == EditionHyperion {
		return hyperionApplications()
	}
	ref := os.Getenv("RYOMARKET_APPS_FLAKE")
	apps := make([]model.Application, 0, len(ryokuEntries))
	for _, e := range ryokuEntries {
		app := model.Application{ID: "ryoku:" + e.id, Name: e.name, Summary: e.summary, Description: e.description, Developer: "Ryoku", Homepage: e.homepage, Categories: []string{"Ryoku"}, Ecosystem: &model.EcosystemInfo{Status: "Upcoming", Reason: "Run the packaged RyoMarket to use the bundled application recipes.", PackageOutput: "packages.x86_64-linux." + e.id}}
		if ref != "" && runtime.GOARCH == "amd64" {
			app.Ecosystem.Status = "Available"
			app.Ecosystem.Reason = "Packaged independently for your Nix profile."
			notes := ryokuNotes(e.id)
			program := e.id
			if e.id == "ryostore" {
				program = "ryostore-ui"
			}
			if e.id == "rashin" {
				program = "rashin-app"
			}
			app.Variants = map[model.Distro][]model.Variant{model.DistroNixOS: {
				{Kind: model.ProviderNixFlake, Package: e.id, FlakeRef: ref, Entrypoint: program, Official: true, Notes: notes, Recommended: true},
				{Kind: model.ProviderNixProfile, Package: e.id, FlakeRef: ref, ProfileName: e.id, Official: true, Notes: notes},
				{Kind: model.ProviderNixConfig, Package: e.id, FlakeRef: ref, Official: true, Notes: notes},
			}}
		}
		apps = append(apps, app)
	}
	return apps
}

func ryokuNotes(id string) string {
	switch id {
	case "ryotunes":
		return "Local music is ready to use. Spotify requires the optional Soloist backend setup."
	case "rashin":
		return "The launcher starts its bundled backend. Configure a model provider to use the assistant."
	case "ryoport":
		return "Local VMs need KVM access. GPU passthrough and remote hosts need separate system or credential setup."
	}
	return ""
}

// hyperionApplications installs only through nix profile, the one transaction that honours FlakeRef.
func hyperionApplications() []model.Application {
	byID := make(map[string]ryokuEntry, len(ryokuEntries))
	for _, e := range ryokuEntries {
		byID[e.id] = e
	}
	apps := make([]model.Application, 0, len(hyperionOutputs))
	for _, id := range hyperionOutputs {
		e := byID[id]
		notes := ryokuNotes(id)
		if id == "ryotunes" {
			e.summary = "Your music, in the Hyperion desktop."
			e.description = "The Ryotunes QML launcher built for Ryoku Hyperion v2, packaged from the Hyperion NixOS flake."
			notes = ""
		}
		app := model.Application{ID: "ryoku:" + id, Name: e.name, Summary: e.summary, Description: e.description, Developer: "Ryoku", Homepage: hyperionHomepage, Categories: []string{"Ryoku"}, Ecosystem: &model.EcosystemInfo{Status: "Upcoming", Reason: "Hyperion packages currently target x86_64-linux.", PackageOutput: "packages.x86_64-linux." + id}}
		if runtime.GOARCH == "amd64" {
			app.Ecosystem.Status = "Available"
			app.Ecosystem.Reason = "Packaged by the Ryoku Hyperion v2 NixOS flake."
			app.Variants = map[model.Distro][]model.Variant{model.DistroNixOS: {
				{Kind: model.ProviderNixProfile, Package: id, FlakeRef: HyperionFlakeRef, ProfileName: id, Label: "Hyperion v2", Official: true, Recommended: true, Notes: notes},
			}}
		}
		apps = append(apps, app)
	}
	return apps
}
