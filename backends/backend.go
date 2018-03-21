package backends

import (
	"reflect"

	"fmt"
	"github.com/Graylog2/collector-sidecar/api/graylog"
	"github.com/Graylog2/collector-sidecar/common"
	"github.com/Graylog2/collector-sidecar/system"
	"os/exec"
	"strings"
)

type Backend struct {
	Enabled              *bool
	Id                   string
	Name                 string
	ServiceType          string
	OperatingSystem      string
	ExecutablePath       string
	ConfigurationPath    string
	ExecuteParameters    []string
	ValidationParameters []string
	Template             string
	backendStatus        system.Status
}

func BackendFromResponse(response graylog.ResponseCollectorBackend) *Backend {
	return &Backend{
		Enabled:              common.NewTrue(),
		Id:                   response.Id,
		Name:                 response.Name,
		ServiceType:          response.ServiceType,
		OperatingSystem:      response.OperatingSystem,
		ExecutablePath:       response.ExecutablePath,
		ConfigurationPath:    response.ConfigurationPath,
		ExecuteParameters:    response.ExecuteParameters,
		ValidationParameters: response.ValidationParameters,
		backendStatus:        system.Status{},
	}
}

func (b *Backend) Equals(a *Backend) bool {
	return reflect.DeepEqual(a, b)
}

func (b *Backend) ValidatePreconditions() bool {
	return true
}

func (b *Backend) ValidateConfigurationFile() bool {
	if b.ValidationParameters == nil {
		log.Errorf("[%s] No parameters configured to validate configuration!", b.Name)
		return false
	}

	var parameters []string
	for _, parameter := range b.ValidationParameters {
		if strings.Contains(parameter, "%s") {
			parameters = append(parameters, fmt.Sprintf(parameter, b.ConfigurationPath))
		} else {
			parameters = append(parameters, parameter)
		}
	}
	output, err := exec.Command(b.ExecutablePath, parameters...).CombinedOutput()
	if err != nil {
		soutput := string(output)
		log.Errorf("[%s] Error during configuration validation: %s %s", b.Name, soutput, err)
		return false
	}

	return true
}
