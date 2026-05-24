<script>
  const API       = import.meta.env.VITE_API_URL       ?? 'http://localhost:8080';
  const WAREHOUSE = import.meta.env.VITE_WAREHOUSE_URL ?? 'http://localhost:8081';

  const slug = window.location.pathname.replace(/^\//, '') || null;

  let page        = $state(null);
  let datasources = $state([]);
  let error       = $state(null);

  function groupDatapoints(dps) {
    const groups = [];
    for (const dp of dps) {
      const last = groups.at(-1);
      if (last && last.tag === 'li' && dp.tag === 'li') {
        last.items.push(dp);
      } else {
        groups.push({ tag: dp.tag, items: [dp] });
      }
    }
    return groups;
  }
  let theme       = $state(localStorage.getItem('theme') ?? 'dark');

  $effect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  });

  function toggleTheme() {
    theme = theme === 'dark' ? 'light' : 'dark';
  }

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

<header>
  {#if page}<span class="page-title">{page.title}</span>{/if}
  <button class="theme-toggle" onclick={toggleTheme} aria-label="Toggle theme">
    {theme === 'dark' ? '☀' : '☾'}
  </button>
</header>

<main>
  {#if !slug}
    <p class="status">No slug in URL — navigate to <code>/my-page-slug</code></p>
  {:else if error}
    <p class="status error">{error}</p>
  {:else if !page}
    <p class="status">Loading…</p>
  {:else if datasources.length === 0}
    <p class="status">No content.</p>
  {:else}
    <div class="feed">
      {#each datasources as ds (ds.id)}
        <article class="post">
          {#each groupDatapoints(ds.datapoints) as group}
            {#if group.tag === 'li'}
              <ul>
                {#each group.items as dp (dp.id)}
                  <li>{dp.content}</li>
                {/each}
              </ul>
            {:else}
              {#each group.items as dp (dp.id)}
                <svelte:element this={dp.tag}>{dp.content}</svelte:element>
              {/each}
            {/if}
          {/each}
        </article>
      {/each}
    </div>
  {/if}
</main>

<style>
  :global(:root[data-theme="dark"]) {
    --bg:       #282828;
    --bg1:      #3c3836;
    --bg2:      #504945;
    --fg:       #ebdbb2;
    --fg2:      #d5c4a1;
    --fg3:      #bdae93;
    --yellow:   #d79921;
    --orange:   #d65d0e;
    --aqua:     #689d6a;
    --border:   #504945;
  }

  :global(:root[data-theme="light"]) {
    --bg:       #fbf1c7;
    --bg1:      #f2e5bc;
    --bg2:      #ebdbb2;
    --fg:       #3c3836;
    --fg2:      #504945;
    --fg3:      #665c54;
    --yellow:   #b57614;
    --orange:   #af3a03;
    --aqua:     #427b58;
    --border:   #d5c4a1;
  }

  :global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }

  :global(body) {
    background: var(--bg);
    color: var(--fg);
    font-family: Georgia, 'Times New Roman', serif;
    font-size: 17px;
    line-height: 1.7;
    transition: background 0.2s, color 0.2s;
  }

  :global(h1) {
    font-size: 2rem;
    font-weight: 700;
    color: var(--yellow);
    line-height: 1.2;
    margin-bottom: 0.4rem;
  }

  :global(h2) {
    font-size: 1.15rem;
    font-weight: 400;
    color: var(--fg3);
    font-style: italic;
    margin-bottom: 1rem;
  }

  :global(p) {
    color: var(--fg2);
    margin-bottom: 0.75rem;
  }

  :global(ul) {
    padding-left: 1.4rem;
    margin-bottom: 0.75rem;
  }

  :global(li) {
    color: var(--fg2);
    margin-bottom: 0.3rem;
  }

  header {
    position: sticky;
    top: 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 24px;
    background: var(--bg1);
    border-bottom: 1px solid var(--border);
    z-index: 10;
  }

  .page-title {
    font-family: system-ui, sans-serif;
    font-size: 13px;
    font-weight: 600;
    color: var(--fg3);
    letter-spacing: 0.05em;
    text-transform: uppercase;
  }

  .theme-toggle {
    background: none;
    border: 1px solid var(--border);
    color: var(--fg3);
    border-radius: 6px;
    padding: 4px 10px;
    font-size: 14px;
    cursor: pointer;
    line-height: 1;
    transition: border-color 0.15s, color 0.15s;
  }

  .theme-toggle:hover {
    color: var(--fg);
    border-color: var(--fg3);
  }

  main {
    max-width: 680px;
    margin: 0 auto;
    padding: 40px 24px;
  }

  .feed {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .post {
    background: var(--bg1);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 28px 32px;
  }

  .status {
    color: var(--fg3);
    font-family: system-ui, sans-serif;
    font-size: 14px;
  }

  .status code {
    background: var(--bg1);
    padding: 1px 6px;
    border-radius: 4px;
    font-size: 13px;
  }

  .error { color: var(--orange); }
</style>
