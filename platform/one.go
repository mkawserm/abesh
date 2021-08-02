package platform

import (
	"context"
	"errors"
	"github.com/mkawserm/abesh/conf"
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
var ErrAuthorizerNotRegistered = errors.New("the requested authorizer has not been registered")

type EventData struct {
	State      uint8 /*0 break 1 input 2 output*/
	ContractId string
	Event      *model.Event
}

type EventDataChannel chan EventData

type One struct {
	triggersCapability    map[string]iface.ITrigger
	authorizersCapability map[string]iface.IAuthorizer
	consumersCapability   map[string]iface.IConsumer
	capabilityRegistry    *registry.CapabilityRegistry

	sourceSinkMap map[string][]string

	eventDataChannel EventDataChannel
}

func (o *One) GetConsumers(contractId string) []iface.IConsumer {
	var ok bool
	var v []string
	var c []iface.IConsumer

	v, ok = o.sourceSinkMap[contractId]
	if !ok {
		return nil
	}

	c = make([]iface.IConsumer, len(v))
	for index, sv := range v {
		c[index], ok = o.consumersCapability[sv]
	}

	return c
}

func (o *One) TransmitInputEvent(contractId string, event *model.Event) error {
	o.eventDataChannel <- EventData{
		State:      1,
		ContractId: contractId,
		Event:      event,
	}
	return nil
}

func (o *One) TransmitOutputEvent(contractId string, event *model.Event) error {
	o.eventDataChannel <- EventData{
		State:      2,
		ContractId: contractId,
		Event:      event,
	}

	return nil
}

func (o *One) SetupCapabilities(manifest *model.Manifest) error {
	var err error

	// configure all capability
	// separate triggersCapability, authorizersCapability, consumersCapability from other
	// capabilities
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

			err = newCapabilityTrigger.SetConfigMap(v.Values)
			if err != nil {
				return err
			}

			err = newCapabilityTrigger.Setup()
			if err != nil {
				return err
			}

			err = newCapabilityTrigger.AddEventTransmitter(o)
			if err != nil {
				return err
			}

			o.triggersCapability[newCapability.ContractId()] = newCapabilityTrigger
		} else if capability.Category() == string(constant.CategoryAuthorizer) {
			newCapability := capability.New()
			newCapabilityAuthorizer := newCapability.(iface.IAuthorizer)

			err = newCapabilityAuthorizer.SetConfigMap(v.Values)
			if err != nil {
				return err
			}

			err = newCapabilityAuthorizer.Setup()
			if err != nil {
				return err
			}

			o.authorizersCapability[newCapability.ContractId()] = newCapabilityAuthorizer
		} else if capability.Category() == string(constant.CategoryConsumer) {
			newCapability := capability.New()
			newCapabilityConsumer := newCapability.(iface.IConsumer)

			err = newCapabilityConsumer.SetConfigMap(v.Values)
			if err != nil {
				return err
			}

			err = newCapabilityConsumer.Setup()
			if err != nil {
				return err
			}

			o.consumersCapability[newCapability.ContractId()] = newCapabilityConsumer
		} else {
			newCapability := capability.New()

			err = newCapability.SetConfigMap(v.Values)
			if err != nil {
				return err
			}

			err = newCapability.Setup()
			if err != nil {
				return err
			}

			o.capabilityRegistry.RegisterCapability(newCapability)
		}
	}

	return nil
}

