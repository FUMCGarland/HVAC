<script>
	import { onMount } from 'svelte';
	import { invalidateAll } from '$app/navigation';
	import {
		Badge,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		A
	} from 'flowbite-svelte';

	export let data;
	let rooms = data.Rooms.filter((r) => r.ShellyID);

	// how do we sort the data for display
	const sortBy = { col: 'Name', ascending: true };

	// refresh every 30 seconds
	onMount(() => {
		const interval = setInterval(async () => {
			await invalidateAll();
			tablesort(sortBy.col);
			tablesort(sortBy.col); // twice to keep the ascending correct
		}, 30000);

		return () => {
			clearInterval(interval);
		};
	});

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
		rooms = rooms.sort(sortcallback);
	};

	function formatDate(d) {
		if (d == '0001-01-01T00:00:00Z') return '';
		return new Date(Date.parse(d)).toLocaleString();
	}
</script>

<Heading tag="h2">Rooms</Heading>
<Table>
	<TableHead>
		<TableHeadCell on:click={tablesort('Name')}>Room Name</TableHeadCell>
		<TableHeadCell on:click={tablesort('ShellyID')}>ShellyID</TableHeadCell>
		<TableHeadCell on:click={tablesort('LastUpdate')}>LastUpdate</TableHeadCell>
		<TableHeadCell on:click={tablesort('Battery')}>Battery</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each rooms as room}
			<TableBodyRow>
				<TableBodyCell><A href="/room/{room.ID}">{room.Name}</A></TableBodyCell>
				<TableBodyCell>{room.ShellyID}</TableBodyCell>
				<TableBodyCell>{formatDate(room.LastUpdate)}</TableBodyCell>

				{#if room.Battery > 45}
					<TableBodyCell><Badge color="green">{room.Battery}</Badge></TableBodyCell>
				{/if}
				{#if room.Battery <= 45 && room.Battery > 15}
					<TableBodyCell><Badge color="yellow">{room.Battery}</Badge></TableBodyCell>
				{/if}
				{#if room.Battery <= 15}
					<TableBodyCell><Badge color="red">{room.Battery}</Badge></TableBodyCell>
				{/if}
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
