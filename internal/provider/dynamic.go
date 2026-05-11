package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dynamicObjectFromGo(input map[string]any) types.Dynamic {
	attributeTypes := make(map[string]attr.Type, len(input))
	attributes := make(map[string]attr.Value, len(input))

	for key, value := range input {
		attrValue := attrValueFromGo(value)
		attributeTypes[key] = attrValue.Type(context.Background())
		attributes[key] = attrValue
	}

	return types.DynamicValue(types.ObjectValueMust(attributeTypes, attributes))
}

func attrValueFromGo(value any) attr.Value {
	switch v := value.(type) {
	case string:
		return types.StringValue(v)
	case bool:
		return types.BoolValue(v)
	case int:
		return types.NumberValue(big.NewFloat(float64(v)))
	case int64:
		return types.NumberValue(big.NewFloat(float64(v)))
	case float64:
		return types.NumberValue(big.NewFloat(v))
	case *big.Float:
		return types.NumberValue(v)
	default:
		return types.StringValue(fmt.Sprintf("%v", v))
	}
}

func goMapFromDynamic(value types.Dynamic) (map[string]any, error) {
	if value.IsNull() || value.IsUnknown() || value.IsUnderlyingValueNull() || value.IsUnderlyingValueUnknown() {
		return map[string]any{}, nil
	}

	objectValue, ok := value.UnderlyingValue().(types.Object)
	if !ok {
		return nil, fmt.Errorf("expected dynamic object, got %T", value.UnderlyingValue())
	}

	out := make(map[string]any, len(objectValue.Attributes()))
	for key, attrValue := range objectValue.Attributes() {
		goValue, err := goValueFromAttr(attrValue)
		if err != nil {
			return nil, fmt.Errorf("attribute %q: %w", key, err)
		}

		out[key] = goValue
	}

	return out, nil
}

func goValueFromDynamic(value types.Dynamic) (any, error) {
	if value.IsNull() || value.IsUnknown() || value.IsUnderlyingValueNull() || value.IsUnderlyingValueUnknown() {
		return nil, nil
	}

	return goValueFromAttr(value.UnderlyingValue())
}

func goValueFromAttr(value attr.Value) (any, error) {
	switch v := value.(type) {
	case types.String:
		return v.ValueString(), nil
	case types.Bool:
		return v.ValueBool(), nil
	case types.Number:
		return v.ValueBigFloat(), nil
	case types.Object:
		out := make(map[string]any, len(v.Attributes()))
		for key, nested := range v.Attributes() {
			goValue, err := goValueFromAttr(nested)
			if err != nil {
				return nil, fmt.Errorf("attribute %q: %w", key, err)
			}

			out[key] = goValue
		}

		return out, nil
	case types.Dynamic:
		return goValueFromDynamic(v)
	default:
		return nil, fmt.Errorf("unsupported dynamic value type %T", value)
	}
}
