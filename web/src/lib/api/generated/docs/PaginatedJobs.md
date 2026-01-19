
# PaginatedJobs


## Properties

Name | Type
------------ | -------------
`items` | [Array&lt;Job&gt;](Job.md)
`page` | number
`pageSize` | number
`totalCount` | number
`totalPages` | number

## Example

```typescript
import type { PaginatedJobs } from ''

// TODO: Update the object below with actual values
const example = {
  "items": null,
  "page": 1,
  "pageSize": 25,
  "totalCount": 120,
  "totalPages": 5,
} satisfies PaginatedJobs

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PaginatedJobs
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


