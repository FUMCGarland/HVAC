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
		<TableBodyRow>
			<TableBodyCell>Humidity</TableBodyCell>
			<TableBodyCell>{data.Humidity}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Occupied</TableBodyCell>
			<TableBodyCell>{data.Occupied}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Last Update</TableBodyCell>
			<TableBodyCell>{data.LastUpdate}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Battery</TableBodyCell>
			<TableBodyCell>{data.Battery}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>ShellyID (sensor)</TableBodyCell>
			<TableBodyCell><A href="/shelly/">{data.ShellyID}</A></TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
