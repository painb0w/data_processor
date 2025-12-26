package bootstrap

import "github.com/painb0w/data_processor/internal/services/protoconverter"

func InitProtoConverter() *protoconverter.ProtoConverter {
	return protoconverter.NewProtoConvert()
}
