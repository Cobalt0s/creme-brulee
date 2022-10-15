package config

import "fmt"

type JaegerTraceConfig struct {
	Address             string
	Enabled             bool
	SamplingProbability float64
}

func (c JaegerTraceConfig) GetURL() string {
	return fmt.Sprintf("http://%v:14268/api/traces", c.Address)
}

func NewJaegerTraceConfig() (*JaegerTraceConfig, error) {

	envTraceEnable, err := GetEnvBoolWithDefault("TRACE_ENABLE", false)
	if err != nil {
		return nil, err
	}
	envTraceProbability, err := GetEnvFloat64WithDefault("TRACE_PROBABILITY", 0)
	if err != nil {
		return nil, err
	}

	t := &JaegerTraceConfig{
		Enabled:             envTraceEnable,
		SamplingProbability: envTraceProbability,
	}
	if envTraceEnable {
		envTraceHost, err := GetEnv("TRACE_HOST")
		if err != nil {
			return nil, err
		}
		t.Address = envTraceHost
	}
	return t, nil
}
