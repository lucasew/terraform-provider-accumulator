resource "accumulator_item" "b" {
  group = accumulator_group.main.id
  key = "b"
  value = "B"
}