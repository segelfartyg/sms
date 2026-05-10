export interface Page {
	id: string;
	title: string;
	slug: string;
	boxes?: Box[];
	created_at: string;
	updated_at: string;
}

export interface Datasource {
	id: string;
	type: string;
	description: string;
	content: string;
}

export interface Datapoint {
	id: string;
	datasource_id: string;
	tag: string;
	content: string;
	description: string;
	position: number;
}

export interface Box {
	id: string;
	page_id: string;
	type: string;
	content: Record<string, unknown>;
	position: number;
	datasource_id?: string;
	created_at: string;
	updated_at: string;
}

const BASE      = import.meta.env.VITE_API_URL       ?? 'http://localhost:8080';
const WAREHOUSE = import.meta.env.VITE_WAREHOUSE_URL ?? 'http://localhost:8081';

async function req<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${BASE}${path}`, {
		headers: { 'Content-Type': 'application/json', ...init?.headers },
		...init
	});
	if (!res.ok) {
		const msg = await res.text().catch(() => res.statusText);
		throw new Error(msg);
	}
	if (res.status === 204) return undefined as T;
	return res.json();
}

async function wreq<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${WAREHOUSE}${path}`, {
		headers: { 'Content-Type': 'application/json', ...init?.headers },
		...init
	});
	if (!res.ok) {
		const msg = await res.text().catch(() => res.statusText);
		throw new Error(msg);
	}
	if (res.status === 204) return undefined as T;
	return res.json();
}

export const api = {
	pages: {
		list: () => req<Page[]>('/pages'),
		get: (id: string) => req<Page>(`/pages/${id}`),
		create: (title: string, slug: string) =>
			req<Page>('/pages', { method: 'POST', body: JSON.stringify({ title, slug }) }),
		update: (id: string, title: string, slug: string) =>
			req<Page>(`/pages/${id}`, { method: 'PUT', body: JSON.stringify({ title, slug }) }),
		delete: (id: string) => req<void>(`/pages/${id}`, { method: 'DELETE' })
	},
	boxes: {
		create: (
			pageId: string,
			type: string,
			content: Record<string, unknown>,
			position: number,
			datasourceId?: string
		) =>
			req<Box>(`/pages/${pageId}/boxes`, {
				method: 'POST',
				body: JSON.stringify({ type, content, position, datasource_id: datasourceId ?? null })
			}),
		update: (
			pageId: string,
			id: string,
			type: string,
			content: Record<string, unknown>,
			position: number,
			datasourceId?: string
		) =>
			req<Box>(`/pages/${pageId}/boxes/${id}`, {
				method: 'PUT',
				body: JSON.stringify({ type, content, position, datasource_id: datasourceId ?? null })
			}),
		delete: (pageId: string, id: string) =>
			req<void>(`/pages/${pageId}/boxes/${id}`, { method: 'DELETE' })
	},
	datasources: {
		list: () => req<Datasource[]>('/datasources'),
		create: (description: string, type: string) =>
			wreq<Datasource>('/datasources', {
				method: 'POST',
				body: JSON.stringify({ description, type, content: '' })
			}),
		delete: (id: string) => wreq<void>(`/datasources/${id}`, { method: 'DELETE' })
	},
	datapoints: {
		list: (datasourceId: string) =>
			wreq<Datapoint[]>(`/datasources/${datasourceId}/datapoints`),
		create: (datasourceId: string, tag: string, content: string, description: string) =>
			wreq<Datapoint>(`/datasources/${datasourceId}/datapoints`, {
				method: 'POST',
				body: JSON.stringify({ tag, content, description })
			}),
		update: (datasourceId: string, id: string, dp: Pick<Datapoint, 'tag' | 'content' | 'description' | 'position'>) =>
			wreq<Datapoint>(`/datasources/${datasourceId}/datapoints/${id}`, {
				method: 'PUT',
				body: JSON.stringify(dp)
			}),
		delete: (datasourceId: string, id: string) =>
			wreq<void>(`/datasources/${datasourceId}/datapoints/${id}`, { method: 'DELETE' })
	}
};
