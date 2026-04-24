walk(
  if type == "object"
     and .additionalProperties
     and .additionalProperties.type == "object"
     and (.additionalProperties.additionalProperties | not)
  then .additionalProperties = true
  # Array schemas without items default to string — Go codegen requires items to be present
  elif type == "object"
     and .type == "array"
     and (.items | not)
  then .items = {"type": "string"}
  else .
  end
)
# Duplicate path with {projectName} but empty parameters — remove to avoid validation error
| del(.paths["/v1/project/{projectName}/preference/"])

# The upstream spec tags some operations with multiple tags, causing the Go generator to
# emit the same request struct in several files (e.g. api_test_case.go AND
# api_batch_operations.go), which fails compilation with "redeclared in this block".
# Keep only the first tag so every operation lands in exactly one generated file.
| .paths |= with_entries(
    .value |= with_entries(
      if (.value | type == "object") and (.value.tags | type == "array") and (.value.tags | length > 1) then
        .value.tags = [.value.tags[0]]
      else .
      end
    )
  )

# Fix: the Go openapi-generator does not call the parent's ToMap() from a child's
# ToMap() when allOf inheritance is used, so discriminator fields (e.g.
# manualScenarioType) are silently dropped from MarshalJSON output.
# Flatten allOf[ref + inline-properties] schemas into a single properties block
# so every field ends up in one ToMap() and the serialized JSON is complete.
| .components.schemas as $schemas
| .components.schemas |= (
    to_entries | map(
        if (
            .value | type == "object"
            and has("allOf")
            and (.allOf | length == 2)
            and (.allOf[0] | has("$ref"))
            and (.allOf[1] | has("properties"))
        ) then
            .value.allOf[0]["$ref"] as $ref |
            ($ref | ltrimstr("#/components/schemas/")) as $parentKey |
            ($schemas[$parentKey].properties // {}) as $parentProps |
            ($schemas[$parentKey].required // []) as $parentRequired |
            .value |= {
                required: ((.required // []) + $parentRequired | unique),
                type: "object",
                properties: ($parentProps + .allOf[1].properties)
            }
        else .
        end
    ) | from_entries
)
