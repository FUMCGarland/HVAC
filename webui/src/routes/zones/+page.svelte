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
</script>

<Heading tag="h2">Zones</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Heating Unoccupied Temp</TableHeadCell>
		<TableHeadCell>Heating Occupied Temp</TableHeadCell>
		<TableHeadCell>Cooling Unoccupied Temp</TableHeadCell>
		<TableHeadCell>Cooling Occupied</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Zones as zone}
			<TableBodyRow>
				<TableBodyCell><A href="/zone/{zone.ID}">{zone.Name}</A></TableBodyCell>
				<TableBodyCell>{zone.Targets.HeatingUnoccupiedTemp}</TableBodyCell>
				<TableBodyCell>{zone.Targets.HeatingOccupiedTemp}</TableBodyCell>
				<TableBodyCell>{zone.Targets.CoolingUnoccupiedTemp}</TableBodyCell>
				<TableBodyCell>{zone.Targets.CoolingOccupiedTemp}</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
