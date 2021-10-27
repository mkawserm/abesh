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
	"runtime"
	"syscall"
	"time"
)

var ErrCapabilityNotFound = errors.New("capability is not found in the global registry")
var ErrTriggerNotRegistered = errors.New("the requested trigger has not been registered")
var ErrAuthorizerNotRegistered = errors.New("the requested authorizer has not been registered")
var ErrRPCNotRegistered = errors.New("the requested rpc has not been registered")
var ErrServiceNotRegistered = errors.New("the requested service has not been registered")

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
	servicesCapability    map[string]iface.IService

	capabilityRegistry *registry.CapabilityRegistry

	sourceSinkMap map[string][]string

	eventDataChannel EventDataChannel

	startCapabilityList []iface.ICapability
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

func (o *One) callSetConfigMap(capability iface.ICapability, values model.ConfigMap) error {
	v, ok := capability.(iface.ISetConfigMap)
	logger.L(constant.Name).Debug("callSetConfigMap info",
		zap.String("contract_id", capability.ContractId()),
		zap.Bool("ok", ok))

	if ok {
		return v.SetConfigMap(values)
	}
	return nil
}

func (o *One) callSetEventTransmitter(capability iface.ICapability) error {
	v, ok := capability.(iface.ISetEventTransmitter)
	logger.L(constant.Name).Debug("callSetEventTransmitter info",
		zap.String("contract_id", capability.ContractId()),
		zap.Bool("ok", ok))

	if ok {
		return v.SetEventTransmitter(o)
	}
	return nil
}

