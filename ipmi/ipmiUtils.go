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
	"time"
)

// GetSerialNo : Get serial number from IPMI node
func GetSerialNo(ipmiIP string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Systems/", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
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
				logger.Logger.Fatal(err)
			}

			serialNo := ipmiNode.Members[0].OdataID[len("/redfish/v1/Systems/"):]

			return serialNo, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

// GetUUID : Get UUID from IPMI node
func GetUUID(ipmiIP string, serialNo string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Systems/"+serialNo, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
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
				logger.Logger.Fatal(err)
			}

			uuid := ipmiNodeInfo.UUID

			return uuid, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

// GetPowerState : Get power status from IPMI node
func GetPowerState(ipmiIP string, serialNo string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Systems/"+serialNo, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
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
				logger.Logger.Fatal(err)
			}

			powerState := ipmiNodeInfo.PowerState

			return powerState, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}

// GetProcessors : Get count of CPU processors from IPMI node
func GetProcessors(ipmiIP string, serialNo string) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Systems/"+serialNo+"/Processors", nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
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
				logger.Logger.Fatal(err)
			}

			count := len(processors.Members)

			return count, nil
		}

		return 0, err
	}

	return 0, errors.New("http response returned error code")
}

// GetProcessorsCores : Get count of cores for selected processor from IPMI node
func GetProcessorsCores(ipmiIP string, serialNo string, processors int) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	coreSum := 0

	for i := 1; i <= processors; i++ {
		req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Systems/"+serialNo+"/Processors/CPU"+strconv.Itoa(i), nil)
		if err != nil {
			return 0, err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err := client.Do(req)

		if err != nil {
			return 0, err
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				str := string(respBody)

				var cpu ipmiCPU
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					logger.Logger.Fatal(err)
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
func GetProcessorsThreads(ipmiIP string, serialNo string, processors int) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	threadSum := 0

	for i := 1; i <= processors; i++ {
		req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Systems/"+serialNo+"/Processors/CPU"+strconv.Itoa(i), nil)
		if err != nil {
			return 0, err
		}
		req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
		resp, err := client.Do(req)

		if err != nil {
			return 0, err
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				str := string(respBody)

				var cpu ipmiCPU
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					logger.Logger.Fatal(err)
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

// GetTotalSystemMemory : Get total system memory from IPMI node
func GetTotalSystemMemory(ipmiIP string, serialNo string) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Systems/"+serialNo, nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
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
				logger.Logger.Fatal(err)
			}

			memoryGiB := ipmiNodeInfo.MemorySummary.TotalSystemMemoryGiB

			return memoryGiB, nil
		}

		return 0, err
	}

	return 0, errors.New("http response returned error code")

}

// ChangePowerState : Change power status for selected IPMI node
func ChangePowerState(ipmiIP string, serialNo string, state string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}

	resetType := ipmiResetType{state}
	jsonBytes, err := json.Marshal(resetType)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://"+ipmiIP+"/redfish/v1/Systems/"+serialNo+"/Actions/ComputerSystem.Reset", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)

	resp, err := client.Do(req)

	if err != nil {
		return "", err
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
				logger.Logger.Fatal(err)
			}

			return str, nil
		}

		return "", err

	}

	return "", errors.New("http response returned error code")
}

// GetBMCNICMac : Get MAC address of BMC interface from IPMI node
func GetBMCNICMac(ipmiIP string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("GET", "https://"+ipmiIP+"/redfish/v1/Managers/BMC/EthernetInterfaces/3", nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.Ipmi.Username, config.Ipmi.Password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var bmcNIC ipmiBMCNIC
			err = json.Unmarshal([]byte(str), &bmcNIC)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			macAddress := bmcNIC.MACAddress

			return macAddress, nil
		}

		return "", err
	}

	return "", errors.New("http response returned error code")
}
