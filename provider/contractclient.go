package provider

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ContractClient exposes the provider HTTP decode layer for controller-side
// integration tests (open-unifi admin API contract). It is not part of the
// supported Terraform provider surface.
type ContractClient struct {
	c *apiClient
}

// NewContractClient wires the provider client to an arbitrary RoundTripper
// (for in-process adminapi handlers in tests).
func NewContractClient(baseURL, token string, transport http.RoundTripper) *ContractClient {
	c := newAPIClient(baseURL, token, false).withHTTPClient(&http.Client{Transport: transport})
	return &ContractClient{c: c}
}

func (cc *ContractClient) CheckConnectivity(ctx context.Context) error {
	return cc.c.checkConnectivity(ctx)
}

func (cc *ContractClient) Do(ctx context.Context, method, path string, body any, out any) error {
	return cc.c.do(ctx, method, path, body, out)
}

func (cc *ContractClient) GetDevice(ctx context.Context, mac string) (*ContractDeviceView, error) {
	dv, err := cc.c.getDevice(ctx, mac)
	if err != nil {
		return nil, err
	}
	return contractDeviceViewFromWire(dv), nil
}

func (cc *ContractClient) ListDevices(ctx context.Context) ([]ContractDevice, error) {
	devs, err := cc.c.listDevices(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ContractDevice, len(devs))
	for i, d := range devs {
		out[i] = contractDeviceFromWire(d)
	}
	return out, nil
}

func (cc *ContractClient) CreateWireless(ctx context.Context, mac string, entry *ContractWirelessEntry) error {
	return cc.c.createWireless(ctx, mac, contractWirelessToWire(entry))
}

func (cc *ContractClient) GetWireless(ctx context.Context, mac, name string) (*ContractWirelessEntry, error) {
	e, err := cc.c.getWireless(ctx, mac, name)
	if err != nil {
		return nil, err
	}
	return contractWirelessFromWire(e), nil
}

func (cc *ContractClient) UpdateWireless(ctx context.Context, mac, name string, entry *ContractWirelessEntry) error {
	return cc.c.updateWireless(ctx, mac, name, contractWirelessToWire(entry))
}

// ContractDevice is the list-element wire decode shape.
type ContractDevice struct {
	Mac      string
	Name     string
	State    int
	LastSeen int64
	IP       string
}

// ContractDeviceView is the single-device GET wire decode shape.
type ContractDeviceView struct {
	Mac      string
	Name     string
	State    int
	LastSeen int64
	IP       string
}

// ContractWirelessEntry is one per-device WLAN item on the wire.
type ContractWirelessEntry struct {
	ID         string
	Name       string
	SSID       string
	Security   string
	Passphrase string
	VLAN       int
	Enabled    bool
	Band       string
}

// ContractDeviceModel mirrors device resource state for contract tests.
type ContractDeviceModel struct {
	Mac      types.String
	Name     types.String
	SiteID   types.String
	State    types.String
	IP       types.String
	LastSeen types.String
	Firmware types.String
}

// ContractStateName maps numeric server state to Terraform vocabulary.
func ContractStateName(n int) string { return stateName(n) }

// ContractErrNotFound reports HTTP 404 API errors from the client.
func ContractErrNotFound(err error) bool { return errNotFound(err) }

// ContractAPIError inspects typed API errors from Do/Get*.
func ContractAPIError(err error) (status int, notFound bool, ok bool) {
	var ae *apiError
	if !errors.As(err, &ae) {
		return 0, false, false
	}
	return ae.status, ae.NotFound(), true
}

// ContractApplyDevice copies server wire truth into a contract device model.
func ContractApplyDevice(m *ContractDeviceModel, dev *ContractDeviceView, siteID string) {
	wire := &deviceView{
		Mac: dev.Mac, Name: dev.Name, State: dev.State,
		IP: dev.IP, LastSeen: dev.LastSeen,
	}
	dm := &deviceModel{
		Mac: m.Mac, Name: m.Name, SiteID: m.SiteID,
		State: m.State, IP: m.IP, LastSeen: m.LastSeen, Firmware: m.Firmware,
	}
	applyDevice(dm, wire, siteID)
	m.SiteID = dm.SiteID
	m.State = dm.State
	m.IP = dm.IP
	m.LastSeen = dm.LastSeen
	m.Firmware = dm.Firmware
	m.Name = dm.Name
}

// ContractDeviceUpdateBody builds a PATCH body for device updates (contract tests).
func ContractDeviceUpdateBody(plan, state ContractDeviceModel) map[string]any {
	return deviceUpdateBody(deviceModel{
		Name: plan.Name, SiteID: plan.SiteID,
	}, deviceModel{
		Name: state.Name, SiteID: state.SiteID,
	})
}

// ContractDeviceWireTypes returns reflect types for device list vs GET wire structs.
func ContractDeviceWireTypes() (listElem, single reflect.Type) {
	return reflect.TypeOf(device{}), reflect.TypeOf(deviceView{})
}

// ContractNormalizeMAC is the device resource MAC normalizer.
func ContractNormalizeMAC(s string) (string, error) { return normalizeMAC(s) }

func contractDeviceFromWire(d device) ContractDevice {
	return ContractDevice{
		Mac: d.Mac, Name: d.Name, State: d.State, LastSeen: d.LastSeen, IP: d.IP,
	}
}

func contractDeviceViewFromWire(d *deviceView) *ContractDeviceView {
	if d == nil {
		return nil
	}
	return &ContractDeviceView{
		Mac: d.Mac, Name: d.Name, State: d.State, LastSeen: d.LastSeen, IP: d.IP,
	}
}

func contractWirelessFromWire(e *wirelessEntry) *ContractWirelessEntry {
	if e == nil {
		return nil
	}
	return &ContractWirelessEntry{
		ID: e.ID, Name: e.Name, SSID: e.SSID, Security: e.Security,
		Passphrase: e.Passphrase, VLAN: e.VLAN, Enabled: e.Enabled, Band: e.Band,
	}
}

func contractWirelessToWire(e *ContractWirelessEntry) *wirelessEntry {
	if e == nil {
		return nil
	}
	return &wirelessEntry{
		ID: e.ID, Name: e.Name, SSID: e.SSID, Security: e.Security,
		Passphrase: e.Passphrase, VLAN: e.VLAN, Enabled: e.Enabled, Band: e.Band,
	}
}
