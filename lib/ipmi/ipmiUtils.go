package ipmi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetSerialNo : Get serial number from IPMI node
func GetSerialNo(bmcIP string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("GetSerialNo(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("GetSerialNo(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			str := string(respBody)

			var ipmiNodes ipmiMembers
			err = json.Unmarshal([]byte(str), &ipmiNodes)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			serialNo := ipmiNodes.Members[0].OdataID[len("/redfish/v1/Systems/"):]

			_ = resp.Body.Close()
			return serialNo, nil
		}
	}

	return "", errors.New("GetSerialNo(): retry count exceeded for " + bmcIP)
}

// GetUUID : Get UUID from IPMI node
func GetUUID(bmcIP string, serialNo string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("GetUUID(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("GetUUID(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			str := string(respBody)

			var ipmiNodeInfo ipmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			uuid := ipmiNodeInfo.UUID

			_ = resp.Body.Close()
			return uuid, nil
		}
	}

	return "", errors.New("GetUUID(): retry count exceeded for " + bmcIP)
}

// GetPowerState : Get power status from IPMI node
func GetPowerState(bmcIP string, serialNo string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("GetPowerState(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("GetPowerState(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {

			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			str := string(respBody)

			var ipmiNodeInfo ipmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			powerState := ipmiNodeInfo.PowerState

			_ = resp.Body.Close()
			return powerState, nil
		}

	}

	return "", errors.New("GetPowerState(): retry count exceeded for " + bmcIP)
}

// GetTotalSystemMemory : Get total system memory from IPMI node
func GetTotalSystemMemory(bmcIP string, serialNo string) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo, nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("GetTotalSystemMemory(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("GetTotalSystemMemory(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return 0, err
			}

			str := string(respBody)

			var ipmiNodeInfo ipmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				_ = resp.Body.Close()
				return 0, err
			}

			memoryGiB := ipmiNodeInfo.MemorySummary.TotalSystemMemoryGiB

			_ = resp.Body.Close()
			return memoryGiB, nil
		}
	}

	return 0, errors.New("GetTotalSystemMemory(): retry count exceeded for " + bmcIP)
}

// ChangePowerState : Change power status for selected IPMI node
func ChangePowerState(bmcIP string, serialNo string, state string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}

	resetType := ipmiResetType{state}
	jsonBytes, err := json.Marshal(resetType)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Actions/ComputerSystem.Reset", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("ChangePowerState(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("ChangePowerState(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				// Check response
				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					_ = resp.Body.Close()
					return "", err
				}

				str := string(respBody)

				_ = resp.Body.Close()
				return str, nil
			}

			_ = resp.Body.Close()
			return "", err
		}
	}

	return "", errors.New("ChangePowerState(): retry count exceeded for " + bmcIP)
}

// GetNICMac : Get MAC address of selected NIC from IPMI node
func GetNICMac(bmcIP string, nicNO int, isBMC bool) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Managers/BMC/EthernetInterfaces/"+strconv.Itoa(nicNO), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("GetNICMac(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("GetNICMac(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			str := string(respBody)

			var nic ipmiNIC
			err = json.Unmarshal([]byte(str), &nic)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			var macAddress = ""
			if isBMC {
				macAddress = nic.MACAddress
			} else {
				macParts := strings.Split(nic.MACAddress, "-")
				if len(macParts) != 6 {
					_ = resp.Body.Close()
					return "", errors.New("GetNICMac(): Invalid mac address")
				}

				lastPart := lastMacOffset(macParts[len(macParts)-1], nicNO+1)
				for i, part := range macParts {
					if i == len(macParts)-1 {
						break
					}
					macAddress += part + "-"
				}
				macAddress += lastPart
			}

			_ = resp.Body.Close()
			return macAddress, nil
		}
	}

	return "", errors.New("GetNICMac(): retry count exceeded for " + bmcIP)
}

func getMemberCounts(bmcIP string, serialNo string, member string) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/"+member, nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("Get" + member + "(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("Get" + member + "(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {

			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return 0, err
			}

			str := string(respBody)

			var members ipmiMembers
			err = json.Unmarshal([]byte(str), &members)
			if err != nil {
				_ = resp.Body.Close()
				return 0, err
			}

			count := len(members.Members)

			_ = resp.Body.Close()
			return count, nil
		}

	}

	return 0, errors.New("Get" + member + "(): retry count exceeded for " + bmcIP)
}

func getNumCPU(bmcIP string, serialNo string) (int, error) {
	return getMemberCounts(bmcIP, serialNo, "Processors")
}

