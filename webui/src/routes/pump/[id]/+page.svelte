<script>
	import {
		Badge,
		Button,
		Table,
		TableBodyRow,
		TableBody,
		TableBodyCell,
		Heading,
		P,
		A,
		Hr
	} from 'flowbite-svelte';
	import { pumpStop } from '$lib/hvac';

	export let data;

	function sm(s) {
		if (sm == 0) return 'heating';
		return 'cooling';
	}
</script>

<Heading tag="h2">Pump {data.ID}: {data.Name}</Heading>
<Table>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell>ID</TableBodyCell>
			<TableBodyCell><A href="/pump/{data.ID}">{data.ID}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Name</TableBodyCell>
			<TableBodyCell>{data.Name}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Loop</TableBodyCell>
			<TableBodyCell><A href="/loop/{data.Loop}">{data.Loop}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>SystemMode</TableBodyCell>
			<TableBodyCell>{sm(data.SystemMode)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Running</TableBodyCell>
			{#if data.Running}
				<TableBodyCell
					><Button
						on:click={() => {
							pumpStop(data.ID);
						}}>Stop</Button
					></TableBodyCell
				>
			{/if}
			{#if !data.Running}
				<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
			{/if}
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Current Start Time</TableBodyCell>
			<TableBodyCell>{data.CurrentStartTime}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Last Start Time</TableBodyCell>
			<TableBodyCell>{data.LastStartTime}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Last Stop Time</TableBodyCell>
			<TableBodyCell>{data.LastStopTime}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Total Run Time</TableBodyCell>
			<TableBodyCell>{data.Runtime}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
