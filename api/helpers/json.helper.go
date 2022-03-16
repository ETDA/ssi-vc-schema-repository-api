package helpers

import (
	"encoding/json"
	"github.com/thedevsaddam/gojsonq/v2"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

func SetJSONValue(raw *json.RawMessage, field string, value interface{}) (*json.RawMessage, error) {
	var mapper core.Map
	err := json.Unmarshal(*raw, &mapper)
	if err != nil {
		return nil, err
	}
	mapper[field] = value

	b, err := json.Marshal(mapper)
	if err != nil {
		return nil, err
	}

	var newRawMessage json.RawMessage
	err = utils.Copy(&newRawMessage, b)
	return &newRawMessage, err
}

func GetJSONValue(raw *json.RawMessage, field string) interface{} {
	if raw == nil {
		return nil
	}

	return gojsonq.New().FromString(string(*raw)).Find(field)
}
