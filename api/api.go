package api

import (
	"common"
	"device_info"
	"encoding/json"
	"fmt"
	"gablogger"
	"net/http"
	"network_info"
	"peripherals"
	"socketxp"
	"strconv"
)

var Logger = gablogger.Logger()

func Start() {
	if common.API_PORT == "" {
		Logger.Error("API port is undefined. Please set a valid port.")
		return
	}

	if !isValidPort(common.API_PORT) {
		Logger.Error("API port is invalid. Please specify a valid port number (1-65535).")
		return
	}

	mux := http.NewServeMux()
	routes := map[string]func() (string, error){
		"/diagnosis/modem/signal-strength":  peripherals.GetModemSignalStrength,
		"/diagnosis/modem/sim-card-type":    peripherals.GetSIMCardType,
		"/diagnosis/modem/sim-card-iccid":   peripherals.GetSIMCardICCID,
		"/diagnosis/modem/sim-card-carrier": peripherals.GetSIMCardCarrier,
		"/diagnosis/power/source":           peripherals.GetPowerSource,
		"/diagnosis/power/bms":              peripherals.GetHasBMS,
		"/diagnosis/power/battery-level":    peripherals.GetBatteryLevel,
		"/diagnosis/stm32/firmware-version": peripherals.GetFirmwareVersion,
		"/diagnosis/stm32/temperature":      peripherals.GetSTM32Temperature,
		"/diagnosis/fabrication/pcb-batch":  peripherals.GetPCBBatch,
		"/diagnosis/fabrication/pcb-review": peripherals.GetPCBReview,
		"/diagnosis/socketxp/status":        socketxp.IsConnected,
		"/diagnosis/device/serial-number":   device_info.GetDeviceId,
		"/diagnosis/device/os-version":      device_info.GetOSVersion,
		"/diagnosis/network/priority-route": network_info.GetPriorityRoute,
		"/diagnosis/network/modem":          network_info.GetModemInterfaceStatus,
		"/diagnosis/network/wired":          network_info.GetWiredInterfaceStatus,
	}

	// Register the routes with the router
	for path, handler := range routes {
		pathCopy := path // Copy to avoid variable capture in loop
		handlerCopy := handler
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			handleGenericRequest(w, r, handlerCopy, pathCopy)
		})
	}

	Logger.Debug("Starting API Server in port: ", common.API_PORT)
	err := http.ListenAndServe(":"+common.API_PORT, mux)
	if err != nil {
		Logger.Error("Error starting the server:", err)
	}
	Logger.Debug("Server iniciado")
}

func handleGenericRequest(w http.ResponseWriter, r *http.Request, handler func() (string, error), path string) {
	Logger.Debugf("Received request at %s", path)

	response, err := handler()
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		Logger.Errorf("Error processing request at %s: %v", path, err)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dataResponse := map[string]interface{}{"data": response}

	if err := json.NewEncoder(w).Encode(dataResponse); err != nil {
		Logger.Errorf("Error encoding response for request at %s: %v", path, err)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Logger.Debugf("Response sent for request at %s: %v", path, dataResponse)
}

func writeJSONError(w http.ResponseWriter, reason string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"reason": "%s"}`, reason)
}

func isValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	return err == nil && port >= 1 && port <= 65535
}
