# Terraform provider to allow one to define a map variable by parts

## The problem
Sometimes you have one resource, one config file or one something that is built from stuff on other files. For example: a shared hosts file generated from nodes defined in different places, or a ts-proxy service list from many small services. The idea here is when setting up a specific service one could only edit the file of that service, without concentration on specific files or a `svc_foo-setting` that gets joined as `services = [ svc_foo-setting ]`. This is suboptimal and this projects intends to solve that.

## The approach
Essentially, the idea here is to define two resources:

- The group: a map[string]something that holds the final value and gets referenced by items
- The item: a edge in this graph which relates one key of one group into one value

## The challenges
- Make group depend on all the items
- Integrate this into Terraform without making it gross and clunky


## Example

```terraform
resource "accumulator_group" "example" {
  name = "example-accumulator" // this generats a UUID internally
  type = string
}

resource "accumulator_item" "item" {
    group = accumulator_group.example.id
    key = "item"
    value = "example value"
}
```

In this case, `accumulator_group.example.data` should give `{"item": "example value"}` and validate wether ethe value is a string.