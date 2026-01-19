
# JobInput


## Properties

Name | Type
------------ | -------------
`url` | string
`interval` | string

## Example

```typescript
import type { JobInput } from ''

// TODO: Update the object below with actual values
const example = {
  "url": my.shop.com/product-a,
  "interval": 5s,
} satisfies JobInput

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as JobInput
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


