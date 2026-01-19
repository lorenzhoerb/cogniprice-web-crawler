<script lang="ts">
  import { jobState } from '$lib/stores/jobState.svelte';
	import { onMount } from 'svelte';
  import { debounce } from 'lodash-es';



  onMount(() => {
    jobState.load({pageSize: 4});
  })

  const nextPage = () => jobState.nextPage();
  const prevPage = () => jobState.prevPage();

  // Debounced load function
  const debouncedLoad = debounce((url: string) => {
    jobState.load({ pageSize: 4, url });
  }, 300); // 300ms debounce
</script>

<div class="max-w-4xl mx-auto p-4">
  <h1 class="text-2xl font-bold mb-4">Jobs</h1>

  <input
  placeholder="example.com"
  type="text"
  class="border rounded px-3 py-2 w-full max-w-sm"
  oninput={(e) => debouncedLoad((e.target as HTMLInputElement).value)}
/>

  {#if jobState.loading}
    <p class="text-gray-500">Loading jobs…</p>
  {:else if jobState.error}
    <p class="text-red-500 font-semibold">Error: {jobState.error}</p>
  {:else}
    {#if jobState.jobs.length === 0}
      <p class="text-gray-600">No jobs found.</p>
    {:else}
      <ul class="space-y-3">
        {#each jobState.jobs as job}
          <li class="p-4 bg-white shadow rounded-md flex justify-between items-center hover:bg-gray-50 transition">
            <div>
              <p class="font-medium text-gray-800">{job.url}</p>
              <p class="text-sm text-gray-500">
                Next run: {job.nextRunAt} — Status: {job.status}
              </p>
            </div>
          </li>
        {/each}
      </ul>

      <!-- Pagination -->
      <div class="flex items-center justify-between mt-4">
        <button
          class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
          onclick={prevPage}
          disabled={jobState.page <= 1}
        >
          Prev
        </button>

        <p class="text-gray-700">
          Page {jobState.page} of {jobState.totalPages ?? '?'}
        </p>

        <button
          class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
          onclick={nextPage}
          disabled={jobState.totalPages !== null && jobState.page >= jobState.totalPages}
        >
          Next
        </button>
      </div>
    {/if}
  {/if}
</div>