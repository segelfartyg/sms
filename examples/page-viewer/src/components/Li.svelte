<script>
  const { datapoints = [] } = $props();
  const items = $derived(datapoints.filter((d) => d.tag === 'li'));
</script>

{#if items.length}
  <ul>
    {#each items as item (item.id)}
      <li>{item.content}</li>
    {/each}
  </ul>
{/if}
