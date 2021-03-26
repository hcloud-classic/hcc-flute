package ipmi

type ipmiMembers struct {
	OdataContext      string `json:"@odata.context"`
	OdataID           string `json:"@odata.id"`
	OdataType         string `json:"@odata.type"`
	Name              string `json:"Name"`
	MembersOdataCount int    `json:"Members@odata.count"`
	Members           []struct {
		OdataID string `json:"@odata.id"`
	} `json:"Members"`
}

type ipmiNodeInfo struct {
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	ID           string `json:"ID"`
	SerialNumber string `json:"SerialNumber"`
	Name         string `json:"Name"`
	Model        string `json:"Model"`
	Manufacturer string `json:"Manufacturer"`
	PartNumber   string `json:"PartNumber"`
	AssetTag     string `json:"AssetTag"`
	SKU          string `json:"SKU"`
	SystemType   string `json:"SystemType"`
	Description  string `json:"Description"`
	UUID         string `json:"UUID"`
	Status       struct {
		State        string `json:"State"`
		Health       string `json:"Health"`
		HealthRollup string `json:"HealthRollup"`
	} `json:"Status"`
	IndicatorLED string `json:"IndicatorLED"`
	PowerState   string `json:"PowerState"`
	Boot         struct {
		BootSourceOverrideEnabled                      string   `json:"BootSourceOverrideEnabled"`
		BootSourceOverrideTarget                       string   `json:"BootSourceOverrideTarget"`
		BootSourceOverrideMode                         string   `json:"BootSourceOverrideMode"`
		BootSourceOverrideTargetRedfishAllowableValues []string `json:"BootSourceOverrideTarget@Redfish.AllowableValues"`
	} `json:"Boot"`
	BiosVersion      string `json:"BiosVersion"`
	ProcessorSummary struct {
		Count  int    `json:"Count"`
		Model  string `json:"Model"`
		Status struct {
			State        string `json:"State"`
			Health       string `json:"Health"`
			HealthRollup string `json:"HealthRollup"`
		} `json:"Status"`
	} `json:"ProcessorSummary"`
	MemorySummary struct {
		TotalSystemMemoryGiB int `json:"TotalSystemMemoryGiB"`
		Status               struct {
			State        string `json:"State"`
			Health       string `json:"Health"`
			HealthRollup string `json:"HealthRollup"`
		} `json:"Status"`
	} `json:"MemorySummary"`
	Bios struct {
		OdataID string `json:"@odata.id"`
	} `json:"Bios"`
	Processors struct {
		OdataID string `json:"@odata.id"`
	} `json:"Processors"`
	EthernetInterfaces struct {
		OdataID string `json:"@odata.id"`
	} `json:"EthernetInterfaces"`
	LogServices struct {
		OdataID string `json:"@odata.id"`
	} `json:"LogServices"`
	Memory struct {
		OdataID string `json:"@odata.id"`
	} `json:"Memory"`
	Storage struct {
		OdataID string `json:"@odata.id"`
	} `json:"Storage"`
	Links struct {
		Chassis []struct {
			OdataID string `json:"@odata.id"`
		} `json:"Chassis"`
		ManagedBy []struct {
			OdataID string `json:"@odata.id"`
		} `json:"ManagedBy"`
	} `json:"Links"`
	Oem struct {
		IntelRackScale struct {
			OdataType                         string        `json:"@odata.type"`
			ProcessorSockets                  int           `json:"ProcessorSockets"`
			MemorySockets                     int           `json:"MemorySockets"`
			UserModeEnabled                   bool          `json:"UserModeEnabled"`
			TrustedExecutionTechnologyEnabled bool          `json:"TrustedExecutionTechnologyEnabled"`
			PciDevices                        []interface{} `json:"PciDevices"`
			PCIeConnectionID                  []interface{} `json:"PCIeConnectionId"`
			Metrics                           struct {
				OdataID string `json:"@odata.id"`
			} `json:"Metrics"`
		} `json:"Intel_RackScale"`
	} `json:"Oem"`
	Actions struct {
		ComputerSystemReset struct {
			Target                          string   `json:"target"`
			ResetTypeRedfishAllowableValues []string `json:"ResetType@Redfish.AllowableValues"`
		} `json:"#ComputerSystem.Reset"`
		Oem struct {
			IntelOemChangeTPMState struct {
				Target                              string   `json:"target"`
				InterfaceTypeRedfishAllowableValues []string `json:"InterfaceType@Redfish.AllowableValues"`
			} `json:"#Intel.Oem.ChangeTPMState"`
		} `json:"Oem"`
	} `json:"Actions"`
}