func (o *One) SetupTriggers(service iface.IService, triggerManifests []*model.TriggerManifest) error {
	var err error

	for _, t := range triggerManifests {
		trigger := o.triggersCapability[t.ContractId]
		if trigger == nil {
			return ErrTriggerNotRegistered
		}

		var authorizer iface.IAuthorizer
		var expression string
		var authorizationHandler iface.AuthorizationHandler

		if t.Authorizer != nil {
			authorizer = o.authorizersCapability[t.Authorizer.ContractId]
			expression = t.Authorizer.Expression
			if authorizer == nil {
				return ErrAuthorizerNotRegistered
			}

			authorizationHandler = authorizer.IsAuthorized
		}

		logger.L(constant.Name).Debug("adding service to trigger", zap.String("trigger_contract_id", t.ContractId))
		err = trigger.AddService(authorizationHandler,
			expression,
			t.Values,
			o.capabilityRegistry,
			service)

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *One) SetupConsumers(manifest *model.Manifest) error {
	for _, cm := range manifest.Consumers {
		v, ok := o.sourceSinkMap[cm.Source]

		if !ok {
			v = make([]string, 0)
		}

		if len(v) == 0 {
			v = append(v, cm.Sink)
		} else {
			if Search(len(v), func(index int) bool {
				return v[index] == cm.Sink
			}) == -1 {
				v = append(v, cm.Sink)
			}
		}

		o.sourceSinkMap[cm.Source] = v
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

		err = service.SetConfigMap(s.Values)
		if err != nil {
			return err
		}

		err = service.Setup()
		if err != nil {
			return err
		}

		logger.L(constant.Name).Debug("configuring triggersCapability", zap.String("contract_id", service.ContractId()))
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

	o.triggersCapability = make(map[string]iface.ITrigger)
	o.authorizersCapability = make(map[string]iface.IAuthorizer)
	o.consumersCapability = make(map[string]iface.IConsumer)
	o.capabilityRegistry = registry.NewCapabilityRegistry()

	o.sourceSinkMap = make(map[string][]string)
	o.eventDataChannel = make(EventDataChannel, conf.EnvironmentConfigIns().EventBufferSize)

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

	logger.L(constant.Name).Debug("configuring consumers")
	err = o.SetupConsumers(manifest)
	if err != nil {
		return err
	}

	elapsed := time.Since(timerStart)

	logger.L(constant.Name).Info("setup execution time", zap.Duration("seconds", elapsed))
	return nil
}

func (o *One) EventDispatcher() {
	for {
		edc := <-o.eventDataChannel
		if edc.State == 0 {
			break
		}

		consumers := o.GetConsumers(edc.ContractId)
		if consumers == nil {
			logger.L(constant.Name).Debug("no consumer is assigned",
				zap.String("source", edc.ContractId), zap.Any("event", edc))
		}

		for _, consumer := range consumers {
			if consumer == nil {
				continue
			}

			if edc.State == 1 {
				go func() {
					err := consumer.ConsumeInputEvent(edc.ContractId, edc.Event)
					if err != nil {
						logger.L(constant.Name).Error("error while sending input event data to consumer",
							zap.String("source", edc.ContractId))
					}
				}()
			}

			if edc.State == 2 {
				go func() {
					err := consumer.ConsumeOutputEvent(edc.ContractId, edc.Event)
					if err != nil {
						logger.L(constant.Name).Error("error while sending output event data to consumer",
							zap.String("source", edc.ContractId))
					}
				}()
			}
		}

	} // FOR LOOP END
}

func (o *One) Run() {
	timerStart := time.Now()

	idleChan := make(chan struct{})

	// EVENT DISPATCHER
	go o.EventDispatcher()

	// SHUTDOWN HANDLER
	go func() {
		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
		sig := <-signChan
		o.eventDataChannel <- EventData{State: 0}
		logger.L(constant.Name).Info("shutdown signal received",
			zap.String("signal", sig.String()))

		logger.L(constant.Name).Info("preparing for shutdown")
		logger.L(constant.Name).Info("closing all triggersCapability")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for _, t := range o.triggersCapability {
			err := t.Stop(ctx)
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}

		logger.L(constant.Name).Info("closed all triggersCapability")

		// close event data channel
		close(o.eventDataChannel)

		// Actual shutdown trigger.
		close(idleChan)
	}()

	for _, t := range o.triggersCapability {
		trigger := t

		// start all triggersCapability
		go func() {
			err := trigger.Start(context.TODO())
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}()
	}

	logger.L(constant.Name).Info("all triggersCapability are executed")
	elapsed := time.Since(timerStart)

	logger.L(constant.Name).Info("run execution time", zap.Duration("seconds", elapsed))
	// Blocking until the shutdown is complete
	<-idleChan

	logger.L(constant.Name).Info("shutdown complete")
}

func Search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}
