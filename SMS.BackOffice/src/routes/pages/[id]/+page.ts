import { api } from '$lib/api';

export const load = async ({ params }) => {
	const page = await api.pages.get(params.id);
	return { page };
};
