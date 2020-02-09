package systemctl

import (
	"os/exec"
	"strings"
)

// IsEnabled checks if a systemd service is enabled
func IsEnabled(service string) (bool, error) {
	return boolCheck(service, "enabled", "disabled")
}

// IsActive checks if a systemd service is active
func IsActive(service string) (bool, error) {
	return boolCheck(service, "active", "inactive")
}

// IsEnabled checks if a systemd service is enabled
func boolCheck(service, positive, negative string) (bool, error) {

	bytes, err := exec.Command("/bin/systemctl", "is-"+positive, service).Output()
	status := strings.TrimSuffix(string(bytes), "\n")

	if status == positive {
		return true, nil
	}

	if status == negative {
		return false, nil
	}

	return false, err
}
