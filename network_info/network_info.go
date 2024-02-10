package network_info

import (
	"fmt"
	"gablogger"
	"os/exec"

	"github.com/vishvananda/netlink"
)

var Logger = gablogger.Logger()

func GetPriorityRoute() (string, error) {
	routes, err := netlink.RouteList(nil, int(netlink.FAMILY_ALL))
	if err != nil {
		return "", fmt.Errorf("cannot get list of routes")
	}

	var selectedRoute netlink.Route
	minMetric := -1
	for _, route := range routes {
		if route.Dst == nil { // filter default route
			if minMetric == -1 || route.Priority < minMetric {
				minMetric = route.Priority
				selectedRoute = route
			}
		}
	}

	iface, err := netlink.LinkByIndex(selectedRoute.LinkIndex)
	if err != nil {
		return "", fmt.Errorf("cannot retrieve interface name from priority route")
	}

	return iface.Attrs().Name, nil
}

func getIpFromInterface(interfaceName string) (string, error) {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		errorMessage := fmt.Sprintf("Interface '%s' not found: %v", interfaceName, err)
		return "not found", fmt.Errorf(errorMessage)
	}

	addrs, err := netlink.AddrList(iface, netlink.FAMILY_ALL)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting IP addresses for interface '%s': %v", interfaceName, err)
		return "without ip", fmt.Errorf(errorMessage)
	}

	if len(addrs) > 0 {
		ipAddress := addrs[0].IP.String()
		return ipAddress, nil
	}

	errorMessage := fmt.Sprintf("No IP address found for interface '%s'", interfaceName)
	return "no IP found", fmt.Errorf(errorMessage)
}

func GetWiredInterfaceStatus() (string, error) {
	ip, err := getIpFromInterface("eth0.2")
	if err != nil {
		Logger.Error(err)
		return ip, err
	}
	if checkConnectivity("eth0.2") {
		return "available", nil
	} else {
		return "unavailable", nil
	}
}

func GetModemInterfaceStatus() (string, error) {
	ip, err := getIpFromInterface("wwan0")
	if err != nil {
		Logger.Error(err)
		return ip, err
	}
	if checkConnectivity("wwan0") {
		return "available", nil
	} else {
		return "unavailable", nil
	}
}

func checkConnectivity(networkInterface string) bool {
	_, err := exec.Command("ping", "-I", networkInterface, "-w", "5", "-c", "1", "google.com").CombinedOutput()
	return err == nil
}
