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
</script>

<Heading tag="h2">{data.Loop.Name}</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Pump ID</TableHeadCell>
		<TableHeadCell>Pump Name</TableHeadCell>
	</TableHead>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell><A href="/loop/{data.Loop.ID}">{data.Loop.ID}</A></TableBodyCell>
			<TableBodyCell>{data.Loop.Name}</TableBodyCell>
			<TableBodyCell><A href="/pump/{data.Pump.ID}">{data.Pump.ID}</A></TableBodyCell>
			<TableBodyCell>{data.Pump.Name}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
<Hr />
<Heading tag="h2">Blowers</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Zone</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Blowers as blower}
			<TableBodyRow>
				<TableBodyCell><A href="/blower/{blower.ID}">{blower.ID}</A></TableBodyCell>
				<TableBodyCell>{blower.Name}</TableBodyCell>
				<TableBodyCell><A href="/zone/{blower.Zone}">{blower.Zone}</A></TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
