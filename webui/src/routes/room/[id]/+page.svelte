<script>
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		P,
		A,
		Hr
	} from 'flowbite-svelte';

	export let data;
</script>

<Heading tag="h2">{data.Name}</Heading>
<Table>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell>Room</TableBodyCell>
			<TableBodyCell><A href="/room/{data.ID}">{data.Name}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Zone</TableBodyCell>
			<TableBodyCell><A href="/zone/{data.Zone}">{data.Zone}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Temp</TableBodyCell>
			<TableBodyCell>{data.Temperature}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
