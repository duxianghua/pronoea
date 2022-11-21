package utils

import (
	"os"
)

func GetCurrentNamespace() string {
	byte, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	namespace := string(byte)
	if err != nil {
		namespace = "default"
	}
	return namespace
}
