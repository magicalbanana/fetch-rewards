package handler

const compareVersionSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "description": "Schema to compare semantic versions",
  "type": "object",
  "properties": {
    "data": {
      "type": "object",
      "properties": {
        "compare_from": {
          "description": "Semantic version to be compared from",
          "type": "string",
          "pattern": "^([0-9]+)(\\.[0-9]+)?(\\.[0-9]+)?(-([0-9A-Za-z\\-]+(\\.[0-9A-Za-z\\-]+)*))?(\\+([0-9A-Za-z\\-]+(\\.[0-9A-Za-z\\-]+)*))?$",
          "minLength": 1
        },
        "compare_to": {
          "description": "Semantic version to be compared to",
          "type": "string",
          "pattern": "^([0-9]+)(\\.[0-9]+)?(\\.[0-9]+)?(-([0-9A-Za-z\\-]+(\\.[0-9A-Za-z\\-]+)*))?(\\+([0-9A-Za-z\\-]+(\\.[0-9A-Za-z\\-]+)*))?$",
          "minLength": 1
        }
      },
      "required": [
        "compare_from",
        "compare_to"
      ],
      "additionalProperties": false
    }
  },
  "required": ["data"],
  "additionalProperties": false
}
`