// GetProcessorsCores : Get count of cores for selected processor from IPMI node
func GetProcessorsCores(bmcIP string, serialNo string, processors int) (int, error) {
	coreSum := 0

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}

	for i := 1; i <= processors; i++ {
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors/CPU"+strconv.Itoa(i), nil)
		if err != nil {
			return 0, err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

		var j = 0
		for ; j < int(config.Ipmi.RequestRetry); j++ {
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
				if err != nil {
					logger.Logger.Println(err)
				} else {
					_ = resp.Body.Close()
					logger.Logger.Println("GetProcessorsCores(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
				}
				logger.Logger.Println("GetProcessorsCores(): Retrying for " + bmcIP + " " + strconv.Itoa(j+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
				continue
			} else {
				// Check response
				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					_ = resp.Body.Close()
					return 0, err
				}

				str := string(respBody)

				var cpu ipmiProcessor
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					_ = resp.Body.Close()
					return 0, err
				}

				totalCores := cpu.TotalCores

				coreSum += totalCores

				_ = resp.Body.Close()
				break
			}
		}
		if j == int(config.Ipmi.RequestRetry) {
			return 0, errors.New("GetProcessorsCores(): retry count exceeded for " + bmcIP)
		}
	}

	return coreSum, nil
}

func getCPUDetails(bmcIP string, serialNo string, ID string) (*cpu, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors/"+ID, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("getCPUDetails(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("getCPUDetails(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return nil, err
			}

			str := string(respBody)

			var ipmiProcessor ipmiProcessor
			err = json.Unmarshal([]byte(str), &ipmiProcessor)
			if err != nil {
				_ = resp.Body.Close()
				return nil, err
			}
			_ = resp.Body.Close()

			cpu := cpu{
				ID:          ipmiProcessor.ID,
				Socket:      ipmiProcessor.Socket,
				Manufacture: ipmiProcessor.Manufacturer,
				Model:       ipmiProcessor.ProcessorID.VendorID,
				MaxSpeedMHz: ipmiProcessor.MaxSpeedMHz,
				Cores:       ipmiProcessor.TotalCores,
				Threads:     ipmiProcessor.TotalThreads,
				Status: status{
					State:        ipmiProcessor.Status.State,
					Health:       ipmiProcessor.Status.Health,
					HealthRollup: ipmiProcessor.Status.HealthRollup,
				},
			}

			return &cpu, nil
		}
	}

	return nil, errors.New("getCPUDetails(): retry count exceeded for " + bmcIP)
}

func getCPUsDetail(bmcIP string, serialNo string) (*[]cpu, error) {
	cpuNum, err := getNumCPU(bmcIP, serialNo)
	if err != nil {
		return nil, err
	}

	var cpus []cpu
	for i := 1; i <= cpuNum; i++ {
		cpu, err := getCPUDetails(bmcIP, serialNo, "CPU"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		cpus = append(cpus, *cpu)
	}

	return &cpus, nil
}

func getNumMemory(bmcIP string, serialNo string) (int, error) {
	return getMemberCounts(bmcIP, serialNo, "Memory")
}

func getMemoryDetails(bmcIP string, serialNo string, ID string) (*memory, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Memory/"+ID, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			if err != nil {
				logger.Logger.Println(err)
			} else {
				_ = resp.Body.Close()
				logger.Logger.Println("getMemoryDetails(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("getMemoryDetails(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return nil, err
			}

			str := string(respBody)

			var ipmiMemory ipmiMemory
			err = json.Unmarshal([]byte(str), &ipmiMemory)
			if err != nil {
				_ = resp.Body.Close()
				return nil, err
			}
			_ = resp.Body.Close()

			memory := memory{
				ID:            ipmiMemory.ID,
				CapacityMB:    ipmiMemory.Capacitymib,
				Manufacture:   ipmiMemory.Manufacturer,
				SerialNumber:  ipmiMemory.Serialnumber,
				PartNumber:    ipmiMemory.Partnumber,
				DeviceLocator: ipmiMemory.Devicelocator,
				SpeedMhz:      ipmiMemory.Operatingspeedmhz,
				Status: status{
					State:        ipmiMemory.Status.State,
					Health:       ipmiMemory.Status.Health,
					HealthRollup: ipmiMemory.Status.Healthrollup,
				},
			}

			return &memory, nil
		}
	}

	return nil, errors.New("getMemoryDetails(): retry count exceeded for " + bmcIP)
}

func getMemoriesDetail(bmcIP string, serialNo string) (*[]memory, error) {
	memoryNum, err := getNumMemory(bmcIP, serialNo)
	if err != nil {
		return nil, err
	}

	var memories []memory
	for i := 1; i <= memoryNum; i++ {
		memory, err := getMemoryDetails(bmcIP, serialNo, "Memory"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		memories = append(memories, *memory)
	}

	return &memories, nil
}

// GetNodeDetailData : Get JSON data of node's detail info
func GetNodeDetailData(bmcIP string, serialNo string) (string, error) {
	memoriesDetail, err := getMemoriesDetail(bmcIP, serialNo)
	if err != nil {
		return "", err
	}

	getCPUsDetail, err := getCPUsDetail(bmcIP, serialNo)
	if err != nil {
		return "", err
	}

	nodeDetailData := nodeDetailData{
		Memories: *memoriesDetail,
		CPUs:     *getCPUsDetail,
	}

	result, err := json.Marshal(nodeDetailData)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
