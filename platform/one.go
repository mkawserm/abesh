package platform

import (
	"errors"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"go.uber.org/zap"
)

var ErrCapabilityNotFound = errors.New("capability is not found in the global registry")
var ErrCapabilityCategoryIsNotService = errors.New("capability category is not service")
var ErrTriggerNotRegistered = errors.New("the requested trigger has not been registered")

type One struct {
	triggers           map[string]iface.ITrigger
	capabilityRegistry *registry.CapabilityRegistry
}

func (o *One) SetupCapabilities(manifest *model.Manifest) error {
	var err error

	// configure all capability
	// separate triggers from other
	// capability
	for _, v := range manifest.Capabilities {
		capability := registry.GlobalRegistry().GetCapability(v.ContractId)
		if capability == nil {
			logger.L(constant.Name).Error("capability not found",
				zap.String("contract_id", v.ContractId),
			)
			return ErrCapabilityNotFound
		}

		if capability.Category() == string(constant.CategoryTrigger) {
			newCapability := capability.New()
			newCapabilityTrigger := newCapability.(iface.ITrigger)

			err = newCapabilityTrigger.Setup()
			if err != nil {
				return err
			}

			err = newCapabilityTrigger.SetValues(v.Values)
			if err != nil {
				return err
			}

			o.triggers[newCapability.ContractId()] = newCapabilityTrigger
		} else {
			newCapability := capability.New()
			err = newCapability.Setup()
			if err != nil {
				return err
			}

			err = newCapability.SetValues(v.Values)
			if err != nil {
				return err
			}

			o.capabilityRegistry.RegisterCapability(newCapability)
		}
	}

	return nil
}

func (o *One) SetupTriggers(service iface.IService, triggers []*model.TriggerManifest) error {
	var err error

	for _, t := range triggers {
		trigger := o.triggers[t.ContractId]
		if trigger == nil {
			return ErrTriggerNotRegistered
		}

		err = trigger.AddService(o.capabilityRegistry, t.TriggerValues, service)

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *One) SetupServices(manifest *model.Manifest) error {
	var err error

	//Configure services
	for _, s := range manifest.Services {
		capability := registry.GlobalRegistry().GetCapability(s.ContractId)

		if capability == nil {
			return ErrCapabilityNotFound
		}

		service := capability.New().(iface.IService)

		if service == nil {
			return ErrCapabilityCategoryIsNotService
		}

		err = service.Setup()
		if err != nil {
			return err
		}

		err = service.SetValues(s.Values)
		if err != nil {
			return err
		}

		err = o.SetupTriggers(service, s.Triggers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *One) Setup(manifest *model.Manifest) error {
	var err error
	o.triggers = make(map[string]iface.ITrigger)
	o.capabilityRegistry = registry.NewCapabilityRegistry()

	err = o.SetupCapabilities(manifest)
	if err != nil {
		return err
	}

	err = o.SetupServices(manifest)
	if err != nil {
		return err
	}

	return nil
}

func (o *One) Run() error {

	return nil
}
