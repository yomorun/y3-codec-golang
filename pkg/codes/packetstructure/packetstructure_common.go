package packetstructure

import (
	"reflect"
	"strings"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

var (
	startingToken byte = 0x01
	logger             = utils.Logger.WithPrefix(utils.DefaultLogger, "yomoCodec::packetStructure")
)

type field struct {
	field reflect.StructField
	val   reflect.Value
}

// fieldNameByTag: get fieldName or tagName
func fieldNameByTag(tagName string, field reflect.StructField) string {
	fieldName := field.Name

	tagValue := field.Tag.Get(tagName)
	tagValue = strings.SplitN(tagValue, ",", 2)[0] // TODO: 考虑0x10:name的结构用于提升使用体验
	if tagValue != "" {
		fieldName = tagValue
	}

	return fieldName
}