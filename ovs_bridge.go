// Copyright 2018 Paul Greenberg (greenpau@outlook.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ovsdbclient

import (
	"fmt"
)

// OvsBridge represents an OVS bridge. The data help by the data
// structure is the same as the output of `ovs-vsctl list Bridge`
// command.
type OvsBridge struct {
	UUID                string
	Name                string
	AutoAttach          []string // TODO: unverified data type
	Controller          []string // TODO: unverified data type
	DatapathName        string   // reference from ovs-appctl dpif/show
	DatapathID          string
	DatapathType        string
	DatapathVersion     string
	ExternalIDs         map[string]string
	FailMode            string
	FloodVlans          []string          // TODO: unverified data type
	FlowTables          map[string]string // TODO: unverified data type
	Ipfix               []string          // TODO: unverified data type
	McastSnoopingEnable bool
	Mirrors             []string // TODO: unverified data type
	Netflow             []string // TODO: unverified data type
	OtherConfig         map[string]string
	Ports               []string
	Protocols           []string // TODO: unverified data type
	RstpEnable          bool
	RstpStatus          map[string]string // TODO: unverified data type
	Sflow               []string          // TODO: unverified data type
	Status              map[string]string // TODO: unverified data type
	StpEnable           bool
}

// GetDbBridges returns a list of bridges from the Bridge table of OVS database.
func (cli *OvsClient) GetDbBridges() ([]*OvsBridge, error) {
	bridges := []*OvsBridge{}
	query := "SELECT * FROM Bridge"
	result, error := cli.Database.Vswitch.Client.Transact(cli.Database.Vswitch.Name, query)
	if error != nil {
		return nil, fmt.Errorf("the '%s' query failed: %s", query, error)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("the '%s' query did not return any rows", query)
	}
	for _, row := range result.Rows {
		bridge := &OvsBridge{}
		data, dataType, error := row.GetColumnValue("_uuid", result.Columns)
		if error != nil {
			return nil, fmt.Errorf("couldn't get bridge '_uuid': %s", error)
		}
		if dataType != "string" {
			return nil, fmt.Errorf("bridge '_uuid' is not string")
		}
		bridge.UUID = data.(string)

		if data, dataType, error := row.GetColumnValue("name", result.Columns); error == nil {
			if dataType == "string" {
				bridge.Name = data.(string)
			} else {
				return nil, fmt.Errorf("'name' of bridge %s is not string", bridge.UUID)
			}
		}
		if data, dataType, error := row.GetColumnValue("ports", result.Columns); error == nil {
			if dataType == "[]string" {
				bridge.Ports = data.([]string)
			} else if dataType == "string" {
				bridge.Ports = append(bridge.Ports, data.(string))
			} else {
				return nil, fmt.Errorf("'ports' of bridge %s is not array of strings and not string", bridge.UUID)
			}
		}
		bridges = append(bridges, bridge)
	}
	return bridges, nil

}

// GetAllInterfacesToBridges returns all interfaces which are related with ports of the bridge
func (cli *OvsClient) GetAllInterfacesToBridges() (map[string]string, error) {
	bridges, error := cli.GetDbBridges()
	if error != nil {
		return nil, fmt.Errorf("couldn't get list of bridges: %s", error)
	}
	allPortsToInterfaces, error := cli.GetAllPortsToInterfaces()
	if error != nil {
		return nil, fmt.Errorf("couldn't get list of ports: %s", error)
	}
	allInterfacesToBridges := make(map[string]string)
	for _, bridge := range bridges {
		for _, bridgePort := range bridge.Ports {
			for _, interfaceUUID := range allPortsToInterfaces[bridgePort]{
						allInterfacesToBridges[interfaceUUID] = bridge.Name
			}
		}
	}
	return allInterfacesToBridges, nil
}
