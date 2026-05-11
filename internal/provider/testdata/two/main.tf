data "accumulator_group" "main" {}

output "expected" {
    value = {
        "a": "A"
        "b": "B"
    }
}

output "got" {
 value = data.accumulator_group.main.value
}

resource "accumulator_item" "b" {
  group = data.accumulator_group.main.id
  key = "b"
  value = "B"
}

resource "accumulator_item" "a" {
  group = data.accumulator_group.main.id
  key = "a"
  value = "A"
}
