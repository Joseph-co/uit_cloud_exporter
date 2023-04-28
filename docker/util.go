package docker

import (
	"os"
)

func GetHostName() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostName, nil
}

func GetIpAddr() (string, error) {

	return "", nil
}