func (o *One) callSetCapabilityRegistry(capability iface.ICapability) error {
	v, ok := capability.(iface.ISetCapabilityRegistry)

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

func (o *One) configureCapabilities(manifest *model.Manifest) error {
	var err error

	// configure all capability
	// separate triggers, authorizers, consumers, rpcs from other
	// capabilities
	for _, v := range manifest.Capabilities {

		capability := registry.GlobalRegistry().GetCapability(v.ContractId)

		if capability == nil {
			logger.L(constant.Name).Error("capability not found",
				zap.String("contract_id", v.ContractId),
			)
			return ErrCapabilityNotFound
		}

		contractIdAssign := capability.ContractId()
		if len(v.NewContractId) != 0 {
			contractIdAssign = v.NewContractId
		}

		newCapability := capability.New()
		err = o.callSetConfigMap(newCapability, v.Values)
		if err != nil {
			return err
		}

		err = o.callSetEventTransmitter(newCapability)
		if err != nil {
			return err
		}

		if capability.Category() == string(constant.CategoryTrigger) {
			newCapabilityTrigger := newCapability.(iface.ITrigger)
			o.triggersCapability[contractIdAssign] = newCapabilityTrigger
		} else if capability.Category() == string(constant.CategoryRPC) {
			newCapabilityRPC := newCapability.(iface.IRPC)
			o.rpcsCapability[contractIdAssign] = newCapabilityRPC
		} else if capability.Category() == string(constant.CategoryService) {
			newCapabilityService := newCapability.(iface.IService)
			o.servicesCapability[contractIdAssign] = newCapabilityService
		} else if capability.Category() == string(constant.CategoryAuthorizer) {
			newCapabilityAuthorizer := newCapability.(iface.IAuthorizer)
			o.authorizersCapability[contractIdAssign] = newCapabilityAuthorizer
		} else if capability.Category() == string(constant.CategoryConsumer) {
			newCapabilityConsumer := newCapability.(iface.IConsumer)
			o.consumersCapability[contractIdAssign] = newCapabilityConsumer
		} else {
			o.capabilityRegistry.RegisterCapability(contractIdAssign, newCapability)
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

	// for rpcs
	for _, c := range o.rpcsCapability {
		if errLocal := o.callSetCapabilityRegistry(c); errLocal != nil {
			return errLocal
		}
		if errLocal := o.callSetup(c); errLocal != nil {
			return errLocal
		}
	}

	// for services
	for _, c := range o.servicesCapability {
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

func (o *One) configureConsumers(manifest *model.Manifest) error {
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

func (o *One) configureTriggers(manifest *model.Manifest) error {
	//Configuring triggers
	for _, s := range manifest.Triggers {
		trigger := o.triggersCapability[s.Trigger]
		if trigger == nil {
			logger.L(constant.Name).Error("trigger not found", zap.String("contract_id", s.Trigger))
			return ErrTriggerNotRegistered
		}

		service := o.servicesCapability[s.Service]
		if service == nil {
			logger.L(constant.Name).Error("service not found", zap.String("contract_id", s.Service))
			return ErrServiceNotRegistered
		}

		var authorizer iface.IAuthorizer
		if len(s.Authorizer) != 0 {
			authorizer = o.authorizersCapability[s.Authorizer]
			if authorizer == nil {
				logger.L(constant.Name).Error("authorizer not found", zap.String("contract_id", s.Authorizer))
				return ErrAuthorizerNotRegistered
			}
		}

		logger.L(constant.Name).Debug("trigger information",
			zap.Any("trigger", s))

		if errLocal := trigger.AddService(authorizer,
			s.AuthorizerExpression,
			s.TriggerValues,
			service); errLocal != nil {
			return errLocal
		}

	}

	return nil
}

func (o *One) configureRPCS(manifest *model.Manifest) error {
	// configuring RPCS
	for _, s := range manifest.RPCS {
		rpc := o.rpcsCapability[s.RPC]
		if rpc == nil {
			return ErrRPCNotRegistered
		}

		if len(s.Authorizer) != 0 {
			authorizer := o.authorizersCapability[s.Authorizer]
			if authorizer == nil {
				return ErrAuthorizerNotRegistered
			}
			if errLocal := rpc.AddAuthorizer(authorizer, s.AuthorizerExpression, s.Method); errLocal != nil {
				return errLocal
			}
			logger.L(constant.Name).Debug("authorizers setup complete", zap.String("contract_id", rpc.ContractId()))
		} else {
			logger.L(constant.Name).Debug("no authorizer defined", zap.String("contract_id", rpc.ContractId()))
		}
	}

	return nil
}

func (o *One) Setup(manifest *model.Manifest) error {
	timerStart := time.Now()
	/* SYSTEM INFORMATION */
	logger.L(constant.Name).Info("Number of cpu", zap.Int("cpu", runtime.NumCPU()))
	logger.L(constant.Name).Info("Number of go routine", zap.Int("goroutine", runtime.NumGoroutine()))

	var err error

	/* INIT ALL DATA */
	o.triggersCapability = make(map[string]iface.ITrigger)
	o.authorizersCapability = make(map[string]iface.IAuthorizer)
	o.consumersCapability = make(map[string]iface.IConsumer)
	o.rpcsCapability = make(map[string]iface.IRPC)
	o.servicesCapability = make(map[string]iface.IService)
	o.startCapabilityList = make([]iface.ICapability, 0, 100)
	o.capabilityRegistry = registry.NewCapabilityRegistry()

	o.sourceSinkMap = make(map[string][]string)
	o.eventDataChannel = make(EventDataChannel, conf.EnvironmentConfigIns().EventBufferSize)
	/* INIT ALL DATA COMPLETE */

	/* CONFIGURE */
	logger.L(constant.Name).Debug("configuring capabilities")
	err = o.configureCapabilities(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("configuring triggers")
	err = o.configureTriggers(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("configuring rpcs")
	err = o.configureRPCS(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("configuring consumers")
	err = o.configureConsumers(manifest)
	if err != nil {
		return err
	}

	logger.L(constant.Name).Debug("assign start capabilities")
	for _, value := range manifest.Start {
		if v, ok := o.triggersCapability[value]; ok {
			o.startCapabilityList = append(o.startCapabilityList, v)
		}

		if v, ok := o.rpcsCapability[value]; ok {
			o.startCapabilityList = append(o.startCapabilityList, v)
		}

		if v, ok := o.servicesCapability[value]; ok {
			o.startCapabilityList = append(o.startCapabilityList, v)
		}

		if v, ok := o.consumersCapability[value]; ok {
			o.startCapabilityList = append(o.startCapabilityList, v)
		}

		if v, ok := o.authorizersCapability[value]; ok {
			o.startCapabilityList = append(o.startCapabilityList, v)
		}

		if v := o.capabilityRegistry.Capability(value); v != nil {
			o.startCapabilityList = append(o.startCapabilityList, v)
		}
	}
	logger.L(constant.Name).Debug("assign start capabilities done")

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

		for index := range consumers {
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

		//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancel()

		logger.L(constant.Name).Info("closing all capabilities")
		for _, c := range o.startCapabilityList {
			err := o.callStop(context.Background(), c)
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}
		logger.L(constant.Name).Info("closed all capabilities")

		// close event data channel
		close(o.eventDataChannel)

		// Actual shutdown trigger.
		close(idleChan)
	}()

	logger.L(constant.Name).Info("starting all")
	// start all capabilities which has start method
	for _, c := range o.startCapabilityList {
		capability := c
		go func() {
			err := o.callStart(context.Background(), capability)
			if err != nil {
				logger.L(constant.Name).Error(err.Error())
			}
		}()
	}

	logger.L(constant.Name).Info("all started")

	elapsed := time.Since(timerStart)

	logger.L(constant.Name).Info("run execution time", zap.Duration("seconds", elapsed))
	// Blocking until the shutdown is complete

	logger.L(constant.Name).Info("Number of go routine", zap.Int("goroutine", runtime.NumGoroutine()))
	<-idleChan

	logger.L(constant.Name).Info("shutdown complete")
}

func (o *One) callStart(context context.Context, capability iface.ICapability) error {
	v, ok := capability.(iface.IStart)
	if ok {
		return v.Start(context)
	}
	return nil
}

func (o *One) callStop(context context.Context, capability iface.ICapability) error {
	v, ok := capability.(iface.IStop)
	if ok {
		return v.Stop(context)
	}
	return nil
}

func Search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}
