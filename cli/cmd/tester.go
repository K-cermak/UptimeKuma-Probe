package cmd

import (
	"UptimeKumaProbeCLI/db"
	"UptimeKumaProbeCLI/helpers"
	"UptimeKumaProbeCLI/utils"
	"os/exec"
	"runtime"
)

func PingTest(address string, timeout string) {
	count, correct := helpers.StrToInt(timeout)
	if !correct {
		helpers.PrintError(true, "Invalid timeout value")
	}

	helpers.PrintInfo("Pinging " + address + " with timeout " + timeout + " ms")
	ret := utils.PingAddress(address, count, true)
	if ret {
		helpers.PrintSuccess("Ping successful")
	} else {
		helpers.PrintWarning("Ping failed")
	}
}

func HttpTest(address string, timeout string) {
	count, correct := helpers.StrToInt(timeout)
	if !correct {
		helpers.PrintError(true, "Invalid timeout value")
	}

	ignoreSslStr := db.GetValue("ignore_ssl_errors")
	ignoreSsl := false
	if ignoreSslStr == "true" {
		ignoreSsl = true
	}

	helpers.PrintInfo("Performing HTTP request to " + address + " with timeout " + timeout + " ms")
	ret := utils.CheckHTTP(address, count, "", "", ignoreSsl, true)
	if ret {
		helpers.PrintSuccess("HTTP request successful")
	} else {
		helpers.PrintWarning("HTTP request failed")
	}
}

func ApiTest(testType string) {
	if testType != "service" && testType != "http" {
		helpers.PrintError(true, "Invalid type, expected <service> or <http>")
	}

	if testType == "http" {
		apiPort := db.GetValue("api_port")
		HttpTest("http://127.0.0.1:"+apiPort, "5000")

	} else if testType == "service" {
		if runtime.GOOS != "linux" {
			helpers.PrintError(true, "This service testing is only available on Linux")
		}

		cmd := exec.Command("systemctl", "is-active", "kprobe")
		output, err := cmd.CombinedOutput()
		if err != nil {
			helpers.PrintError(true, "Failed to check service status ("+err.Error()+")")
		}

		if string(output) == "active\n" {
			helpers.PrintSuccess("Service is active")
		} else {
			helpers.PrintWarning("Service is not active")
		}
	}
}

func ApiRestart() {
	if runtime.GOOS != "linux" {
		helpers.PrintError(true, "This service testing is only available on Linux")
	}

	cmd := exec.Command("systemctl", "restart", "kprobe")
	err := cmd.Run()
	if err != nil {
		helpers.PrintError(true, "Failed to restart service ("+err.Error()+"). Tip: Use 'sudo' to run the command")
	}

	helpers.PrintSuccess("Service restarted")
}