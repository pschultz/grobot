package testAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fgrosse/grobot/log"
	"reflect"
)

func EqualJson(expected interface{}) *EqualJsonMatcher {
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}

	var parsedObject interface{}
	err = json.Unmarshal(expectedJson, &parsedObject)
	if err != nil {
		panic(err)
	}

	return &EqualJsonMatcher{
		Expected: parsedObject,
	}
}

func EqualJsonString(expectedJson string) *EqualJsonMatcher {
	var parsedObject interface{}
	err := json.Unmarshal([]byte(expectedJson), &parsedObject)
	if err != nil {
		panic(err)
	}

	return &EqualJsonMatcher{
		Expected: parsedObject,
	}
}

type EqualJsonMatcher struct {
	Expected interface{}
}

func (m *EqualJsonMatcher) Matches(actual interface{}) (success bool) {
	if actual == nil && m.Expected == nil {
		return false
	}

	if actual == nil {
		return false
	}

	var parsedObject interface{}
	var err error
	switch value := actual.(type) {
	case []byte:
		err = json.Unmarshal(value, &parsedObject)
	case *json.RawMessage:
		err = json.Unmarshal(*value, &parsedObject)
	case string:
		err = json.Unmarshal([]byte(value), &parsedObject)
	default:
		err = errors.New(fmt.Sprintf("EqualJsonMatcher can not handle type: %T", actual))
	}
	if err != nil {
		panic(err)
	}

	if reflect.DeepEqual(m.Expected, parsedObject) {
		return true
	} else {
		log.Debug("Expected\n%#v\nto equal\n%#v", m.Expected, parsedObject)
		return false
	}
}

func (m *EqualJsonMatcher) String() string {
	return fmt.Sprintf("matching %#v", m.Expected)
}
