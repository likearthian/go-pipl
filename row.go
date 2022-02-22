package pipl

import (
	"encoding/json"
	"fmt"
)

type Row = map[string]interface{}

func NewRows(v interface{}) ([]Row, error) {
	var objects []map[string]interface{}
	if v == nil {
		return nil, fmt.Errorf("receive nil, expect object or objects")
	}

	buf, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v. %s", err.Error())
	}

	var vv interface{}
	if err := json.Unmarshal(buf, &vv); err != nil {
		return nil, fmt.Errorf("failed to serialized v. %s", err.Error())
	}

	switch obj := vv.(type) {
	case []interface{}:
		for _, o := range obj {
			objects = append(objects, o.(map[string]interface{}))
		}
	case map[string]interface{}:
		objects = []map[string]interface{}{obj}
	case []map[string]interface{}:
		objects = obj
	}

	return objects, nil
}
