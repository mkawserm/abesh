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
var ErrCapabilityCategoryIsNotRPC = errors.New("capability category is not rpc")

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
	rpcsCapability        map[string]iface.IRPC

	capabilityRegistry *registry.CapabilityRegistry

	sourceSinkMap map[string][]string

	eventDataChannel EventDataChannel
}

func (o *One) GetTriggersCapability() map[string]iface.ITrigger {
	return o.triggersCapability
}

func (o *One) GetAuthorizersCapability() map[string]iface.IAuthorizer {
	return o.authorizersCapability
}

func (o *One) GetConsumersCapability() map[string]iface.IConsumer {
	return o.consumersCapability
}

func (o *One) GetCapabilityRegistry() map[string]iface.ICapability {
	return o.capabilityRegistry.Iterator()
}

//func (o *One) GetRPCCapability() map[string]iface.IRPC {
//	return o.rpcsCapability
//}

func (o *One) getConsumers(contractId string) []iface.IConsumer {
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

func (o *One) callSetConfigMap(capability iface.ICapability, values map[string]string) error {
	v, ok := capability.(iface.IConfigMapSetter)
	logger.L(constant.Name).Debug("callSetConfigMap info",
		zap.String("contract_id", capability.ContractId()),
		zap.Bool("ok", ok))

	if ok {
		return v.SetConfigMap(values)
	}
	return nil
}

func (o *One) callSetCapabilityRegistry(capability iface.ICapability) error {
	v, ok := capability.(iface.ICapabilityRegistrySetter)

	logger.L(constant.Name).Debug("callSetCapabilityRegistry info",
		zap.String("contract_id", capability.ContractId()),
		zap.Bool("ok", ok))

	if ok {
		return v.SetCapabilityRegistry(o.capabilityRegistry)
	}
	return nil
}

func (o *One) callSetup(capability iface.ICapability) error {
	v, ok := capability.(iface.ISetup)
	logger.L(constant.Name).Debug("callSetup info",
		zap.String("contract_id", capability.ContractId()),
		zap.Bool("ok", ok))

	if ok {
		return v.Setup()
	}

	return nil
}

func (o *One) setupCapabilities(manifest *model.Manifest) error {
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

		newCapability := capability.New()
		err = o.callSetConfigMap(newCapability, v.Values)
		if err != nil {
			return err
		}

		if capability.Category() == string(constant.CategoryTrigger) {
			newCapabilityTrigger := newCapability.(iface.ITrigger)

			err = newCapabilityTrigger.AddEventTransmitter(o)
			if err != nil {
				return err
			}

			o.triggersCapability[newCapability.ContractId()] = newCapabilityTrigger
		} else if capability.Category() == string(constant.CategoryRPC) {
			newCapabilityRPC := newCapability.(iface.IRPC)

			err = newCapabilityRPC.AddEventTransmitter(o)
			if err != nil {
				return err
			}

			o.rpcsCapability[newCapability.ContractId()] = newCapabilityRPC
		} else if capability.Category() == string(constant.CategoryService) {
			// skip service
			continue
		} else if capability.Category() == string(constant.CategoryAuthorizer) {
			newCapabilityAuthorizer := newCapability.(iface.IAuthorizer)

			o.authorizersCapability[newCapability.ContractId()] = newCapabilityAuthorizer
		} else if capability.Category() == string(constant.CategoryConsumer) {
			newCapabilityConsumer := newCapability.(iface.IConsumer)

			o.consumersCapability[newCapability.ContractId()] = newCapabilityConsumer
		} else {
			o.capabilityRegistry.RegisterCapability(newCapability)
		}
	}

	// for triggers
	for _, c := range o.triggersCapability {
		if errLocal := o.callSetCapabilityRegistry(c); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetup(c); errLocal != nil {
			return errLocal
		}
	}

	// for authorizers
	for _, c := range o.authorizersCapability {
		if errLocal := o.callSetCapabilityRegistry(c); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetup(c); errLocal != nil {
			return errLocal
		}
	}

	// for consumers
	for _, c := range o.consumersCapability {
		if errLocal := o.callSetCapabilityRegistry(c); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetup(c); errLocal != nil {
			return errLocal
		}
	}

	// rpcs
	//for _, c := range o.rpcsCapability {
	//	if errLocal := o.callSetCapabilityRegistry(c); errLocal != nil {
	//		return errLocal
	//	}
	//
	//	if errLocal := o.callSetup(c); errLocal != nil {
	//		return errLocal
	//	}
	//}

	// other capabilities
	for _, c := range o.capabilityRegistry.Iterator() {
		if errLocal := o.callSetCapabilityRegistry(c); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetup(c); errLocal != nil {
			return errLocal
		}
	}

	return nil
}

