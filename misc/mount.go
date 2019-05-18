package misc

import (
	"os/exec"
)

func IsMounted(path string) (bool, error) {

	err := exec.Command("/bin/findmnt", path).Run()

	if err == nil {
		return true, nil
	}

	switch err.(type) {
	case *exec.ExitError:
		return false, nil
	default:
		return false, err
	}
}
