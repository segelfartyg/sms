<script lang="ts">
	import { api, type Datasource } from '$lib/api';

	let { data } = $props();

	let datasources = $state<Datasource[]>([]);
	$effect(() => { datasources = data.datasources ?? []; });

	let showForm = $state(false);
	let description = $state('');
	let type = $state('content');
	let error = $state('');
	let saving = $state(false);

	async function create() {
		if (!description.trim()) { error = 'Description is required.'; return; }
		saving = true; error = '';
		try {
			const ds = await api.datasources.create(description.trim(), type.trim() || 'content');
			datasources = [ds, ...datasources];
			description = ''; type = 'content'; showForm = false;
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Something went wrong.';
		} finally {
			saving = false;
		}
	}

	async function remove(id: string) {
		if (!confirm('Delete this datasource and all its datapoints?')) return;
		try {
			await api.datasources.delete(id);
			datasources = datasources.filter((d) => d.id !== id);
		} catch (e: unknown) {
			alert(e instanceof Error ? e.message : 'Delete failed.');
		}
	}
</script>

<div class="page-header">
	<h1>Content</h1>
	{#if !showForm}
		<button class="primary" onclick={() => (showForm = true)}>+ New Source</button>
	{/if}
</div>

{#if showForm}
	<div class="card form-card">
		<h2>New Datasource</h2>
		<div class="form-grid">
			<div class="field">
				<label for="ds-desc">Description</label>
				<input id="ds-desc" bind:value={description} placeholder="e.g. Homepage hero text" />
			</div>
			<div class="field">
				<label for="ds-type">Type</label>
				<input id="ds-type" bind:value={type} placeholder="content" />
			</div>
		</div>
		{#if error}<p class="error-msg">{error}</p>{/if}
		<div class="form-actions">
			<button class="ghost" onclick={() => { showForm = false; error = ''; }}>Cancel</button>
			<button class="primary" onclick={create} disabled={saving}>
				{saving ? 'Creating…' : 'Create'}
			</button>
		</div>
	</div>
{/if}

{#if datasources.length === 0}
	<div class="empty">No datasources yet. Create one to get started.</div>
{:else}
	<div class="table-wrap card">
		<table>
			<thead>
				<tr>
					<th>Description</th>
					<th>Type</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each datasources as ds (ds.id)}
					<tr>
						<td><a href="/content/{ds.id}" class="ds-link">{ds.description}</a></td>
						<td class="muted mono">{ds.type}</td>
						<td class="actions">
							<a href="/content/{ds.id}" class="ghost btn-sm">Manage</a>
							<button class="danger btn-sm" onclick={() => remove(ds.id)}>Delete</button>
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

	h1 { font-size: 22px; font-weight: 700; }

	.card {
		background: #fff;
		border-radius: 10px;
		border: 1px solid #e5e7eb;
	}

	.form-card { padding: 20px; margin-bottom: 24px; }
	.form-card h2 { font-size: 15px; font-weight: 600; margin-bottom: 16px; }

	.form-grid {
		display: grid;
		grid-template-columns: 1fr 160px;
		gap: 12px;
	}

	.form-actions {
		display: flex;
		gap: 8px;
		justify-content: flex-end;
		margin-top: 16px;
	}

	.table-wrap { overflow: hidden; }

	table { width: 100%; border-collapse: collapse; }
	thead tr { border-bottom: 1px solid #e5e7eb; }
	th {
		text-align: left;
		padding: 10px 16px;
		font-size: 11px;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #6b7280;
	}
	td { padding: 12px 16px; border-bottom: 1px solid #f3f4f6; }
	tbody tr:last-child td { border-bottom: none; }

	.ds-link { font-weight: 500; color: #4f46e5; text-decoration: none; }
	.ds-link:hover { text-decoration: underline; }

	.muted { color: #9ca3af; font-size: 12px; }
	.mono { font-family: monospace; }

	.actions { display: flex; gap: 6px; justify-content: flex-end; }

	.empty {
		text-align: center;
		padding: 48px;
		color: #9ca3af;
		background: #fff;
		border-radius: 10px;
		border: 1px solid #e5e7eb;
	}
</style>
