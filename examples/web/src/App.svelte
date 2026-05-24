<script>
  const API       = import.meta.env.VITE_API_URL       ?? 'http://localhost:8080';
  const WAREHOUSE = import.meta.env.VITE_WAREHOUSE_URL ?? 'http://localhost:8081';

  const path    = window.location.pathname;
  const dsMatch = path.match(/^\/ds\/([^/]+)$/);
  const dsId    = dsMatch?.[1] ?? null;
  const slug    = !dsId ? (path.replace(/^\//, '') || null) : null;

  // ── Feed state ──────────────────────────────────────────────────
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

  // ── Datasource detail state ─────────────────────────────────────
  let dsDetail      = $state(null);
  let dsDetailError = $state(null);

  $effect(() => {
    if (!dsId) return;
    fetch(`${WAREHOUSE}/datasources/${dsId}`)
      .then((r) => {
        if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
        return r.json();
      })
      .then((data) => (dsDetail = data))
      .catch((e) => (dsDetailError = e.message));
  });

  // ── View ────────────────────────────────────────────────────────
  let view = $state(localStorage.getItem('view') ?? 'feed');
  $effect(() => { localStorage.setItem('view', view); });

  // ── Shared helpers ──────────────────────────────────────────────
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

  function fmtDate(iso) {
    return new Intl.DateTimeFormat(undefined, {
      dateStyle: 'medium', timeStyle: 'short'
    }).format(new Date(iso));
  }

  const THEMES = ['dark', 'light', 'bw'];
  let theme = $state(localStorage.getItem('theme') ?? 'dark');

  $effect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  });

  function toggleTheme() {
    const i = THEMES.indexOf(theme);
    theme = THEMES[(i + 1) % THEMES.length];
  }

  const THEME_LABELS = { dark: '☀', light: '☾', bw: '◑' };
</script>

