resource "accumulator_group" "main" {}

output "expected" {
    value = {
        "a": "A"
        "b": "B"
    }
}

output "got" {
 value = accumulator_group.main.value
}

resource "accumulator_item" "b" {
  group = accumulator_group.main.id
  key = "b"
  value = "B"
}

resource "accumulator_item" "a" {
  group = accumulator_group.main.id
  key = "a"
  value = "A"
}
