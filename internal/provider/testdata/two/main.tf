resource "accumulator_group" "main" { 
}

output "expected" {
    value = {
        "a": "A"
        "b": "B"
    }
}

output "got" {
 value = accumulator_group.main.value
}