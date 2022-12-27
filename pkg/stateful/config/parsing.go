package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"
)

type MissingENV struct {
	Name string
}

func (e MissingENV) Error() string {
	return fmt.Sprintf("missing env var: %s", e.Name)
}

func GetEnvFloat64WithDefault(envName string, defaultValue float64) (float64, error) {
	val, err := GetEnvFloat64(envName)
	switch err.(type) {
	case MissingENV:
		return defaultValue, nil
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

func GetEnvFloat64(envName string) (float64, error) {
	val, err := GetEnv(envName)
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func GetEnvInt64WithDefault(envName string, defaultValue int64) (int64, error) {
	val, err := GetEnvInt64(envName)
	switch err.(type) {
	case MissingENV:
		return defaultValue, nil
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

func GetEnvInt64(envName string) (int64, error) {
	val, err := GetEnv(envName)
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func GetEnvInt(envName string) (int, error) {
	val, err := GetEnvInt64(envName)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func GetEnvIntWithDefault(envName string, defaultValue int64) (int, error) {
	val, err := GetEnvInt64WithDefault(envName, defaultValue)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func GetEnvBoolWithDefault(envName string, defaultValue bool) (bool, error) {
	val, err := GetEnvBool(envName)
	switch err.(type) {
	case MissingENV:
		return defaultValue, nil
	}
	if err != nil {
		return false, err
	}
	return val, nil
}

func GetEnvBool(envName string) (bool, error) {
	val, err := GetEnv(envName)
	if err != nil {
		return false, err
	}
	result, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return result, nil
}

func GetEnvWithDefault(envName, defaultValue string) (string, error) {
	val, err := GetEnv(envName)
	switch err.(type) {
	case MissingENV:
		return defaultValue, nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetEnv(envName string) (string, error) {
	val, ok := os.LookupEnv(envName)
	if ok {
		return val, nil
	}

	return "", MissingENV{envName}
}

func GetEnvDuration(envName string) (time.Duration, error) {
	val, err := GetEnv(envName)
	if err != nil {
		return 0, err
	}

	return time.ParseDuration(val)
}

func GetEnvDecodedB64(envName string) ([]byte, error) {
	variable, err := GetEnv(envName)
	if err != nil {
		return nil, err
	}
	decodedValue, err := base64.StdEncoding.DecodeString(variable)
	if err != nil {
		return nil, fmt.Errorf("cannot b64 decode %v", envName)
	}
	return decodedValue, nil
}
