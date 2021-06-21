package ipmi

import (
	"hcloud-flute/logger"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetSerialNo(ipmiIp string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Systems/", nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNode IpmiNode
			err = json.Unmarshal([]byte(str), &ipmiNode)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			serialNo := ipmiNode.Members[0].OdataID[len("/redfish/v1/Systems/"):]

			return serialNo, nil
		} else {
			return "", err
		}
	} else {
		return "", errors.New("HTTP response returned error!")
	}
}

func GetUuid(ipmiIp string, serialNo string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Systems/" + serialNo, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNodeInfo IpmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			uuid := ipmiNodeInfo.UUID

			return uuid, nil
		} else {
			return "", err
		}
	} else {
		return "", errors.New("HTTP response returned error!")
	}
}

func GetPowerState(ipmiIp string, serialNo string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Systems/" + serialNo, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNodeInfo IpmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			powerState := ipmiNodeInfo.PowerState

			return powerState, nil
		} else {
			return "", err
		}
	} else {
		return "", errors.New("HTTP response returned error!")
	}
}

func GetProcessors(ipmiIp string, serialNo string) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Systems/" + serialNo + "/Processors", nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var processors Processors
			err = json.Unmarshal([]byte(str), &processors)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			count := len(processors.Members)

			return count, nil
		} else {
			return 0, err
		}
	} else {
		return 0, errors.New("HTTP response returned error!")
	}
}

func GetProcessorsCores(ipmiIp string, serialNo string, processors int) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	coreSum := 0

	for i := 1; i <= processors; i++ {
		req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Systems/" + serialNo + "/Processors/CPU" + strconv.Itoa(i), nil)
		req.SetBasicAuth(username, password)
		resp, err := client.Do(req)

		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				str := string(respBody)

				var cpu Cpu
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					logger.Logger.Fatal(err)
				}

				totalCores := cpu.TotalCores

				coreSum += totalCores
			} else {
				return 0, err
			}
		} else {
			return 0, errors.New("HTTP response returned error!")
		}
	}

	return coreSum, nil
}

func GetProcessorsThreads(ipmiIp string, serialNo string, processors int) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	threadSum := 0

	for i := 1; i <= processors; i++ {
		req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Systems/" + serialNo + "/Processors/CPU" + strconv.Itoa(i), nil)
		req.SetBasicAuth(username, password)
		resp, err := client.Do(req)

		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			// Check response
			respBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				str := string(respBody)

				var cpu Cpu
				err = json.Unmarshal([]byte(str), &cpu)
				if err != nil {
					logger.Logger.Fatal(err)
				}

				totalThreads := cpu.TotalThreads

				threadSum += totalThreads
			} else {
				return 0, err
			}
		} else {
			return 0, errors.New("HTTP response returned error!")
		}
	}

	return threadSum, nil
}

func GetTotalSystemMemory(ipmiIp string, serialNo string) (int, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Systems/" + serialNo, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var ipmiNodeInfo IpmiNodeInfo
			err = json.Unmarshal([]byte(str), &ipmiNodeInfo)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			memoryGiB := ipmiNodeInfo.MemorySummary.TotalSystemMemoryGiB

			return memoryGiB, nil
		} else {
			return 0, err
		}
	} else {
		return 0, errors.New("HTTP response returned error!")
	}
}

func GetInfo(ipmiIp string) {
	serialNo, err := GetSerialNo(ipmiIp)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("SerialNo: " + serialNo)

	uuid, err := GetUuid(ipmiIp, serialNo)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("UUID: " + uuid)

	powerState, err := GetPowerState(ipmiIp, serialNo)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("Status: " + powerState)

	processors, err := GetProcessors(ipmiIp, serialNo)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("Processors: " + strconv.Itoa(processors))

	cores, err := GetProcessorsCores(ipmiIp, serialNo, processors)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("Cores: " + strconv.Itoa(cores))

	threads, err := GetProcessorsThreads(ipmiIp, serialNo, processors)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("Threads: " + strconv.Itoa(threads))

	memory, err := GetTotalSystemMemory(ipmiIp, serialNo)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("Memory: " + strconv.Itoa(memory) + "GiB")

	mac, err := GetBMCNICMac(ipmiIp)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.Println("MAC Address: " + mac)
}

func ChangePowerState(ipmiIp string, serialNo string, state string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}

	resetType := ResetType{state}
	jsonBytes, err := json.Marshal(resetType)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://" + ipmiIp + "/redfish/v1/Systems/" + serialNo + "/Actions/ComputerSystem.Reset", bytes.NewBuffer(jsonBytes))
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var processors Processors
			err = json.Unmarshal([]byte(str), &processors)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			return str, nil
		} else {
			return "", err
		}
	} else {
		return "", errors.New("HTTP response returned error!")
	}
}

func GetBMCNICMac(ipmiIp string) (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://" + ipmiIp + "/redfish/v1/Managers/BMC/EthernetInterfaces/3", nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		// Check response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)

			var bmcNIC BMCNIC
			err = json.Unmarshal([]byte(str), &bmcNIC)
			if err != nil {
				logger.Logger.Fatal(err)
			}

			macAddress := bmcNIC.MACAddress

			return macAddress, nil
		} else {
			return "", err
		}
	} else {
		return "", errors.New("HTTP response returned error!")
	}
}