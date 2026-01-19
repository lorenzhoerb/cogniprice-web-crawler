
# ModelError


## Properties

Name | Type
------------ | -------------
`message` | string
`status` | string
`code` | number
`errors` | [Array&lt;ErrorErrorsInner&gt;](ErrorErrorsInner.md)

## Example

```typescript
import type { ModelError } from ''

// TODO: Update the object below with actual values
const example = {
  "message": The interval field must be a valid duration like '5s',
  "status": INVALID_ARGUMENT,
  "code": 400,
  "errors": null,
} satisfies ModelError

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelError
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


