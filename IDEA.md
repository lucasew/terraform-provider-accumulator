# Terraform provider to build a map incrementally

## The problem
Sometimes a resource, config file, or generated value is assembled from fragments owned by different files or services. For example: a shared hosts file generated from nodes defined in different places, or a ts-proxy service list built from many small services.

The goal is to let each service define only its own contribution locally, without forcing everything into one central file or into an extra aggregation layer such as `services = [svc_foo_setting]`.

## The approach
This provider defines two resources:

- `accumulator_group`: represents the final accumulated map
- `accumulator_item`: contributes one key/value pair to a group

The resulting value is exposed as a computed attribute on the group.

## Semantics
- Each item contributes exactly one entry to the target group.
- The final group value is a map keyed by item keys.
- Duplicate keys in the same group are invalid.
- Group typing is optional.
- If `type` is set, all item values must conform to it.
- If `type` is omitted, item values may be of any type.

## The challenges
- Integrate the aggregation model into Terraform without making usage clunky.
- Make the final group value reflect all items associated with it.

## Example

```terraform
resource "accumulator_group" "example" {
  name = "example-accumulator" // this generates a UUID internally
  type = string
}

resource "accumulator_item" "item" {
  group = accumulator_group.example.id
  key   = "item"
  value = "example value"
}
```

In this case, `accumulator_group.example.data` should produce `{"item": "example value"}` and validate that the value is a string.
