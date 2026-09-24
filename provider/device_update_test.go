package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDeviceUpdateClearsName(t *testing.T) {
	state := deviceModel{Name: types.StringValue("old"), SiteID: types.StringValue("default")}
	plan := state
	plan.Name = types.StringNull()
	body := deviceUpdateBody(plan, state)
	if got := body["name"]; got != "" {
		t.Fatalf("name update = %#v, want explicit empty name", body)
	}
	if _, ok := body["site_id"]; ok {
		t.Fatalf("unchanged site_id should be absent: %#v", body)
	}
}

func TestImportAndNormalizeHelpers(t *testing.T) {
	if got, err := normalizeMAC("AA-BB-CC-DD-EE-FF"); err != nil || got != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("normalize MAC: %q %v", got, err)
	}
	if _, err := normalizeMAC("not-a-mac"); err == nil {
		t.Fatal("invalid MAC accepted")
	}
}
