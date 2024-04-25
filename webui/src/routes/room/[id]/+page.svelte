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
    console.log(data);
</script>

<Heading tag="h2">{data.Name}</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Zone</TableHeadCell>
		<TableHeadCell>Temp</TableHeadCell>
	</TableHead>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell><A href="/room/{data.ID}">{data.ID}</A></TableBodyCell>
			<TableBodyCell>{data.Name}</TableBodyCell>
			<TableBodyCell><A href="/zone/{data.Zone}">{data.Zone}</A></TableBodyCell>
			<TableBodyCell>{data.Temperature}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