func (o *One) setupTriggers(service iface.IService, triggerManifests []*model.TriggerManifest) error {
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
			service)

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *One) setupConsumers(manifest *model.Manifest) error {
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

func (o *One) setupServices(manifest *model.Manifest) error {
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

		if service.Category() != string(constant.CategoryService) {
			return ErrCapabilityCategoryIsNotService
		}

		if errLocal := o.callSetConfigMap(service, s.Values); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetCapabilityRegistry(service); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetup(service); errLocal != nil {
			return errLocal
		}

		logger.L(constant.Name).Debug("configuring triggersCapability", zap.String("contract_id", service.ContractId()))
		err = o.setupTriggers(service, s.Triggers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *One) setupRPCAuthorizer(authorizer iface.IAddAuthorizer, manifest *model.RPCAuthorizerManifest) error {
	return authorizer.AddAuthorizer(manifest.Method, manifest.ContractId, manifest.Expression)
}

func (o *One) setupRPCS(manifest *model.Manifest) error {
	// configuring RPCS
	for _, s := range manifest.RPCS {
		capability := registry.GlobalRegistry().GetCapability(s.ContractId)

		if capability == nil {
			return ErrCapabilityNotFound
		}

		rpc := capability.New().(iface.IRPC)

		if rpc == nil {
			return ErrCapabilityCategoryIsNotRPC
		}

		if rpc.Category() != string(constant.CategoryRPC) {
			return ErrCapabilityCategoryIsNotRPC
		}

		if errLocal := o.callSetConfigMap(rpc, s.Values); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetCapabilityRegistry(rpc); errLocal != nil {
			return errLocal
		}

		if errLocal := o.callSetup(rpc); errLocal != nil {
			return errLocal
		}

		rpcAddAuthorizer := rpc.(iface.IAddAuthorizer)
		if rpcAddAuthorizer == nil {
			logger.L(constant.Name).Debug("add authorizer interface is not satisfied so no need to configure authorizers",
				zap.String("contract_id", rpc.ContractId()))
			continue
		}

		logger.L(constant.Name).Debug("configuring authorizers", zap.String("contract_id", rpc.ContractId()))
		for _, a := range s.Authorizers {
			if errLocal := o.setupRPCAuthorizer(rpcAddAuthorizer, a); errLocal != nil {
				return errLocal
			}
		}
		logger.L(constant.Name).Debug("authorizers setup complete", zap.String("contract_id", rpc.ContractId()))
	}

	return nil
}

func (o *One) Setup(manifest *model.Manifest) error {
	timerStart := time.Now()
	var err error

	o.triggersCapability = make(map[string]iface.ITrigger)
	o.authorizersCapability = make(map[string]iface.IAuthorizer)
	o.consumersCapability = make(map[string]iface.IConsumer)
	o.rpcsCapability = make(map[string]iface.IRPC)
	o.capabilityRegistry = registry.NewCapabilityRegistry()

	o.sourceSinkMap = make(map[string][]string)
	o.eventDataChannel = make(EventDataChannel, conf.EnvironmentConfigIns().EventBufferSize)

	logger.L(constant.Name).Debug("configuring capabilities")
	err = o.setupCapabilities(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("configuring services")
	err = o.setupServices(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("configuring rpcs")
	err = o.setupRPCS(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("configuring consumers")
	err = o.setupConsumers(manifest)
	if err != nil {
		return err
	}

	elapsed := time.Since(timerStart)

	logger.L(constant.Name).Info("setup execution time", zap.Duration("seconds", elapsed))
	return nil
}

func (o *One) eventDispatcher() {
	for {
		edc := <-o.eventDataChannel
		if edc.State == 0 {
			break
		}

		consumers := o.getConsumers(edc.ContractId)
		if consumers == nil {
			logger.L(constant.Name).Debug("no consumer is assigned",
				zap.String("source", edc.ContractId), zap.Any("event", edc))
		}

		for index, _ := range consumers {
			// NOTE: may have issues regarding
			// go routine check later
			i := index

			if edc.State == 1 {
				go func() {
					consumer := consumers[i]
					if consumer == nil {
						return
					}
					err := consumer.ConsumeInputEvent(edc.ContractId, edc.Event)
					if err != nil {
						logger.L(constant.Name).Error("error while sending input event data to consumer",
							zap.String("source", edc.ContractId))
					}
				}()
			}

			if edc.State == 2 {
				go func() {
					consumer := consumers[i]
					if consumer == nil {
						return
					}
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
	go o.eventDispatcher()

	// SHUTDOWN HANDLER
	go func() {
		signChan := make(chan os.Signal, 1)
		signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
		sig := <-signChan
		o.eventDataChannel <- EventData{State: 0}
		logger.L(constant.Name).Info("shutdown signal received",
			zap.String("signal", sig.String()))

		logger.L(constant.Name).Info("preparing for shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.L(constant.Name).Info("closing all triggerCapabilities")
		for _, t := range o.triggersCapability {
			err := t.Stop(ctx)
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}
		logger.L(constant.Name).Info("closed all triggerCapabilities")

		logger.L(constant.Name).Info("closing all rpcs")
		for _, t := range o.rpcsCapability {
			err := t.Stop(ctx)
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}
		logger.L(constant.Name).Info("closed all rpcs")

		// close event data channel
		close(o.eventDataChannel)

		// Actual shutdown trigger.
		close(idleChan)
	}()

	for _, t := range o.triggersCapability {
		trigger := t

		// start all triggerCapability
		go func() {
			err := trigger.Start(context.TODO())
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}()
	}

	logger.L(constant.Name).Info("all triggerCapabilities are executed")

	for _, r := range o.rpcsCapability {
		rpc := r
		// start all rpcCapability
		go func() {
			err := rpc.Start(context.TODO())
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}()
	}
	logger.L(constant.Name).Info("all rpcCapabilities are executed")

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
