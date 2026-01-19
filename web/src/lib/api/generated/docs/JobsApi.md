# JobsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**apiV1JobsGet**](JobsApi.md#apiv1jobsget) | **GET** /api/v1/jobs | List jobs |
| [**apiV1JobsIdDelete**](JobsApi.md#apiv1jobsiddelete) | **DELETE** /api/v1/jobs/{id} | Delete a job |
| [**apiV1JobsIdGet**](JobsApi.md#apiv1jobsidget) | **GET** /api/v1/jobs/{id} | Get a job by ID |
| [**apiV1JobsIdPausePost**](JobsApi.md#apiv1jobsidpausepost) | **POST** /api/v1/jobs/{id}/pause | Pause a job |
| [**apiV1JobsIdResumePost**](JobsApi.md#apiv1jobsidresumepost) | **POST** /api/v1/jobs/{id}/resume | Resume a job |
| [**apiV1JobsPost**](JobsApi.md#apiv1jobspost) | **POST** /api/v1/jobs | Create a new job |



## apiV1JobsGet

> PaginatedJobs apiV1JobsGet(page, pageSize, status, url)

List jobs

Retrieves all jobs currently managed by the scheduler, paginated and filterable.

### Example

```ts
import {
  Configuration,
  JobsApi,
} from '';
import type { ApiV1JobsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new JobsApi();

  const body = {
    // number | Current page number (optional)
    page: 56,
    // number | Number of items per page (optional)
    pageSize: 56,
    // JobStatus | Filter jobs by status (optional)
    status: ...,
    // string | Filter jobs by URL substring match (optional)
    url: my.shop.com,
  } satisfies ApiV1JobsGetRequest;

  try {
    const data = await api.apiV1JobsGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **page** | `number` | Current page number | [Optional] [Defaults to `1`] |
| **pageSize** | `number` | Number of items per page | [Optional] [Defaults to `25`] |
| **status** | `JobStatus` | Filter jobs by status | [Optional] [Defaults to `undefined`] [Enum: SCHEDULED, IN_PROGRESS, PAUSED, FAILED] |
| **url** | `string` | Filter jobs by URL substring match | [Optional] [Defaults to `undefined`] |

### Return type

[**PaginatedJobs**](PaginatedJobs.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Paginated list of jobs |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV1JobsIdDelete

> apiV1JobsIdDelete(id)

Delete a job

Deletes a job by its unique ID.

### Example

```ts
import {
  Configuration,
  JobsApi,
} from '';
import type { ApiV1JobsIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new JobsApi();

  const body = {
    // number | Unique ID of the job
    id: 56,
  } satisfies ApiV1JobsIdDeleteRequest;

  try {
    const data = await api.apiV1JobsIdDelete(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Unique ID of the job | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Job deleted successfully (no content) |  -  |
| **404** | Job not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV1JobsIdGet

> Job apiV1JobsIdGet(id)

Get a job by ID

Retrieves a single job by its unique ID.

### Example

```ts
import {
  Configuration,
  JobsApi,
} from '';
import type { ApiV1JobsIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new JobsApi();

  const body = {
    // number | Unique ID of the job
    id: 56,
  } satisfies ApiV1JobsIdGetRequest;

  try {
    const data = await api.apiV1JobsIdGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Unique ID of the job | [Defaults to `undefined`] |

### Return type

[**Job**](Job.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A single job |  -  |
| **404** | Job not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV1JobsIdPausePost

> Job apiV1JobsIdPausePost(id)

Pause a job

Requests to pause a scheduled or running job.

### Example

```ts
import {
  Configuration,
  JobsApi,
} from '';
import type { ApiV1JobsIdPausePostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new JobsApi();

  const body = {
    // number | Unique ID of the job
    id: 56,
  } satisfies ApiV1JobsIdPausePostRequest;

  try {
    const data = await api.apiV1JobsIdPausePost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Unique ID of the job | [Defaults to `undefined`] |

### Return type

[**Job**](Job.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Job paused successfully |  -  |
| **400** | Invalid request |  -  |
| **404** | Job not found |  -  |
| **409** | Job cannot be paused from current state |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV1JobsIdResumePost

> Job apiV1JobsIdResumePost(id)

Resume a job

Requests to resume a paused job.

### Example

```ts
import {
  Configuration,
  JobsApi,
} from '';
import type { ApiV1JobsIdResumePostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new JobsApi();

  const body = {
    // number | Unique ID of the job
    id: 56,
  } satisfies ApiV1JobsIdResumePostRequest;

  try {
    const data = await api.apiV1JobsIdResumePost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `number` | Unique ID of the job | [Defaults to `undefined`] |

### Return type

[**Job**](Job.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Job resumed successfully |  -  |
| **400** | Invalid request |  -  |
| **404** | Job not found |  -  |
| **409** | Job cannot be resumed from current state |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## apiV1JobsPost

> Job apiV1JobsPost(jobInput)

Create a new job

Adds a new job to the scheduler.

### Example

```ts
import {
  Configuration,
  JobsApi,
} from '';
import type { ApiV1JobsPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new JobsApi();

  const body = {
    // JobInput
    jobInput: ...,
  } satisfies ApiV1JobsPostRequest;

  try {
    const data = await api.apiV1JobsPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **jobInput** | [JobInput](JobInput.md) |  | |

### Return type

[**Job**](Job.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Job created successfully |  -  |
| **400** | Invalid request |  -  |
| **409** | Conflict (job already exists) |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

