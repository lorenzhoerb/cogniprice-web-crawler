import { Configuration, JobsApi, JobStatus, type PaginatedJobs } from './generated';

export interface JobFilter {
	page?: number;
	pageSize?: number;
	status?: JobStatus;
	url?: string;
}

export class JobService {
	private readonly apiClient;
	constructor(private readonly basePath = 'http://localhost:8082') {
		this.apiClient = new JobsApi(new Configuration({ basePath: basePath }));
	}

	async getJobs(jobFilter: JobFilter): Promise<PaginatedJobs> {
		return this.apiClient.apiV1JobsGet(jobFilter);
	}

	async getJobById(id: number) {
		return this.apiClient.apiV1JobsIdGet({ id });
	}

	async createJob(payload: { url: string; interval?: string }) {
		return this.apiClient.apiV1JobsPost({ jobInput: payload });
	}

	async deleteJob(id: number) {
		return this.apiClient.apiV1JobsIdDelete({ id });
	}
}

export const jobService = new JobService();