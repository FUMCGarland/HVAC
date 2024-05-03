<script>
	import {
		Table,
		TableBodyRow,
		TableBody,
		TableBodyCell,
		Heading,
		P,
		A,
		Hr
	} from 'flowbite-svelte';

	export let data;
	const e = new TextEncoder();
	const loops = e.encode(atob(data.Loops));
</script>

<Heading tag="h2">Chiller {data.ID}: {data.Name}</Heading>
<Table>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell>ID</TableBodyCell>
			<TableBodyCell><A href="/chiller/{data.ID}">{data.ID}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Name</TableBodyCell>
			<TableBodyCell>{data.Name}</TableBodyCell>
		</TableBodyRow>
		{#each loops as l}
			<TableBodyRow>
				<TableBodyCell>Loop</TableBodyCell>
				<TableBodyCell><A href="/loop/{l}">{l}</A></TableBodyCell>
			</TableBodyRow>
		{/each}
		<TableBodyRow>
			<TableBodyCell>Running</TableBodyCell>
			<TableBodyCell>{data.Running}</TableBodyCell>
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
