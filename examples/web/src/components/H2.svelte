<script>
  const { datapoints = [] } = $props();
  const dp = $derived(datapoints.find((d) => d.tag === 'h2'));
</script>

{#if dp}
  <h2>{dp.content}</h2>
{/if}
