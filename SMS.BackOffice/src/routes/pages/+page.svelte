<script lang="ts">
	import { api, type Page } from '$lib/api';

	let { data } = $props();

	let pages = $state<Page[]>([]);
	$effect(() => { pages = data.pages; });
	let showForm = $state(false);
	let title = $state('');
	let slug = $state('');
	let error = $state('');
	let saving = $state(false);

	async function createPage() {
		if (!title.trim() || !slug.trim()) {
			error = 'Title and slug are required.';
			return;
		}
		saving = true;
		error = '';
		try {
			const page = await api.pages.create(title.trim(), slug.trim());
			pages = [page, ...pages];
			title = '';
			slug = '';
			showForm = false;
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Something went wrong.';
		} finally {
			saving = false;
		}
	}

	async function deletePage(id: string) {
		if (!confirm('Delete this page and all its boxes?')) return;
		try {
			await api.pages.delete(id);
			pages = pages.filter((p) => p.id !== id);
		} catch (e: unknown) {
			alert(e instanceof Error ? e.message : 'Delete failed.');
		}
	}

	function cancelForm() {
		showForm = false;
		title = '';
		slug = '';
		error = '';
	}
</script>

<div class="page-header">
	<h1>Pages</h1>
	{#if !showForm}
		<button class="primary" onclick={() => (showForm = true)}>+ New Page</button>
	{/if}
</div>

{#if showForm}
	<div class="card form-card">
		<h2>New Page</h2>
		<div class="form-grid">
			<div class="field">
				<label for="title">Title</label>
				<input id="title" bind:value={title} placeholder="My Page" />
			</div>
			<div class="field">
				<label for="slug">Slug</label>
				<input id="slug" bind:value={slug} placeholder="my-page" />
			</div>
		</div>
		{#if error}<p class="error-msg">{error}</p>{/if}
		<div class="form-actions">
			<button class="ghost" onclick={cancelForm}>Cancel</button>
			<button class="primary" onclick={createPage} disabled={saving}>
				{saving ? 'Creating…' : 'Create Page'}
			</button>
		</div>
	</div>
{/if}

{#if pages.length === 0}
	<div class="empty">No pages yet. Create one to get started.</div>
{:else}
	<div class="table-wrap card">
		<table>
			<thead>
				<tr>
					<th>Title</th>
					<th>Slug</th>
					<th>Created</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each pages as page (page.id)}
					<tr>
						<td>
							<a href="/pages/{page.id}" class="page-link">{page.title}</a>
						</td>
						<td class="slug">{page.slug}</td>
						<td class="muted">{new Date(page.created_at).toLocaleDateString()}</td>
						<td class="actions">
							<a href="/pages/{page.id}" class="ghost btn-sm">Manage</a>
							<button class="danger btn-sm" onclick={() => deletePage(page.id)}>Delete</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

<style>
	.page-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 24px;
	}

	h1 {
		font-size: 22px;
		font-weight: 700;
	}

	.card {
		background: #fff;
		border-radius: 10px;
		border: 1px solid #e5e7eb;
	}

	.form-card {
		padding: 20px;
		margin-bottom: 24px;
	}

	.form-card h2 {
		font-size: 15px;
		font-weight: 600;
		margin-bottom: 16px;
	}

	.form-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 12px;
	}

	.form-actions {
		display: flex;
		gap: 8px;
		justify-content: flex-end;
		margin-top: 16px;
	}

	.table-wrap {
		overflow: hidden;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	thead tr {
		border-bottom: 1px solid #e5e7eb;
	}

	th {
		text-align: left;
		padding: 10px 16px;
		font-size: 11px;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #6b7280;
	}

	td {
		padding: 12px 16px;
		border-bottom: 1px solid #f3f4f6;
	}

	tbody tr:last-child td {
		border-bottom: none;
	}

	.page-link {
		font-weight: 500;
		color: #4f46e5;
		text-decoration: none;
	}

	.page-link:hover {
		text-decoration: underline;
	}

	.slug {
		font-family: monospace;
		font-size: 12px;
		color: #6b7280;
	}

	.muted {
		color: #9ca3af;
		font-size: 12px;
	}

	.actions {
		display: flex;
		gap: 6px;
		justify-content: flex-end;
	}

	:global(.btn-sm) {
		padding: 4px 10px;
		font-size: 12px;
		border-radius: 5px;
		text-decoration: none;
		display: inline-flex;
		align-items: center;
	}

	:global(a.ghost) {
		border: 1px solid #d1d5db;
		color: #374151;
		background: transparent;
		cursor: pointer;
	}

	:global(a.ghost:hover) {
		opacity: 0.8;
	}

	.empty {
		text-align: center;
		padding: 48px;
		color: #9ca3af;
		background: #fff;
		border-radius: 10px;
		border: 1px solid #e5e7eb;
	}
</style>
