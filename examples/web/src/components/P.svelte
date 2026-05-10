<script>
  const { datapoints = [] } = $props();
  const dp = $derived(datapoints.find((d) => d.tag === 'p'));
</script>

{#if dp}
  <p>{dp.content}</p>
{/if}
