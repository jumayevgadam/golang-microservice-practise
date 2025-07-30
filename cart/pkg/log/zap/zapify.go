package zap

import (
	"cart/pkg/log"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// nolint: gocyclo
func zapifyField(field log.Field) zap.Field {
	switch field.Type() {
	case log.FieldTypeNil:
		return zap.Reflect(field.Key(), nil)
	case log.FieldTypeString:
		return zap.String(field.Key(), field.String())
	case log.FieldTypeBinary:
		return zap.Binary(field.Key(), field.Binary())
	case log.FieldTypeBoolean:
		return zap.Bool(field.Key(), field.Bool())
	case log.FieldTypeSigned:
		return zap.Int64(field.Key(), field.Signed())
	case log.FieldTypeUnsigned:
		return zap.Uint64(field.Key(), field.Unsigned())
	case log.FieldTypeFloat:
		return zap.Float64(field.Key(), field.Float())
	case log.FieldTypeTime:
		return zap.Time(field.Key(), field.Time())
	case log.FieldTypeDuration:
		return zap.Duration(field.Key(), field.Duration())
	case log.FieldTypeError:
		return zap.NamedError(field.Key(), field.Error())
	case log.FieldTypeArray:
		return zap.Any(field.Key(), field.Interface())
	case log.FieldTypeAny:
		return zap.Any(field.Key(), field.Interface())
	case log.FieldTypeReflect:
		return zap.Reflect(field.Key(), field.Interface())
	case log.FieldTypeByteString:
		return zap.ByteString(field.Key(), field.Binary())
	case log.FieldTypeStringer:
		return zap.Stringer(field.Key(), field.Interface().(fmt.Stringer))
	default:
		// For when new field type is not added to this func
		panic(fmt.Sprintf("unknown field type: %d", field.Type()))
	}
}

func zapifyFields(fields ...log.Field) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zapifyField(field))
	}

	return zapFields
}
