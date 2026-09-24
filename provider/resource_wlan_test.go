package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// wlanModelForSecurity builds a fully-valid wlanModel with every field set
// explicitly, for testing validateWlan/wlanEntryFromModel security handling.
func wlanModelForSecurity(security, passphrase string) wlanModel {
	return wlanModel{
		DeviceMac:  types.StringValue("aa:bb:cc:dd:ee:ff"),
		Name:       types.StringValue("main"),
		SSID:       types.StringValue("MainNet"),
		Security:   types.StringValue(security),
		Passphrase: types.StringValue(passphrase),
		VLAN:       types.Int64Value(42),
		Enabled:    types.BoolValue(true),
		Band:       types.StringValue("both"),
	}
}

func TestValidateWlanSecurity(t *testing.T) {
	tests := []struct {
		name    string
		model   wlanModel
		wantOK  bool
		wantErr string
	}{
		{name: "wpa-p valid", model: wlanModelForSecurity("wpa-p", "longenough"), wantOK: true},
		{name: "wpa3-p valid", model: wlanModelForSecurity("wpa3-p", "longenough"), wantOK: true},
		{name: "wpa2-wpa3 valid", model: wlanModelForSecurity("wpa2-wpa3", "longenough"), wantOK: true},
		{name: "wpa3-p short passphrase", model: wlanModelForSecurity("wpa3-p", "short"), wantErr: "at least 8 characters"},
		{name: "wpa2-wpa3 short passphrase", model: wlanModelForSecurity("wpa2-wpa3", "short"), wantErr: "at least 8 characters"},
		{name: "wpa-eap unsupported", model: wlanModelForSecurity("wpa-eap", "longenough"), wantErr: "security must be one of open|wpa-p|wpa3-p|wpa2-wpa3"},
		{name: "unknown security", model: wlanModelForSecurity("wpa3", "longenough"), wantErr: "security must be one of open|wpa-p|wpa3-p|wpa2-wpa3"},
		{name: "open valid", model: wlanModelForSecurity("open", ""), wantOK: true},
		{name: "open with passphrase rejected", model: wlanModelForSecurity("open", "longenough"), wantErr: "passphrase must be empty when security is open"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d diag.Diagnostics
			ok := (&wlanResource{}).validateWlan(&tt.model, &d)
			if tt.wantErr != "" {
				if ok || !d.HasError() {
					t.Fatalf("validateWlan ok=%v hasError=%v, want failure with %q", ok, d.HasError(), tt.wantErr)
				}
				found := false
				for _, di := range d.Errors() {
					if strings.Contains(di.Summary()+di.Detail(), tt.wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("diagnostics = %+v, want detail containing %q", d.Errors(), tt.wantErr)
				}
				return
			}
			if !ok || d.HasError() {
				t.Fatalf("validateWlan ok=%v diagnostics=%v, want success", ok, d.Errors())
			}
		})
	}
}

func TestValidateWlanSecurityNullPassphrase(t *testing.T) {
	m := wlanModelForSecurity("wpa3-p", "")
	m.Passphrase = types.StringNull()
	var d diag.Diagnostics
	if ok := (&wlanResource{}).validateWlan(&m, &d); ok || !d.HasError() {
		t.Fatalf("wpa3-p null passphrase: ok=%v hasError=%v, want failure", ok, d.HasError())
	}
	found := false
	for _, di := range d.Errors() {
		if strings.Contains(di.Summary()+di.Detail(), "at least 8 characters") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("diagnostics = %+v, want detail containing %q", d.Errors(), "at least 8 characters")
	}
}

func TestWlanEntryFromModelPassphrase(t *testing.T) {
	t.Run("open clears passphrase", func(t *testing.T) {
		m := wlanModelForSecurity("open", "longenough")
		e := wlanEntryFromModel(&m)
		if e.Security != "open" {
			t.Fatalf("Security = %q, want open", e.Security)
		}
		if e.Passphrase != "" {
			t.Fatalf("Passphrase = %q, want empty", e.Passphrase)
		}
	})
	t.Run("wpa3-p preserves passphrase", func(t *testing.T) {
		m := wlanModelForSecurity("wpa3-p", "longenough")
		e := wlanEntryFromModel(&m)
		if e.Security != "wpa3-p" {
			t.Fatalf("Security = %q, want wpa3-p", e.Security)
		}
		if e.Passphrase != "longenough" {
			t.Fatalf("Passphrase = %q, want longenough", e.Passphrase)
		}
	})
	t.Run("wpa2-wpa3 preserves passphrase", func(t *testing.T) {
		m := wlanModelForSecurity("wpa2-wpa3", "longenough")
		e := wlanEntryFromModel(&m)
		if e.Security != "wpa2-wpa3" {
			t.Fatalf("Security = %q, want wpa2-wpa3", e.Security)
		}
		if e.Passphrase != "longenough" {
			t.Fatalf("Passphrase = %q, want longenough", e.Passphrase)
		}
	})
}
