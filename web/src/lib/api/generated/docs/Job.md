
# Job


## Properties

Name | Type
------------ | -------------
`id` | number
`url` | string
`status` | [JobStatus](JobStatus.md)
`interval` | string
`retryAttempts` | number
`pauseRequested` | boolean
`dispatchedAt` | Date
`nextRunAt` | Date
`createdAt` | Date
`updatedAt` | Date

## Example

```typescript
import type { Job } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "url": https://shopify.com/product/1,
  "status": null,
  "interval": 24h,
  "retryAttempts": null,
  "pauseRequested": null,
  "dispatchedAt": null,
  "nextRunAt": null,
  "createdAt": null,
  "updatedAt": null,
} satisfies Job

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Job
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


