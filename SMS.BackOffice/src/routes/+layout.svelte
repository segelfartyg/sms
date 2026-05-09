<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="shell">
	<header>
		<a href="/pages" class="brand">KMS BackOffice</a>
		<nav>
			<a href="/pages">Pages</a>
			<a href="/content">Content</a>
		</nav>
	</header>
	<main>
		{@render children()}
	</main>
</div>

<style>
	:global(*, *::before, *::after) {
		box-sizing: border-box;
		margin: 0;
		padding: 0;
	}

	:global(body) {
		font-family: system-ui, -apple-system, sans-serif;
		background: #f0f2f5;
		color: #1a1a2e;
		font-size: 14px;
		line-height: 1.5;
	}

	:global(button) {
		cursor: pointer;
		font: inherit;
		border: none;
		border-radius: 6px;
		padding: 6px 14px;
		font-size: 13px;
		transition: opacity 0.15s;
	}

	:global(button:hover) {
		opacity: 0.85;
	}

	:global(button.primary) {
		background: #4f46e5;
		color: #fff;
	}

	:global(button.danger) {
		background: #ef4444;
		color: #fff;
	}

	:global(button.ghost) {
		background: transparent;
		border: 1px solid #d1d5db;
		color: #374151;
	}

	:global(input, textarea, select) {
		font: inherit;
		font-size: 13px;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		padding: 7px 10px;
		width: 100%;
		background: #fff;
		transition: border-color 0.15s;
	}

	:global(input:focus, textarea:focus, select:focus) {
		outline: none;
		border-color: #4f46e5;
	}

	:global(label) {
		display: block;
		font-size: 12px;
		font-weight: 600;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		margin-bottom: 4px;
	}

	:global(.field) {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	:global(.error-msg) {
		color: #ef4444;
		font-size: 12px;
		margin-top: 6px;
	}

	.shell {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	header {
		background: #1a1a2e;
		padding: 0 24px;
		height: 52px;
		display: flex;
		align-items: center;
		gap: 32px;
		position: sticky;
		top: 0;
		z-index: 10;
	}

	.brand {
		color: #fff;
		text-decoration: none;
		font-weight: 700;
		font-size: 15px;
		letter-spacing: -0.01em;
	}

	nav {
		display: flex;
		gap: 4px;
	}

	nav a {
		color: rgba(255, 255, 255, 0.6);
		text-decoration: none;
		font-size: 13px;
		padding: 4px 10px;
		border-radius: 5px;
		transition: color 0.15s, background 0.15s;
	}

	nav a:hover {
		color: #fff;
		background: rgba(255, 255, 255, 0.1);
	}

	main {
		flex: 1;
		padding: 32px 24px;
		max-width: 900px;
		margin: 0 auto;
		width: 100%;
	}
</style>
