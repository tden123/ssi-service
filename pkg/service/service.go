package service

import (
	"github.com/pkg/errors"
	"github.com/tbd54566975/ssi-service/pkg/service/did"
	"github.com/tbd54566975/ssi-service/pkg/service/framework"
	"github.com/tbd54566975/ssi-service/pkg/storage"
	"log"
)

// SSIService represents all services and their dependencies independent of transport
type SSIService struct {
	services []framework.Service
}

// NewSSIService creates a new instance of the SSIS which instantiates all services and their
// dependencies independent of transport.
// TODO(gabe) make service loading config-based
func NewSSIService(log *log.Logger) (*SSIService, error) {
	services, err := instantiateServices(log)
	if err != nil {
		errMsg := "could not instantiate the verifiable credentials service"
		log.Printf(errMsg)
		return nil, errors.Wrap(err, errMsg)
	}
	return &SSIService{services: services}, nil
}

func (ssi *SSIService) GetServices() []framework.Service {
	return ssi.services
}

// instantiateServices begins all instantiates and their dependencies
func instantiateServices(log *log.Logger) ([]framework.Service, error) {
	bolt, err := storage.NewBoltDB()
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate BoltDB")
	}
	didService, err := did.NewDIDService(log, []did.Method{did.KeyMethod}, bolt)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate the DID service")
	}
	return []framework.Service{didService}, nil
}
