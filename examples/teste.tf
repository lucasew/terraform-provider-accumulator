resource "accumulator_group" "example" {
  name = "example-accumulator" // this generats a UUID internally
  type = string
}

resource "accumulator_item" "item" {
    group = accumulator_group.example.id
    value = "example value"
}