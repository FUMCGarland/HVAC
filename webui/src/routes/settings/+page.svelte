<script>
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		A,
	} from 'flowbite-svelte';

	export let data;
</script>

<Heading tag="h2">System Settings</Heading>
<Table>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell>DataLogFile</TableBodyCell>
			<TableBodyCell>{data.DataLogFile}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>HTTPAuthData (file)</TableBodyCell>
			<TableBodyCell>{data.HTTPAuthData}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>HTTPStaticDir</TableBodyCell>
			<TableBodyCell>{data.HTTPStaticDir}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>HTTPaddr</TableBodyCell>
			<TableBodyCell>{data.HTTPaddr}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>MQTT.Auth (file)</TableBodyCell>
			<TableBodyCell>{data.MQTT.Auth}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>MQTT.ID</TableBodyCell>
			<TableBodyCell>{data.MQTT.ID}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>MQTT.ListenAddr</TableBodyCell>
			<TableBodyCell>{data.MQTT.ListenAddr}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>MQTT.Root</TableBodyCell>
			<TableBodyCell>{data.MQTT.Root}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>OpenWeatherMapID</TableBodyCell>
			<TableBodyCell>{data.OpenWeatherMapID}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>OpenWeatherMapKey</TableBodyCell>
			<TableBodyCell>{data.OpenWeatherMapKey}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>StateStore</TableBodyCell>
			<TableBodyCell>{data.StateStore}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
