package utility

import (
	"fmt"
	"github.com/mkawserm/abesh/errors"
	"github.com/mkawserm/abesh/model"
)

func GetLanguage(headers map[string]string) string {
	lang, found := headers["accept-language"]
	if !found {
		lang, found = headers["Accept-Language"]
		if !found {
			lang = "en"
		}
	}

	return lang
}

func GetErrorResponseCode(err *errors.Error) string {
	return fmt.Sprintf("%s_%d", err.GetPrefix(), err.GetCode())
}

func GetSuccessResponseCode(status *model.Status) string {
	return fmt.Sprintf("%s_%d", status.GetPrefix(), status.GetCode())
}

func GetErrorMessage(err *errors.Error, lang string) string {
	message, found := err.GetParams()[lang]
	if found {
		return message
	}

	return err.GetMessage()
}

func GetSuccessMessage(status *model.Status, lang string) string {
	message, found := status.GetParams()[lang]
	if found {
		return message
	}

	return status.GetMessage()
}
