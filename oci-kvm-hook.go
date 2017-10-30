// +build linux

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"strings"
)

// {"version":"","id":"96c1870d5c21c324db0a5b350f5a8f5571cf514205b6d3647b893b6580a05018","pid":27146,"root":"/opt/docker/devicemapper/mnt/a45fd0cb88f52620f94215e8e19b806f787f0c649be8e5c13b737dc2d8278daf/rootfs"}
type State struct {
	Version    string `json:"version"`
	ID         string `json:"id"`
	Pid        int    `json:"pid"`
	Root       string `json:"root"`
	BundlePath string `json:"bundlePath"`
}

type Process struct {
	Env []string `json:"env"`
}

func stringInArray(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

func getProcessDevicesCgroupPath(pid int) (string, error) {

	pid_cgroup_path := fmt.Sprintf("/proc/%d/cgroup", pid)
	pid_cgroup, err := os.Open(pid_cgroup_path)
	if err != nil {
		return "", err
	}
	defer pid_cgroup.Close()

	scanner := bufio.NewScanner(pid_cgroup)
	for scanner.Scan() {
		entry := strings.SplitN(scanner.Text(), ":", 3)
		if entry[0] == "0" {
			// let's just ignore cgroups v2 for now -- e.g. systemd still uses
			// v1 for the devices subsystem
			continue
		}
		controllers := strings.Split(entry[1], ",")
		if stringInArray("devices", controllers) {
			// note systemd leaves a symlink even in the case it's comounted
			return fmt.Sprintf("/sys/fs/cgroup/devices/%s", entry[2]), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func allowKvm(state State) {

	cgroup_path, err := getProcessDevicesCgroupPath(state.Pid)
	if err != nil {
		log.Printf("Failed to get process %d devices cgroup path: %v", state.Pid, err.Error())
		return
	}

	if cgroup_path == "" {
		log.Printf("Process does not belong to any devices cgroup, yay!")
	} else {
		allow_path := fmt.Sprintf("%s/devices.allow", cgroup_path)
		allow, err := os.OpenFile(allow_path, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("Failed to open file for writing: %s: %v", allow_path, err.Error())
			return
		}
		defer allow.Close()

		_, err = allow.WriteString("c 10:232 rwm")
		if err != nil {
			log.Printf("Failed to write to group file: %s: %v", allow_path, err.Error())
			return
		}

		log.Printf("Added kvm whitelist into %s", allow_path)
	}

	// Get info about /dev/kvm
	info, err := os.Stat("/dev/kvm")
	if err != nil {
		log.Printf("Skipping /dev/kvm creation in container because not present on host: %v", err.Error())
		return
	}

	// A mode like 0666 or 0600
	mode := fmt.Sprintf("%#o", info.Mode()&0xFFFF)
	kvm_path := fmt.Sprintf("%s/dev/kvm", state.Root)
	cmd := exec.Command("/usr/bin/nsenter", "--target", fmt.Sprintf("%d", state.Pid), "--mount", "--",
		"/usr/bin/mknod", "-m", mode, kvm_path, "c", "10", "232")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to run mknod: %s: %v: %s", kvm_path, err.Error(), output)
		return
	}

	log.Printf("Allowed /dev/kvm in new container: %s %s", kvm_path, state.ID)
}

func main() {
	var state State

	logwriter, err := syslog.New(syslog.LOG_NOTICE, "oci-kvm-hook")
	if err == nil {
		log.SetOutput(logwriter)
	}

	if err := json.NewDecoder(os.Stdin).Decode(&state); err != nil {
		log.Fatalf("Invalid json passed to OCI hook: %v", err.Error())
	}

	command := map[bool]string{true: "prestart", false: "poststop"}[state.Pid > 0]
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if command == "prestart" {
		allowKvm(state)
	}
}
