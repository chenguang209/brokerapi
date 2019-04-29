package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi/domain"
)

const (
	invalidServiceDetailsErrorKey = "invalid-service-details"
	instanceIDLogKey              = "instance-id"
	serviceIdMissingKey           = "service-id-missing"
	planIdMissingKey              = "plan-id-missing"
	unknownErrorKey               = "unknown-error"
	apiVersionInvalidKey          = "broker-api-version-invalid"

	bindLogKey    = "bind"
	getBindLogKey = "getBinding"

	unbindLogKey               = "unbind"
	lastOperationLogKey        = "lastOperation"
	lastBindingOperationLogKey = "lastBindingOperation"
)

var (
	serviceIdError        = errors.New("service_id missing")
	planIdError           = errors.New("plan_id missing")
	invalidServiceIDError = errors.New("service-id not in the catalog")
	invalidPlanIDError    = errors.New("plan-id not in the catalog")
)

type APIHandler struct {
	ServiceBroker domain.ServiceBroker
	Logger        lager.Logger
}

func (h APIHandler) respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		h.Logger.Error("encoding response", err, lager.Data{"status": status, "response": response})
	}
}

type brokerVersion struct {
	Major int
	Minor int
}

func getAPIVersion(req *http.Request) brokerVersion {
	var version brokerVersion
	apiVersion := req.Header.Get("X-Broker-API-Version")

	fmt.Sscanf(apiVersion, "%d.%d", &version.Major, &version.Minor)

	return version
}
