# terraform-provider-accumulator

PoC of a group + items approach to define parts of a configuration on different files.

The idea is to get rid of index files of some kind, for example, a hosts file generator that depends on the IPs of the machines being set up
or something like that. Like a list of manual references, the idea is to define a group on the place where the index would be defined and items
that add a key on this hashmap.

As the plugin runs as a subprocess of the Terraform tool, it doesn't have actual access to the graph to do the dependency inversion so via the
plugin approach it's essentially impossible in a reliable (no data races/TOCTOU and stuff like this) way.

Ideally, this concept should be core but as life is not a strawberry, it should be implemented first hehe.
