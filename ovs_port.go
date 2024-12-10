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

// OvsPort represents an OVS bridge. The data help by the data
// structure is the same as the output of `ovs-vsctl list Port`
// command.
type OvsPort struct {
	UUID            string
	Name            string
	BondActiveSlave []string // TODO: unverified data type
	BondDowndelay   float64
	BondFakeIface   bool
	BondMode        []string // TODO: unverified data type
	BondUpdelay     float64
	Cvlans          []string // TODO: unverified data type
	ExternalIDs     map[string]string
	FakeBridge      bool
	Interfaces      []string
	Lacp            []string          // TODO: unverified data type
	Mac             []string          // TODO: unverified data type
	OtherConfig     map[string]string // TODO: unverified data type
	Protected       bool
	Qos             []string           // TODO: unverified data type
	RstpStatistics  map[string]float64 // TODO: unverified data type
	RstpStatus      map[string]string  // TODO: unverified data type
	Statistics      map[string]float64 // TODO: unverified data type
	Status          map[string]string  // TODO: unverified data type
	Tag             []string           // TODO: unverified data type
	Trunks          []string           // TODO: unverified data type
	VlanMode        []string           // TODO: unverified data type
}

// GetDbPorts returns a list of ports from the Port table of OVS database.
func (cli *OvsClient) GetDbPorts() ([]*OvsPort, error) {
	ports := []*OvsPort{}
	query := "SELECT * FROM Port"
	result, error := cli.Database.Vswitch.Client.Transact(cli.Database.Vswitch.Name, query)
	if error != nil {
		return nil, fmt.Errorf("the '%s' query failed: %s", query, error)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("the '%s' query did not return any rows", query)
	}
	for _, row := range result.Rows {
		port := &OvsPort{}
		data, dataType, error := row.GetColumnValue("_uuid", result.Columns)
		if error != nil {
			return nil, fmt.Errorf("couldn't get port '_uuid': %s", error)
		}
		if dataType != "string" {
			return nil, fmt.Errorf("port '_uuid' is not string")
		}
		port.UUID = data.(string)

		if data, dataType, error := row.GetColumnValue("name", result.Columns); error == nil {
			if dataType == "string" {
				port.Name = data.(string)
			} else {
				return nil, fmt.Errorf("'name' of port %s is not string", port.UUID)
			}
		}
		if data, dataType, error := row.GetColumnValue("interfaces", result.Columns); error == nil {
			if dataType == "[]string" {
				port.Interfaces = data.([]string)
			} else if dataType == "string" {
				port.Interfaces = append(port.Interfaces, data.(string))
			} else {
				return nil, fmt.Errorf("'interfaces' of port %s is not array of strings and not string", port.UUID)
			}
		}
		ports = append(ports, port)
	}
	return ports, nil

}

// GetAllPortsToInterfaces returns map of all ports, where every port uuid has its all interfaces
func (cli *OvsClient) GetAllPortsToInterfaces() (map[string][]string, error) {
	ports, error := cli.GetDbPorts()
	if error != nil {
		return nil, fmt.Errorf("couldn't get list of ports: %s", error)
	}
	allPortsToInterfaces := make(map[string][]string)
	for _, port := range ports {
		allPortsToInterfaces[port.UUID] = port.Interfaces
	}
	return allPortsToInterfaces, nil
}
