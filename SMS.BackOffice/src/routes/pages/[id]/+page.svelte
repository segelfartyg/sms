<script lang="ts">
	import { goto } from '$app/navigation';
	import { api, type Box, type Datasource, type Page } from '$lib/api';

	let { data } = $props();

	// ── Page state ──────────────────────────────────────────────
	let page: Page = $state({} as Page);
	let editingPage = $state(false);
	let pageTitle = $state('');
	let pageSlug = $state('');
	let pageError = $state('');
	let pageSaving = $state(false);
	$effect(() => { page = data.page; });

	async function savePageEdits() {
		if (!pageTitle.trim() || !pageSlug.trim()) {
			pageError = 'Title and slug are required.';
			return;
		}
		pageSaving = true;
		pageError = '';
		try {
			page = await api.pages.update(page.id, pageTitle.trim(), pageSlug.trim());
			editingPage = false;
		} catch (e: unknown) {
			pageError = e instanceof Error ? e.message : 'Save failed.';
		} finally {
			pageSaving = false;
		}
	}

	function startEditPage() {
		pageTitle = page.title;
		pageSlug = page.slug;
		pageError = '';
		editingPage = true;
	}

	async function deletePage() {
		if (!confirm('Delete this page and all its boxes?')) return;
		await api.pages.delete(page.id);
		goto('/pages');
	}

	async function exportReadme() {
		const dsIds = [...new Set(boxes.map((b) => b.datasource_id).filter(Boolean) as string[])];
		const lines: string[] = [`# ${page.title}`, ``, `**Slug:** \`${page.slug}\``, ``];

		if (dsIds.length === 0) {
			lines.push('_No datasources linked to this page._');
		} else {
			lines.push('## Datasources', '');
			for (const id of dsIds) {
				const ds = datasources.find((d) => d.id === id);
				const label = ds ? `${ds.description || ds.type} (type: \`${ds.type}\`)` : id;
				lines.push(`### ${label}`, '');
				const dps = await api.datapoints.list(id).catch(() => []);
				if (dps.length === 0) {
					lines.push('_No datapoints._', '');
				} else {
					for (const dp of dps) {
						if (dp.tag === 'h1') lines.push(`# ${dp.content}`);
						else if (dp.tag === 'h2') lines.push(`## ${dp.content}`);
						else if (dp.tag === 'li') lines.push(`- ${dp.content}`);
						else lines.push(dp.content);
					}
					lines.push('');
				}
			}
		}

		const blob = new Blob([lines.join('\n')], { type: 'text/markdown' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${page.slug}-datasources.md`;
		a.click();
		URL.revokeObjectURL(url);
	}

	// ── Datasource state ─────────────────────────────────────────
	let datasources = $state<Datasource[]>([]);
	$effect(() => {
		api.datasources.list().then((ds) => { datasources = ds; }).catch(() => {});
	});

	// ── Box state ────────────────────────────────────────────────
	let boxes = $state<Box[]>([]);
	$effect(() => { boxes = [...(data.page.boxes ?? [])].sort((a, b) => a.position - b.position); });
	let showBoxForm = $state(false);
	let editingBox = $state<Box | null>(null);

	let boxPosition = $state(0);
	let boxDatasourceId = $state('');
	let boxError = $state('');
	let boxSaving = $state(false);

	function openNewBox() {
		editingBox = null;
		boxPosition = boxes.length;
		boxDatasourceId = '';
		boxError = '';
		showBoxForm = true;
	}

	function openEditBox(box: Box) {
		editingBox = box;
		boxPosition = box.position;
		boxDatasourceId = box.datasource_id ?? '';
		boxError = '';
		showBoxForm = true;
	}

	function cancelBoxForm() {
		showBoxForm = false;
		editingBox = null;
		boxError = '';
	}

	async function saveBox() {
		boxSaving = true;
		boxError = '';
		try {
			if (editingBox) {
				const updated = await api.boxes.update(
					page.id,
					editingBox.id,
					boxPosition,
					boxDatasourceId || undefined
				);
				boxes = boxes.map((b) => (b.id === updated.id ? updated : b));
			} else {
				const created = await api.boxes.create(
					page.id,
					boxPosition,
					boxDatasourceId || undefined
				);
				boxes = [...boxes, created];
			}
			boxes = [...boxes].sort((a, b) => a.position - b.position);
			showBoxForm = false;
			editingBox = null;
		} catch (e: unknown) {
			boxError = e instanceof Error ? e.message : 'Save failed.';
		} finally {
			boxSaving = false;
		}
	}

	async function deleteBox(box: Box) {
		if (!confirm(`Delete box at position ${box.position}?`)) return;
		await api.boxes.delete(page.id, box.id);
		boxes = boxes.filter((b) => b.id !== box.id);
	}
</script>

<!-- ── Back nav ───────────────────────────────────────────── -->
<a href="/pages" class="back">← Pages</a>

<!-- ── Page header ──────────────────────────────────────────── -->
<div class="section-header">
	{#if editingPage}
		<div class="card inline-form">
			<div class="form-grid">
				<div class="field">
					<label for="p-title">Title</label>
					<input id="p-title" bind:value={pageTitle} />
				</div>
				<div class="field">
					<label for="p-slug">Slug</label>
					<input id="p-slug" bind:value={pageSlug} />
				</div>
			</div>
			{#if pageError}<p class="error-msg">{pageError}</p>{/if}
			<div class="form-actions">
				<button class="ghost" onclick={() => (editingPage = false)}>Cancel</button>
				<button class="primary" onclick={savePageEdits} disabled={pageSaving}>
					{pageSaving ? 'Saving…' : 'Save'}
				</button>
			</div>
		</div>
	{:else}
		<div class="page-title-row">
			<div>
				<h1>{page.title}</h1>
				<span class="slug">{page.slug}</span>
			</div>
			<div class="page-actions">
				<button class="ghost" onclick={exportReadme}>Export README</button>
				<button class="ghost" onclick={startEditPage}>Edit</button>
				<button class="danger" onclick={deletePage}>Delete Page</button>
			</div>
		</div>
	{/if}
</div>

<!-- ── Boxes ─────────────────────────────────────────────────── -->
<div class="boxes-header">
	<h2>Boxes <span class="count">{boxes.length}</span></h2>
	{#if !showBoxForm}
		<button class="primary" onclick={openNewBox}>+ Add Box</button>
	{/if}
</div>

{#if showBoxForm}
	<div class="card box-form">
		<h3>{editingBox ? 'Edit Box' : 'New Box'}</h3>
		<div class="form-row">
			<div class="field">
				<label for="b-datasource">Datasource</label>
				<select id="b-datasource" bind:value={boxDatasourceId}>
					<option value="">— none —</option>
					{#each datasources as ds (ds.id)}
						<option value={ds.id}>{ds.description || ds.type} ({ds.id})</option>
					{/each}
				</select>
			</div>
			<div class="field field-sm">
				<label for="b-pos">Position</label>
				<input id="b-pos" type="number" bind:value={boxPosition} min="0" />
			</div>
		</div>
		{#if boxError}<p class="error-msg">{boxError}</p>{/if}
		<div class="form-actions">
			<button class="ghost" onclick={cancelBoxForm}>Cancel</button>
			<button class="primary" onclick={saveBox} disabled={boxSaving}>
				{boxSaving ? 'Saving…' : editingBox ? 'Save Changes' : 'Add Box'}
			</button>
		</div>
	</div>
{/if}

{#if boxes.length === 0}
	<div class="empty">No boxes yet. Add one above.</div>
{:else}
	<div class="box-list">
		{#each boxes as box (box.id)}
			<div class="card box-card">
				<div class="box-meta">
					<span class="pos-label">pos {box.position}</span>
					{#if box.datasource_id}
						<span class="ds-label">
							{datasources.find((d) => d.id === box.datasource_id)?.description || box.datasource_id}
						</span>
					{/if}
				</div>
				<div class="box-actions">
					<button class="ghost btn-sm" onclick={() => openEditBox(box)}>Edit</button>
					<button class="danger btn-sm" onclick={() => deleteBox(box)}>Delete</button>
				</div>
			</div>
		{/each}
	</div>
{/if}

<style>
	.back {
		display: inline-block;
		color: #6b7280;
		text-decoration: none;
		font-size: 13px;
		margin-bottom: 20px;
	}

	.back:hover {
		color: #1a1a2e;
	}

	.section-header {
		margin-bottom: 32px;
	}

	.page-title-row {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 16px;
	}

	h1 {
		font-size: 22px;
		font-weight: 700;
	}

	.slug {
		font-family: monospace;
		font-size: 12px;
		color: #6b7280;
	}

	.page-actions {
		display: flex;
		gap: 8px;
		flex-shrink: 0;
	}

	.card {
		background: #fff;
		border-radius: 10px;
		border: 1px solid #e5e7eb;
	}

	.inline-form {
		padding: 20px;
	}

	.form-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 12px;
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 100px;
		gap: 12px;
		margin-bottom: 12px;
	}

	.form-actions {
		display: flex;
		gap: 8px;
		justify-content: flex-end;
		margin-top: 16px;
	}

	.boxes-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 16px;
	}

	h2 {
		font-size: 16px;
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.count {
		background: #f3f4f6;
		color: #6b7280;
		font-size: 11px;
		font-weight: 600;
		padding: 1px 7px;
		border-radius: 99px;
	}

	.box-form {
		padding: 20px;
		margin-bottom: 16px;
	}

	.box-form h3 {
		font-size: 14px;
		font-weight: 600;
		margin-bottom: 16px;
	}

	.box-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.box-card {
		padding: 16px;
	}

	.box-meta {
		display: flex;
		align-items: center;
		gap: 10px;
		margin-bottom: 10px;
	}

	.pos-label {
		font-size: 11px;
		color: #9ca3af;
	}

	.ds-label {
		font-size: 11px;
		color: #059669;
		background: #d1fae5;
		padding: 1px 6px;
		border-radius: 4px;
	}

	.box-actions {
		display: flex;
		gap: 6px;
		justify-content: flex-end;
	}

	:global(.btn-sm) {
		padding: 4px 10px;
		font-size: 12px;
		border-radius: 5px;
	}

	.empty {
		text-align: center;
		padding: 40px;
		color: #9ca3af;
		background: #fff;
		border-radius: 10px;
		border: 1px solid #e5e7eb;
	}
</style>
