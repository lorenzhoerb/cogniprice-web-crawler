import type { Job, PaginatedJobs } from "$lib/api/generated";
import { jobService, type JobFilter } from "$lib/api/job.service";

export class JobState {
    loading = $state<boolean>(false);
    error = $state<string | null>();
    jobs = $state<Job[]>([]);
    page = $state(1);
    pageSize = $state(20);
    totalPages = $state<number | null>(null);

    // Remember the last filter used
  private lastFilter: JobFilter = {};

  // Load a page with optional filter
  async load(filter: JobFilter = {}) {
    this.loading = true;
    this.error = null;

    // Merge new filter with last used
    this.lastFilter = { ...this.lastFilter, ...filter };

    try {
      const result: PaginatedJobs = await jobService.getJobs({
        ...this.lastFilter,
        page: filter.page ?? this.page,
        pageSize: filter.pageSize
      });

      // Store only current page entities
      this.jobs = [...result.items];
      this.page = result.page ?? filter.page ?? this.page;
      this.pageSize = result.pageSize ?? this.pageSize;
      this.totalPages = result.totalPages ?? null;
    } catch (e: unknown) {
        if (e instanceof Error) {
            this.error = e.message;
        } else {
            this.error = String(e);
        }
    } finally {
      this.loading = false;
    }
  }

  // Pagination helpers use last filter automatically
  async nextPage() {
    const current = this.page;
    if (this.totalPages && current < this.totalPages) {
      await this.load({ ...this.lastFilter, page: current + 1 });
    }
  }

  async prevPage() {
    const current = this.page;
    if (current > 1) {
      await this.load({ ...this.lastFilter, page: current - 1 });
    }
  }

  // Reset filters if needed
  resetFilters() {
    this.lastFilter = {};
    this.page = 1;
  }
}

export const jobState = new JobState();
