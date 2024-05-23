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
		A,
		Input,
		Button,
		Label
	} from 'flowbite-svelte';
	import { updateZoneTargets } from '$lib/hvac.js';

	export let data;
	let hu = data.Targets.HeatingUnoccupiedTemp;
	let ho = data.Targets.HeatingOccupiedTemp;
	let cu = data.Targets.CoolingUnoccupiedTemp;
	let co = data.Targets.CoolingOccupiedTemp;

	function update() {
		const cmd = JSON.stringify({
			HeatingUnoccupiedTemp: Number(hu),
			HeatingOccupiedTemp: Number(ho),
			CoolingUnoccupiedTemp: Number(cu),
			CoolingOccupiedTemp: Number(co)
		});
		updateZoneTargets(data.ID, cmd);
	}
</script>

<Heading tag="h2">Zone {data.ID}: {data.Name}</Heading>
<form>
	<Table>
		<TableHead>
			<TableHeadCell>Heating Unoccupied</TableHeadCell>
			<TableHeadCell>Heating Occupied</TableHeadCell>
			<TableHeadCell>Cooling Unoccupied</TableHeadCell>
			<TableHeadCell>Cooling Occupied</TableHeadCell>
		</TableHead>
		<TableBody>
			<TableBodyRow>
				<TableBodyCell>{data.Targets.HeatingUnoccupiedTemp}</TableBodyCell>
				<TableBodyCell>{data.Targets.HeatingOccupiedTemp}</TableBodyCell>
				<TableBodyCell>{data.Targets.CoolingUnoccupiedTemp}</TableBodyCell>
				<TableBodyCell>{data.Targets.CoolingOccupiedTemp}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell><Input type="text" bind:value={hu} id="HU" /></TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={ho} id="HO" /></TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={cu} id="CU" /></TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={co} id="CO" /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell colspan="3">&nbsp;</TableBodyCell>
				<TableBodyCell
					><Button
						on:click={() => {
							update();
						}}>Update</Button
					></TableBodyCell
				>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>

<Heading tag="h3">Rooms in Zone {data.ID}: {data.Name}</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Current Temperature</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Rooms as room}
			<TableBodyRow>
				<TableBodyCell><A href="/room/{room.ID}">{room.Name}</A></TableBodyCell>
				<TableBodyCell>{room.Temperature}</TableBodyCell>
			</TableBodyRow>
		{/each}
		<TableBodyRow>
			<TableBodyCell>Zone Average</TableBodyCell>
			<TableBodyCell>{data.AverageTemp}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>

<Heading tag="h3">Blowers for Zone {data.ID}: {data.Name}</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Hot Loop</TableHeadCell>
		<TableHeadCell>Cold Loop</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Blowers as blower}
			<TableBodyRow>
				<TableBodyCell><A href="/blower/{blower.ID}">{blower.Name}</A></TableBodyCell>
				<TableBodyCell><A href="/loop/{blower.HotLoop}">{blower.HotLoop}</A></TableBodyCell>
				<TableBodyCell><A href="/loop/{blower.ColdLoop}">{blower.ColdLoop}</A></TableBodyCell>
				{#if blower.Running}
					<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
				{/if}
				{#if !blower.Running}
					<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
				{/if}
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
