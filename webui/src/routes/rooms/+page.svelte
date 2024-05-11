<script>
	import {
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

	function zoneName(zoneID) {
		const z = data.Zones.filter((z) => z.ID == zoneID);
		return z[0].Name;
	}
</script>

<Heading tag="h2">Rooms</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Zone</TableHeadCell>
		<TableHeadCell>Temperature</TableHeadCell>
		<TableHeadCell>Humidity</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Rooms as room}
			<TableBodyRow>
				<TableBodyCell><A href="/room/{room.ID}">{room.Name}</A></TableBodyCell>
				<TableBodyCell><A href="/zone/{room.Zone}">{zoneName(room.Zone)}</A></TableBodyCell>
				<TableBodyCell>{room.Temperature}</TableBodyCell>
				<TableBodyCell>{room.Humidity}</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
