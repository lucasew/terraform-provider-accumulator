resource "accumulator_item" "a" {
  group = accumulator_group.main.id
  key = "a"
  value = "A"
}