data "accumulator_group" "example" {}

resource "accumulator_item" "item" {
    group = data.accumulator_group.example.id
    key = "example"
    value = "example value"
}

output "expected" {
  value = {
    example = "example value"
  }
}

output "got" {
  value = data.accumulator_group.example.value
}