type ipmiCPU struct {
	OdataContext          string `json:"@odata.context"`
	OdataID               string `json:"@odata.id"`
	OdataType             string `json:"@odata.type"`
	Name                  string `json:"Name"`
	ID                    string `json:"ID"`
	Socket                string `json:"Socket"`
	ProcessorType         string `json:"ProcessorType"`
	ProcessorArchitecture string `json:"ProcessorArchitecture"`
	InstructionSet        string `json:"InstructionSet"`
	Manufacturer          string `json:"Manufacturer"`
	Model                 string `json:"Model"`
	ProcessorID           struct {
		VendorID                string `json:"VendorId"`
		IdentificationRegisters string `json:"IdentificationRegisters"`
		EffectiveFamily         string `json:"EffectiveFamily"`
		EffectiveModel          string `json:"EffectiveModel"`
	} `json:"ProcessorId"`
	MaxSpeedMHz  int `json:"MaxSpeedMHz"`
	TotalCores   int `json:"TotalCores"`
	TotalThreads int `json:"TotalThreads"`
	Status       struct {
		State        string `json:"State"`
		Health       string `json:"Health"`
		HealthRollup string `json:"HealthRollup"`
	} `json:"Status"`
	Oem struct {
		IntelRackScale struct {
			OdataType string `json:"@odata.type"`
			Metrics   struct {
				OdataID string `json:"@odata.id"`
			} `json:"Metrics"`
		} `json:"Intel_RackScale"`
	} `json:"Oem"`
}

type ipmiResetType struct {
	ResetType string `json:"ResetType"`
}

type ipmiNIC struct {
	OdataType    string `json:"@odata.type"`
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	ID           string `json:"ID"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	Status       struct {
		State        string `json:"State"`
		Health       string `json:"Health"`
		HealthRollup string `json:"HealthRollup"`
	} `json:"Status"`
	MACAddress    string   `json:"MACAddress"`
	HostName      string   `json:"HostName"`
	NameServers   []string `json:"NameServers"`
	IPv4Addresses []struct {
		Address       string `json:"Address"`
		SubnetMask    string `json:"SubnetMask"`
		Gateway       string `json:"Gateway"`
		AddressOrigin string `json:"AddressOrigin"`
	} `json:"IPv4Addresses"`
	IPv4StaticAddresses []struct {
		Address    string `json:"Address"`
		SubnetMask string `json:"SubnetMask"`
		Gateway    string `json:"Gateway"`
	} `json:"IPv4StaticAddresses"`
	DHCPv4 struct {
		DHCPEnabled bool `json:"DHCPEnabled"`
	} `json:"DHCPv4"`
	MTUSize                int `json:"MTUSize"`
	MaxIPv6StaticAddresses int `json:"MaxIPv6StaticAddresses"`
	IPv6Addresses          []struct {
		Address       string `json:"Address"`
		PrefixLength  int    `json:"PrefixLength"`
		AddressOrigin string `json:"AddressOrigin"`
	} `json:"IPv6Addresses"`
	IPv6StaticAddresses []struct {
		Address      string `json:"Address"`
		PrefixLength int    `json:"PrefixLength"`
	} `json:"IPv6StaticAddresses"`
	IPv6DefaultGateway        string `json:"IPv6DefaultGateway"`
	IPv6StaticDefaultGateways []struct {
		Address string `json:"Address"`
	} `json:"IPv6StaticDefaultGateways"`
	DHCPv6 struct {
		OperatingMode string `json:"OperatingMode"`
	} `json:"DHCPv6"`
	LinkStatus       string      `json:"LinkStatus"`
	InterfaceEnabled bool        `json:"InterfaceEnabled"`
	SpeedMbps        int         `json:"SpeedMbps"`
	FQDN             interface{} `json:"FQDN"`
	VLAN             struct {
		VLANEnable bool `json:"VLANEnable"`
		VLANID     int  `json:"VLANId"`
	} `json:"VLAN"`
}
