<script>
  import H1 from './components/H1.svelte';
  import H2 from './components/H2.svelte';
  import P  from './components/P.svelte';
  import Li from './components/Li.svelte';

  const API       = import.meta.env.VITE_API_URL       ?? 'http://localhost:8080';
  const WAREHOUSE = import.meta.env.VITE_WAREHOUSE_URL ?? 'http://localhost:8081';

  const slug = window.location.pathname.replace(/^\//, '') || null;

  let page        = $state(null);
  let datasources = $state([]);
  let error       = $state(null);

  $effect(() => {
    if (!slug) return;

    fetch(`${API}/slug/${slug}`)
      .then((r) => {
        if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
        return r.json();
      })
      .then(async (data) => {
        page = data;

        const ids = [...new Set(
          (data.boxes ?? []).map((b) => b.datasource_id).filter(Boolean)
        )];

        datasources = await Promise.all(
          ids.map((id) =>
            fetch(`${WAREHOUSE}/datasources/${id}`)
              .then((r) => {
                if (!r.ok) throw new Error(`datasource ${id}: ${r.status} ${r.statusText}`);
                return r.json();
              })
          )
        );
      })
      .catch((e) => (error = e.message));
  });
</script>

{#if !slug}
  <pre>No slug in URL. Try navigating to /my-page-slug</pre>
{:else if error}
  <pre>Error: {error}</pre>
{:else if page}
  {#each datasources as ds}
    <H1 datapoints={ds.datapoints} />
    <H2 datapoints={ds.datapoints} />
    <P  datapoints={ds.datapoints} />
    <Li datapoints={ds.datapoints} />
  {/each}
{:else}
  <pre>Loading…</pre>
{/if}
