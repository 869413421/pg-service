package encoder

import jsoniter "github.com/json-iterator/go"

var JsonHandler jsoniter.API

func init() {
	JsonHandler = jsoniter.ConfigCompatibleWithStandardLibrary
}
