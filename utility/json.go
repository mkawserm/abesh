package utility

import (
	"encoding/json"
	"github.com/mkawserm/abesh/errors"
	"github.com/mkawserm/abesh/model"
)

func JSONErrorEvent(err *errors.Error, data interface{}, inputMetadata *model.Metadata, contractId string) *model.Event {
	lang := GetLanguage(inputMetadata.Headers)
	r := &model.HTTPResponseModel{
		Code:    GetErrorResponseCode(err),
		Message: GetErrorMessage(err, lang),
		Lang:    lang,
		Data:    data,
	}

	if data == nil {
		r.Data = make(map[string]interface{})
	}

	dataBytes, _ := json.Marshal(r)

	return model.GenerateOutputEvent(inputMetadata,
		contractId,
		"NOT OK",
		err.GetCode(),
		"application/json", dataBytes)
}

func JSONErrorEventHTTP(err *errors.Error, data interface{}, inputMetadata *model.Metadata, contractId string) *model.Event {
	lang := GetLanguage(inputMetadata.Headers)
	r := &model.HTTPResponseModel{
		Code:    GetErrorResponseCode(err),
		Message: GetErrorMessage(err, lang),
		Lang:    lang,
		Data:    data,
	}

	if data == nil {
		r.Data = make(map[string]interface{})
	}

	dataBytes, _ := json.Marshal(r)

	return model.GenerateOutputEvent(inputMetadata,
		contractId,
		"NOT OK",
		399+err.GetCode()%400,
		"application/json", dataBytes)
}

func JSONSuccessEvent(status *model.Status, data interface{}, inputMetadata *model.Metadata, contractId string) *model.Event {
	lang := GetLanguage(inputMetadata.Headers)
	r := &model.HTTPResponseModel{
		Code:    GetSuccessResponseCode(status),
		Message: GetSuccessMessage(status, lang),
		Lang:    lang,
		Data:    data,
	}

	if data == nil {
		r.Data = make(map[string]interface{})
	}

	dataBytes, _ := json.Marshal(r)

	return model.GenerateOutputEvent(inputMetadata,
		contractId,
		"OK",
		status.GetCode(),
		"application/json", dataBytes)
}

func JSONSuccessEventHTTP(status *model.Status, data interface{}, inputMetadata *model.Metadata, contractId string) *model.Event {
	lang := GetLanguage(inputMetadata.Headers)
	r := &model.HTTPResponseModel{
		Code:    GetSuccessResponseCode(status),
		Message: GetSuccessMessage(status, lang),
		Lang:    lang,
		Data:    data,
	}

	if data == nil {
		r.Data = make(map[string]interface{})
	}

	dataBytes, _ := json.Marshal(r)

	return model.GenerateOutputEvent(inputMetadata,
		contractId,
		"OK",
		199+status.GetCode()%200,
		"application/json", dataBytes)
}
