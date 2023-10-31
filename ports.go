package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	portsBasePathV1        = "/v1/ports"
	portsAllowAddressPairs = "allow_address_pairs"
	portsEnableSecurity    = "enable_port_security"
	portsDisableSecurity   = "disable_port_security"
)

// PortsService is an interface for creating and managing Ports with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/ports
type PortsService interface {
	Assign(context.Context, string, *AllowedAddressPairsRequest) (*Port, *Response, error)
	EnablePortSecurity(context.Context, string) (*InstancePortInterface, *Response, error)
	DisablePortSecurity(context.Context, string) (*InstancePortInterface, *Response, error)
}

// PortsServiceOp handles communication with Ports methods of the EdgecenterCloud API.
type PortsServiceOp struct {
	client *Client
}

var _ PortsService = &PortsServiceOp{}

// Port represents an EdgecenterCloud Port.
type Port struct {
	NetworkID           string                `json:"network_id"`
	AllowedAddressPairs []AllowedAddressPairs `json:"allowed_address_pairs"`
	InstanceID          string                `json:"instance_id"`
	PortID              string                `json:"port_id"`
}

// AllowedAddressPairs represents allowed port address pair and/or subnet masks.
type AllowedAddressPairs struct {
	IPAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
}

// AllowedAddressPairsRequest represents a request to assign allowed address pairs for an instance port.
type AllowedAddressPairsRequest struct {
	IPAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
}

// InstancePortInterface represents an instance port interface.
type InstancePortInterface struct {
	FloatingIPDetails   []FloatingIP `json:"floatingip_details"`
	NetworkDetails      Network      `json:"network_details"`
	PortSecurityEnabled bool         `json:"port_security_enabled"`
	PortID              string       `json:"port_id"`
	MacAddress          string       `json:"mac_address"`
	NetworkID           string       `json:"network_id"`
	IPAssignments       []PortIP     `json:"ip_assignments"`
}

// PortIP represents an IPAddress and a SubnetID.
type PortIP struct {
	IPAddress net.IP `json:"ip_address"`
	SubnetID  string `json:"subnet_id"`
}

// Assign allowed address pairs for an instance port.
func (s *PortsServiceOp) Assign(ctx context.Context, portID string, allowedAddressPairsRequest *AllowedAddressPairsRequest) (*Port, *Response, error) {
	if allowedAddressPairsRequest == nil {
		return nil, nil, NewArgError("allowedAddressPairsRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(portsBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, portID, portsAllowAddressPairs)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, allowedAddressPairsRequest)
	if err != nil {
		return nil, nil, err
	}

	port := new(Port)
	resp, err := s.client.Do(ctx, req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// EnablePortSecurity for an instance interface.
func (s *PortsServiceOp) EnablePortSecurity(ctx context.Context, portID string) (*InstancePortInterface, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(portsBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, portID, portsEnableSecurity)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	instancePortInterface := new(InstancePortInterface)
	resp, err := s.client.Do(ctx, req, instancePortInterface)
	if err != nil {
		return nil, resp, err
	}

	return instancePortInterface, resp, err
}

// DisablePortSecurity for an instance interface.
func (s *PortsServiceOp) DisablePortSecurity(ctx context.Context, portID string) (*InstancePortInterface, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(portsBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, portID, portsDisableSecurity)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	instancePortInterface := new(InstancePortInterface)
	resp, err := s.client.Do(ctx, req, instancePortInterface)
	if err != nil {
		return nil, resp, err
	}

	return instancePortInterface, resp, err
}
