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
	import { Heading, P, A } from 'flowbite-svelte';

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

<Heading tag="h2">Rooms</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Zone</TableHeadCell>
		<TableHeadCell>Temperature</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Rooms as room}
			<TableBodyRow>
				<TableBodyCell><A href="/room/{room.ID}">{room.ID}</A></TableBodyCell>
				<TableBodyCell>{room.Name}</TableBodyCell>
				<TableBodyCell><A href="/zone/{room.Zone}">{room.Zone}</A></TableBodyCell>
				<TableBodyCell>{room.Temperature}</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
