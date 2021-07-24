package platform

import (
	"context"
	"errors"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
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

			err = newCapabilityTrigger.SetValues(v.Values)
			if err != nil {
				return err
			}

			err = newCapabilityTrigger.Setup()
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

		logger.L(constant.Name).Debug("adding service to trigger", zap.String("trigger_contract_id", t.ContractId))
		err = trigger.AddService(o.capabilityRegistry, t.Values, service)

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

		err = service.SetValues(s.Values)
		if err != nil {
			return err
		}

		err = service.Setup()
		if err != nil {
			return err
		}

		logger.L(constant.Name).Debug("configuring triggers", zap.String("contract_id", service.ContractId()))
		err = o.SetupTriggers(service, s.Triggers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *One) Setup(manifest *model.Manifest) error {
	timerStart := time.Now()
	var err error
	o.triggers = make(map[string]iface.ITrigger)
	o.capabilityRegistry = registry.NewCapabilityRegistry()
	logger.L(constant.Name).Debug("configuring capabilities")
	err = o.SetupCapabilities(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("configuring services")
	err = o.SetupServices(manifest)
	if err != nil {
		return err
	}

	elapsed := time.Since(timerStart)
	logger.L(constant.Name).Info("setup execution time", zap.Duration("seconds", elapsed))
	return nil
}

func (o *One) Run() {
	timerStart := time.Now()

	idleChan := make(chan struct{})

	go func() {
		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
		sig := <-signChan
		logger.L(constant.Name).Info("shutdown signal received",
			zap.String("signal", sig.String()))

		logger.L(constant.Name).Info("preparing from shutdown")
		logger.L(constant.Name).Info("closing all triggers")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for _, t := range o.triggers {
			err := t.Stop(ctx)
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}

		logger.L(constant.Name).Info("closed all triggers")

		// Actual shutdown trigger.
		close(idleChan)
	}()

	for _, t := range o.triggers {
		trigger := t

		// start all triggers
		go func() {
			err := trigger.Start(context.TODO())
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}()
	}

	logger.L(constant.Name).Info("all triggers are executed")
	elapsed := time.Since(timerStart)

	logger.L(constant.Name).Info("run execution time", zap.Duration("seconds", elapsed))
	// Blocking until the shutdown is complete
	<-idleChan

	logger.L(constant.Name).Info("shutdown complete")
}
