<script>
	import { goto } from '$app/navigation';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		Button,
		A,
		Banner
	} from 'flowbite-svelte';

	export var data;
	const sortBy = { col: 'Name', ascending: true };

	$: tablesort = (column) => {
		if (sortBy.col == column) {
			sortBy.ascending = !sortBy.ascending;
		} else {
			sortBy.col = column;
			sortBy.ascending = true;
		}

		let sm = sortBy.ascending ? 1 : -1;

		let sortcallback = (a, b) =>
			a[column] < b[column] ? -1 * sm : a[column] > b[column] ? 1 * sm : 0;
		data.data = data.data.sort(sortcallback);
	};

	function formatDate(d) {
		let nd = new Date(Date.parse(d));
		nd.setMilliseconds(0);
		nd.setSeconds(0);
		return nd.toLocaleString();
	}
</script>

<Heading tag="h2">Job Next Run Schedule</Heading>
<Table>
	<TableHead>
		<TableHeadCell on:click={tablesort('Name')}>Name</TableHeadCell>
		<TableHeadCell on:click={tablesort('NextRun')}>Next Run (pretty)</TableHeadCell>
		<TableHeadCell on:click={tablesort('NextRun')}>Next Run (raw)</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.data as r}
			<TableBodyRow>
				<TableBodyCell>{r.Name}</TableBodyCell>
				<TableBodyCell>{formatDate(r.NextRun)}</TableBodyCell>
				<TableBodyCell>{r.NextRun}</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
