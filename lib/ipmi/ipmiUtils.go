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

			var ipmiNode ipmiNode
			err = json.Unmarshal([]byte(str), &ipmiNode)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			serialNo := ipmiNode.Members[0].OdataID[len("/redfish/v1/Systems/"):]

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

// GetProcessors : Get count of CPU processors from IPMI node
func GetProcessors(bmcIP string, serialNo string) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors", nil)
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
				logger.Logger.Println("GetProcessors(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("GetProcessors(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {

			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return 0, err
			}

			str := string(respBody)

			var processors ipmiProcessors
			err = json.Unmarshal([]byte(str), &processors)
			if err != nil {
				_ = resp.Body.Close()
				return 0, err
			}

			count := len(processors.Members)

			_ = resp.Body.Close()
			return count, nil
		}

	}

	return 0, errors.New("GetProcessors(): retry count exceeded for " + bmcIP)
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

				var cpu ipmiCPU
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

// GetProcessorsThreads : Get count of threads for selected processor from IPMI node
func GetProcessorsThreads(bmcIP string, serialNo string, processors int) (int, error) {
	threadSum := 0

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
					logger.Logger.Println("GetProcessorsThreads(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
				}
				logger.Logger.Println("GetProcessorsThreads(): Retrying for " + bmcIP + " " + strconv.Itoa(j+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
				continue
			} else {
				// Check response
				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					_ = resp.Body.Close()
					return 0, err
				}

				str := string(respBody)

				var cpu ipmiCPU
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					_ = resp.Body.Close()
					return 0, err
				}

				totalThreads := cpu.TotalThreads

				threadSum += totalThreads

				_ = resp.Body.Close()
				break
			}
		}
		if j == int(config.Ipmi.RequestRetry) {
			return 0, errors.New("GetProcessorsThreads(): retry count exceeded for " + bmcIP)
		}
	}

	return threadSum, nil
}

// GetProcessorModel : Get model of first processor from IPMI node
func GetProcessorModel(bmcIP string, serialNo string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors/CPU1", nil)
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
				logger.Logger.Println("GetProcessorModel(): http response returned error code " + strconv.Itoa(resp.StatusCode) + " for " + bmcIP)
			}
			logger.Logger.Println("GetProcessorModel(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			continue
		} else {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			str := string(respBody)

			var cpu ipmiCPU
			err = json.Unmarshal([]byte(str), &cpu)
			if err != nil {
				_ = resp.Body.Close()
				return "", err
			}

			_ = resp.Body.Close()
			return cpu.ProcessorID.VendorID, nil
		}
	}

	return "", errors.New("GetProcessorModel(): retry count exceeded for " + bmcIP)
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
