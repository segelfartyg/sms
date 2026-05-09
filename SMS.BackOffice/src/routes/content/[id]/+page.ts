import { api } from '$lib/api';

export const load = async ({ params }) => {
	const [datasources, datapoints] = await Promise.all([
		api.datasources.list(),
		api.datapoints.list(params.id)
	]);
	const datasource = datasources.find((d) => d.id === params.id);
	return { datasource, datapoints };
};
