import { api } from '$lib/api';

export const load = async () => {
	const pages = await api.pages.list();
	return { pages };
};
