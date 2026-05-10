<script>
  const { datapoints = [] } = $props();
  const dp = $derived(datapoints.find((d) => d.tag === 'h1'));
</script>

{#if dp}
  <h1>{dp.content}</h1>
{/if}
