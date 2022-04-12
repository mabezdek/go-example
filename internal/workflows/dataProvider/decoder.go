package dataProvider

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DecodeTime() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		}

		return data, nil
	}
}

func DecodeObjectID() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(primitive.ObjectID{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return primitive.ObjectIDFromHex(data.(string))
		}

		return data, nil
	}
}

func Decode(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			DecodeTime(),
			DecodeObjectID(),
		),
		Result: result,
	})

	if err != nil {
		return err
	}

	if data := decoder.Decode(input); data != nil {
		return data
	}

	return nil
}
