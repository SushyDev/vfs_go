package service

import "fmt"

func serviceError(message string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s\n%w", message, err)
}
