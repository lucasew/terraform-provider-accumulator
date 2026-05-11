resource "accumulator_group" "example" {
  name = "example-accumulator" // this generats a UUID internally
  type = string
}

resource "accumulator_item" "item" {
    group = accumulator_group.example.id
    key = "example"
    value = "example value"
}

output "expected" {
  value = {
    example = "example value"
  }
}

output "got" {
  value = accumulator_group.example.value
}