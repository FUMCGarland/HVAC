<script>
	import { onMount } from 'svelte';
	import { invalidateAll } from '$app/navigation';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import { Heading, P, A, Hr } from 'flowbite-svelte';

	export let data;

	onMount(() => {
		const interval = setInterval(() => {
			invalidateAll();
		}, 30000);

		return () => {
			clearInterval(interval);
		};
	});
</script>

<Heading tag="h2">Blower {data.ID}: {data.Name}</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Hot Loop</TableHeadCell>
		<TableHeadCell>Cold Loop</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
		<TableHeadCell>Current Start Time</TableHeadCell>
		<TableHeadCell>Last Start Time</TableHeadCell>
		<TableHeadCell>Last Stop Time</TableHeadCell>
		<TableHeadCell>Total Run Time</TableHeadCell>
		<TableHeadCell>Filter Time</TableHeadCell>
	</TableHead>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell><A href="/blower/{data.ID}">{data.ID}</A></TableBodyCell>
			<TableBodyCell>{data.Name}</TableBodyCell>
			<TableBodyCell><A href="/loop/{data.HotLoop}">{data.HotLoop}</A></TableBodyCell>
			<TableBodyCell><A href="/loop/{data.ColdLoop}">{data.ColdLoop}</A></TableBodyCell>
			<TableBodyCell>{data.Running}</TableBodyCell>
			<TableBodyCell>{data.CurrentStartTime}</TableBodyCell>
			<TableBodyCell>{data.LastStartTime}</TableBodyCell>
			<TableBodyCell>{data.LastStopTime}</TableBodyCell>
			<TableBodyCell>{data.Runtime}</TableBodyCell>
			<TableBodyCell>{data.FilterTime}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
