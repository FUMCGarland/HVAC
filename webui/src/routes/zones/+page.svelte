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

<Heading tag="h2">Zones</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Heating Unoccupied Temp</TableHeadCell>
		<TableHeadCell>Heating Occupied Temp</TableHeadCell>
		<TableHeadCell>Cooling Unoccupied Temp</TableHeadCell>
		<TableHeadCell>Cooling Occupied</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Zones as zone}
			<TableBodyRow>
				<TableBodyCell><A href="/zone/{zone.ID}">{zone.ID}</A></TableBodyCell>
				<TableBodyCell>{zone.Name}</TableBodyCell>
				<TableBodyCell>{zone.Targets.HeatingUnoccupiedTemp}</TableBodyCell>
				<TableBodyCell>{zone.Targets.HeatingOccupiedTemp}</TableBodyCell>
				<TableBodyCell>{zone.Targets.CoolingUnoccupiedTemp}</TableBodyCell>
				<TableBodyCell>{zone.Targets.CoolingOccupiedTemp}</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
