package systemctl

import (
	"os/exec"
)

// Enable a service
// returns true if the service has been enabled, false if it was already enabled
func Enable(service string) (bool, error) {

	// is it enable?
	enabled, err1 := IsEnabled(service)
	if err1 != nil {
		return false, err1
	}

	// already enabled
	if enabled {
		return false, nil
	}

	// enable the service
	if err2 := exec.Command("/bin/systemctl", "enable", service).Run(); err2 != nil {
		return false, err2
	}

	return true, nil
}

// Activate a service
// returns true if the service has been activated, false if it was already active
func Activate(service string) (bool, error) {

	// is it active?
	active, err1 := IsActive(service)
	if err1 != nil {
		return false, err1
	}

	// already active
	if active {
		return false, nil
	}

	// activate the service
	if err2 := exec.Command("/bin/systemctl", "start", service).Run(); err2 != nil {
		return false, err2
	}

	return true, nil
}

// Restart a service
// returns true if the service has been restarted, false if it has only been started
func Restart(service string) (bool, error) {

	activated, err1 := Activate(service)
	if err1 != nil {
		return false, err1
	}

	// Was it only started?
	if activated {
		return false, nil
	}

	// Need a full restart
	if err2 := exec.Command("/bin/systemctl", "restart", service).Run(); err2 != nil {
		return false, err2
	}

	return true, nil
}
