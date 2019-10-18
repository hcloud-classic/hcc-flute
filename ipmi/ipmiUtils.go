package ipmi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"hcc/flute/config"
	"hcc/flute/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetSerialNo : Get serial number from IPMI node
func GetSerialNo(bmcIP string) (string, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/", nil)
		if err != nil {
			return "", err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("GetSerialNo(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNode ipmiNode
			err = json.Unmarshal([]byte(str), &ipmiNode)
			if err != nil {
				return "", err
			}

			serialNo := ipmiNode.Members[0].OdataID[len("/redfish/v1/Systems/"):]

			return serialNo, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

// GetUUID : Get UUID from IPMI node
func GetUUID(bmcIP string, serialNo string) (string, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo, nil)
		if err != nil {
			return "", err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("GetUUID(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNodeInfo ipmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				return "", err
			}

			uuid := ipmiNodeInfo.UUID

			return uuid, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

// GetPowerState : Get power status from IPMI node
func GetPowerState(bmcIP string, serialNo string) (string, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo, nil)
		if err != nil {
			return "", err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("GetPowerState(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNodeInfo ipmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				return "", err
			}

			powerState := ipmiNodeInfo.PowerState

			return powerState, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

// GetProcessors : Get count of CPU processors from IPMI node
func GetProcessors(bmcIP string, serialNo string) (int, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors", nil)
		if err != nil {
			return 0, err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("GetProcessors(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var processors ipmiProcessors
			err = json.Unmarshal([]byte(str), &processors)
			if err != nil {
				return 0, err
			}

			count := len(processors.Members)

			return count, nil
		}

		return 0, err
	}

	return 0, errors.New("http response returned error code")
}

// GetProcessorsCores : Get count of cores for selected processor from IPMI node
func GetProcessorsCores(bmcIP string, serialNo string, processors int) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	coreSum := 0

	for i := 1; i <= processors; i++ {
		var resp *http.Response

		for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
			req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors/CPU"+strconv.Itoa(i), nil)
			if err != nil {
				return 0, err
			}
			req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
			resp, err = client.Do(req)
			if err != nil {
				logger.Logger.Println(err)
				logger.Logger.Println("GetProcessorsCores(): CPU" + strconv.Itoa(i) + ": Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			} else {
				break
			}
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				str := string(respBody)

				var cpu ipmiCPU
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					return 0, err
				}

				totalCores := cpu.TotalCores

				coreSum += totalCores
			}
		} else {
			return 0, errors.New("http response returned error code")
		}
		_ = resp.Body.Close()
	}

	return coreSum, nil
}

// GetProcessorsThreads : Get count of threads for selected processor from IPMI node
func GetProcessorsThreads(bmcIP string, serialNo string, processors int) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
	threadSum := 0

	for i := 1; i <= processors; i++ {
		var resp *http.Response

		for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
			req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors/CPU"+strconv.Itoa(i), nil)
			if err != nil {
				return 0, err
			}
			req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
			resp, err = client.Do(req)
			if err != nil {
				logger.Logger.Println(err)
				logger.Logger.Println("GetProcessorsThreads(): CPU" + strconv.Itoa(i) + ": Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
			} else {
				break
			}
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				str := string(respBody)

				var cpu ipmiCPU
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					return 0, err
				}

				totalThreads := cpu.TotalThreads

				threadSum += totalThreads
			}
		} else {
			return 0, errors.New("http response returned error code")
		}
		_ = resp.Body.Close()
	}

	return threadSum, nil
}

// GetProcessorModel : Get model of first processor from IPMI node
func GetProcessorModel(bmcIP string, serialNo string) (string, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo+"/Processors/CPU1", nil)
		if err != nil {
			return "", err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("GetProcessorModel(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var cpu ipmiCPU
			err = json.Unmarshal([]byte(str), &cpu)
			if err != nil {
				return "", err
			}

			return cpu.ProcessorID.VendorID, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

// GetTotalSystemMemory : Get total system memory from IPMI node
func GetTotalSystemMemory(bmcIP string, serialNo string) (int, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Systems/"+serialNo, nil)
		if err != nil {
			return 0, err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("GetTotalSystemMemory(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNodeInfo ipmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				return 0, err
			}

			memoryGiB := ipmiNodeInfo.MemorySummary.TotalSystemMemoryGiB

			return memoryGiB, nil
		}

		return 0, err
	}

	return 0, errors.New("http response returned error code")

}

// ChangePowerState : Change power status for selected IPMI node
func ChangePowerState(bmcIP string, serialNo string, state string) (string, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
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

		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("ChangePowerState(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var processors ipmiProcessors
			err = json.Unmarshal([]byte(str), &processors)
			if err != nil {
				return "", err
			}

			return str, nil
		}

		return "", err

	}

	return "", errors.New("http response returned error code")
}

// GetNICMac : Get MAC address of selected NIC from IPMI node
func GetNICMac(bmcIP string, nicNO int, isBMC bool) (string, error) {
	var resp *http.Response

	for i := 0; i < int(config.Ipmi.RequestRetry); i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		client := &http.Client{Timeout: time.Duration(config.Ipmi.RequestTimeoutMs) * time.Millisecond}
		req, err := http.NewRequest("GET", "https://"+bmcIP+"/redfish/v1/Managers/BMC/EthernetInterfaces/"+strconv.Itoa(nicNO), nil)
		if err != nil {
			return "", err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err = client.Do(req)
		if err != nil {
			logger.Logger.Println(err)
			logger.Logger.Println("GetNICMac(): Retrying for " + bmcIP + " " + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Ipmi.RequestRetry)))
		} else {
			break
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var nic ipmiNIC
			err = json.Unmarshal([]byte(str), &nic)
			if err != nil {
				return "", err
			}

			var macAddress = ""
			if isBMC {
				macAddress = nic.MACAddress
			} else {
				macParts := strings.Split(nic.MACAddress, "-")
				if len(macParts) != 6 {
					return "", errors.New("invalid mac address")
				}

				lastPart := lastMacOffset(macParts[len(macParts)-1])
				for i, part := range macParts {
					if i == len(macParts)-1 {
						break
					}
					macAddress += part + "-"
				}
				macAddress += lastPart
			}

			return macAddress, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}
