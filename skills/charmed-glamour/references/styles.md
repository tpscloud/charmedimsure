# Glamour v2 style configuration

The `ansi.StyleConfig` fields, the shipped configs, and the JSON format in
Glamour v2.0.1. Return to the [main skill](../SKILL.md) for the custom style
example.

## StyleConfig

A style is an `ansi.StyleConfig` from `"charm.land/glamour/v2/ansi"`.

- Blocks: `Document`, `Paragraph`, `BlockQuote`, `Heading`, `H1` through
  `H6`, `Code`, `CodeBlock`, `List`, `Table`, `DefinitionList`, `HTMLBlock`.
- Primitives: `Text`, `Emph`, `Strong`, `Strikethrough`, `Link`, `LinkText`,
  `Image`, `Item`, `Enumeration`, `Task`, `HorizontalRule`.

Optional fields such as colors and booleans use pointers so unset can
inherit; prefixes and suffixes are strings. The custom style example sets
pointer fields with `new(value)`, which requires Go 1.26 or later; the
example module uses Go 1.26.5. On an older toolchain, assign the value to a
variable and take its address.

## Shipped configs

`"charm.land/glamour/v2/styles"` exports `styles.DarkStyleConfig`,
`LightStyleConfig`, `PinkStyleConfig`, `DraculaStyleConfig`,
`TokyoNightStyleConfig`, `ASCIIStyleConfig`, `NoTTYStyleConfig`, plus
`styles.DefaultStyles` (a map by name). The configs are values, so
assigning one to a variable gives a copy that is safe to edit.

## JSON styles

JSON files use the same field names in snake case (except `HorizontalRule`,
whose key is `hr`) and load with `WithStylesFromJSONFile(path)` or
`WithStylesFromJSONBytes(b)`. The `styles/` directory in the module holds the
shipped JSON as reference. Colors are strings: hex or an ANSI index.
