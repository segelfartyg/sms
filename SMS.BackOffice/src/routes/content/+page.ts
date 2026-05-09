import { api } from '$lib/api';

export const load = async () => {
	const datasources = await api.datasources.list();
	return { datasources };
};
