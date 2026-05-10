<script lang="ts">
	import { api, type Datapoint, type Datasource } from '$lib/api';

	const TAGS = ['h1', 'h2', 'p', 'li'] as const;

	let { data } = $props();

	let datasource = $state<Datasource | undefined>(undefined);
	let datapoints = $state<Datapoint[]>([]);
	$effect(() => { datasource = data.datasource; datapoints = data.datapoints ?? []; });

	let showForm = $state(false);
	let dpTag = $state<string>('h1');
	let dpContent = $state('');
	let dpDescription = $state('');
	let dpError = $state('');
	let dpSaving = $state(false);

	function openForm() {
		dpTag = 'h1'; dpContent = ''; dpDescription = ''; dpError = '';
		showForm = true;
	}

	async function addDatapoint() {
		if (!dpContent.trim()) { dpError = 'Content is required.'; return; }
		dpSaving = true; dpError = '';
		try {
			const dp = await api.datapoints.create(
				datasource!.id,
				dpTag,
				dpContent.trim(),
				dpDescription.trim()
			);
			datapoints = [...datapoints, dp];
			dpContent = ''; dpDescription = ''; showForm = false;
		} catch (e: unknown) {
			dpError = e instanceof Error ? e.message : 'Save failed.';
		} finally {
			dpSaving = false;
		}
	}

	async function move(id: string, direction: -1 | 1) {
		const sorted = [...datapoints].sort((a, b) => a.position - b.position);
		const idx = sorted.findIndex((d) => d.id === id);
		const other = sorted[idx + direction];
		if (!other) return;

		const a = sorted[idx];
		await Promise.all([
			api.datapoints.update(datasource!.id, a.id,     { tag: a.tag,     content: a.content,     description: a.description,     position: other.position }),
			api.datapoints.update(datasource!.id, other.id, { tag: other.tag, content: other.content, description: other.description, position: a.position     }),
		]);

		datapoints = datapoints.map((d) => {
			if (d.id === a.id)     return { ...d, position: other.position };
			if (d.id === other.id) return { ...d, position: a.position };
			return d;
		}).sort((a, b) => a.position - b.position);
	}

	async function remove(id: string) {
		if (!confirm('Delete this datapoint?')) return;
		try {
			await api.datapoints.delete(datasource!.id, id);
			datapoints = datapoints.filter((d) => d.id !== id);
		} catch (e: unknown) {
			alert(e instanceof Error ? e.message : 'Delete failed.');
		}
	}
</script>

<a href="/content" class="back">← Content</a>

<div class="section-header">
	<div>
		<h1>{datasource?.description ?? '…'}</h1>
		<span class="type-label">{datasource?.type}</span>
	</div>
</div>

<div class="dp-header">
	<h2>Datapoints <span class="count">{datapoints.length}</span></h2>
	{#if !showForm}
		<button class="primary" onclick={openForm}>+ Add Datapoint</button>
	{/if}
</div>

{#if showForm}
	<div class="card dp-form">
		<div class="form-row">
			<div class="field field-tag">
				<label for="dp-tag">Tag</label>
				<select id="dp-tag" bind:value={dpTag}>
					{#each TAGS as t}
						<option value={t}>{t}</option>
					{/each}
				</select>
			</div>
			<div class="field field-content">
				<label for="dp-content">Content</label>
				<textarea id="dp-content" bind:value={dpContent} rows="3" placeholder="Enter text…"></textarea>
			</div>
		</div>
		<div class="field">
			<label for="dp-desc">Description <span class="optional">(optional)</span></label>
			<input id="dp-desc" bind:value={dpDescription} placeholder="Internal note about this datapoint" />
		</div>
		{#if dpError}<p class="error-msg">{dpError}</p>{/if}
		<div class="form-actions">
			<button class="ghost" onclick={() => { showForm = false; dpError = ''; }}>Cancel</button>
			<button class="primary" onclick={addDatapoint} disabled={dpSaving}>
				{dpSaving ? 'Saving…' : 'Add'}
			</button>
		</div>
	</div>
{/if}

{#if datapoints.length === 0}
	<div class="empty">No datapoints yet. Add one above.</div>
{:else}
	<div class="dp-list">
		{#each [...datapoints].sort((a, b) => a.position - b.position) as dp, i (dp.id)}
			<div class="card dp-card">
				<div class="dp-meta">
					<span class="tag-badge tag-{dp.tag}">{dp.tag}</span>
					{#if dp.description}<span class="dp-desc">{dp.description}</span>{/if}
				</div>
				<p class="dp-content">{dp.content}</p>
				<div class="dp-actions">
					<div class="order-btns">
						<button class="ghost btn-sm" disabled={i === 0} onclick={() => move(dp.id, -1)}>↑</button>
						<button class="ghost btn-sm" disabled={i === datapoints.length - 1} onclick={() => move(dp.id, 1)}>↓</button>
					</div>
					<button class="danger btn-sm" onclick={() => remove(dp.id)}>Delete</button>
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
	.back:hover { color: #1a1a2e; }

	.section-header { margin-bottom: 28px; }

	h1 { font-size: 22px; font-weight: 700; }

	.type-label {
		font-family: monospace;
		font-size: 12px;
		color: #6b7280;
	}

	.dp-header {
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

	.card {
		background: #fff;
		border-radius: 10px;
		border: 1px solid #e5e7eb;
	}

	.dp-form {
		padding: 20px;
		margin-bottom: 16px;
	}

	.form-row {
		display: grid;
		grid-template-columns: 100px 1fr;
		gap: 12px;
		margin-bottom: 12px;
	}

	.form-actions {
		display: flex;
		gap: 8px;
		justify-content: flex-end;
		margin-top: 16px;
	}

	.optional {
		font-weight: 400;
		text-transform: none;
		letter-spacing: 0;
		color: #9ca3af;
	}

	.dp-list { display: flex; flex-direction: column; gap: 10px; }

	.dp-card { padding: 14px 16px; }

	.dp-meta {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-bottom: 8px;
	}

	.tag-badge {
		font-size: 11px;
		font-weight: 700;
		font-family: monospace;
		padding: 2px 8px;
		border-radius: 4px;
	}

	.tag-h1  { background: #ede9fe; color: #4f46e5; }
	.tag-h2  { background: #dbeafe; color: #1d4ed8; }
	.tag-p   { background: #f3f4f6; color: #374151; }
	.tag-li  { background: #d1fae5; color: #065f46; }

	.dp-desc { font-size: 12px; color: #9ca3af; }

	.dp-content {
		font-size: 14px;
		color: #1a1a2e;
		line-height: 1.5;
		margin-bottom: 10px;
	}

	.dp-actions {
		display: flex;
		justify-content: flex-end;
		align-items: center;
		gap: 8px;
	}

	.order-btns {
		display: flex;
		gap: 4px;
	}

	.order-btns button:disabled {
		opacity: 0.25;
		cursor: default;
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