<header>
  <span class="page-title">
    {#if dsDetail}
      {dsDetail.description || dsDetail.type}
    {:else if page}
      {page.title}
    {/if}
  </span>

  {#if slug && page && !dsId}
    <nav class="view-tabs">
      <button class="tab" class:active={view === 'feed'} onclick={() => (view = 'feed')}>Feed</button>
      <button class="tab" class:active={view === 'table'} onclick={() => (view = 'table')}>Table</button>
    </nav>
  {/if}

  <button class="theme-toggle" onclick={toggleTheme} aria-label="Toggle theme">
    {THEME_LABELS[theme]}
  </button>
</header>

<main>
  {#if dsId}
    <!-- ── Datasource detail view ── -->
    {#if dsDetailError}
      <p class="status error">{dsDetailError}</p>
    {:else if !dsDetail}
      <p class="status">Loading…</p>
    {:else}
      <a class="back" href="#" onclick={(e) => { e.preventDefault(); history.back(); }}>← Back</a>
      <article class="post">
        {#if dsDetail.description}
          <p class="ds-desc">{dsDetail.description}</p>
        {/if}
        {#each groupDatapoints(dsDetail.datapoints ?? []) as group}
          {#if group.tag === 'li'}
            <ul>
              {#each group.items as dp (dp.id)}
                <li>
                  {dp.content}
                  {#if dp.description}<span class="dp-desc"> — {dp.description}</span>{/if}
                </li>
              {/each}
            </ul>
          {:else}
            {#each group.items as dp (dp.id)}
              <svelte:element this={dp.tag}>{dp.content}</svelte:element>
              {#if dp.description}<p class="dp-desc">{dp.description}</p>{/if}
            {/each}
          {/if}
        {/each}
      </article>
    {/if}

  {:else if !slug}
    <p class="status">No slug in URL — navigate to <code>/my-page-slug</code></p>

  {:else if error}
    <p class="status error">{error}</p>

  {:else if !page}
    <p class="status">Loading…</p>

  {:else if datasources.length === 0}
    <p class="status">No content.</p>

  {:else if view === 'table'}
    <!-- ── Table view ── -->
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Description</th>
            <th>Type</th>
            <th>Last synced</th>
          </tr>
        </thead>
        <tbody>
          {#each datasources as ds (ds.id)}
            <tr>
              <td>
                <a class="ds-link" href="/ds/{ds.id}">{ds.description || '—'}</a>
              </td>
              <td class="mono">{ds.type}</td>
              <td class="muted">{fmtDate(ds.updated_at)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

  {:else}
    <!-- ── Feed view ── -->
    <div class="feed">
      {#each datasources as ds (ds.id)}
        <article class="post">
          {#if ds.description}
            <p class="ds-desc">{ds.description}</p>
          {/if}
          {#each groupDatapoints(ds.datapoints) as group}
            {#if group.tag === 'li'}
              <ul>
                {#each group.items as dp (dp.id)}
                  <li>
                    {dp.content}
                    {#if dp.description}<span class="dp-desc"> — {dp.description}</span>{/if}
                  </li>
                {/each}
              </ul>
            {:else}
              {#each group.items as dp (dp.id)}
                {#if dp.tag === 'h1'}
                  <h1><a class="heading-link" href="/ds/{ds.id}">{dp.content}</a></h1>
                {:else}
                  <svelte:element this={dp.tag}>{dp.content}</svelte:element>
                {/if}
                {#if dp.description}<p class="dp-desc">{dp.description}</p>{/if}
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

  :global(:root[data-theme="bw"]) {
    --bg:       #ffffff;
    --bg1:      #f8f8f8;
    --bg2:      #eeeeee;
    --fg:       #000000;
    --fg2:      #111111;
    --fg3:      #555555;
    --yellow:   #000000;
    --orange:   #cc0000;
    --aqua:     #000000;
    --border:   #cccccc;
  }

  :global([data-theme="bw"] body) {
    font-family: 'Times New Roman', Times, serif;
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
    min-width: 80px;
  }

  .view-tabs {
    display: flex;
    gap: 2px;
    background: var(--bg2);
    border-radius: 6px;
    padding: 3px;
  }

  .tab {
    background: none;
    border: none;
    color: var(--fg3);
    font-family: system-ui, sans-serif;
    font-size: 12px;
    font-weight: 500;
    padding: 3px 12px;
    border-radius: 4px;
    cursor: pointer;
    transition: background 0.1s, color 0.1s;
  }

  .tab:hover { color: var(--fg); }

  .tab.active {
    background: var(--bg1);
    color: var(--fg);
    box-shadow: 0 1px 2px rgba(0,0,0,0.15);
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
    min-width: 80px;
    text-align: right;
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

  .back {
    display: inline-block;
    font-family: system-ui, sans-serif;
    font-size: 13px;
    color: var(--fg3);
    text-decoration: none;
    margin-bottom: 28px;
    transition: color 0.1s;
  }

  .back:hover { color: var(--fg); }

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

  .heading-link {
    color: inherit;
    text-decoration: none;
  }

  .heading-link:hover {
    text-decoration: underline;
    text-decoration-color: var(--yellow);
    text-underline-offset: 4px;
  }

  .ds-desc {
    font-family: system-ui, sans-serif;
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.07em;
    text-transform: uppercase;
    color: var(--fg3);
    margin-bottom: 16px;
  }

  .dp-desc {
    font-size: 13px;
    font-style: italic;
    color: var(--fg3);
    margin-top: -0.4rem;
    margin-bottom: 0.75rem;
  }

  /* ── Table view ── */
  .table-wrap {
    background: var(--bg1);
    border: 1px solid var(--border);
    border-radius: 6px;
    overflow: hidden;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: system-ui, sans-serif;
    font-size: 14px;
  }

  thead tr { border-bottom: 1px solid var(--border); }

  th {
    text-align: left;
    padding: 10px 16px;
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.05em;
    text-transform: uppercase;
    color: var(--fg3);
  }

  td {
    padding: 12px 16px;
    border-bottom: 1px solid var(--bg2);
    color: var(--fg2);
  }

  tbody tr:last-child td { border-bottom: none; }

  tbody tr:hover td { background: var(--bg2); }

  .ds-link {
    color: var(--fg);
    font-weight: 500;
    text-decoration: none;
  }

  .ds-link:hover { text-decoration: underline; }

  .mono { font-family: monospace; font-size: 13px; }

  .muted { color: var(--fg3); font-size: 13px; }

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
