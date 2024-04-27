<script>
	import {
		Table,
		TableBody,
		TableBodyRow,
		TableBodyCell,
		Heading,
		P,
		A,
		Hr
	} from 'flowbite-svelte';

	export let data;
	function parseTime(t) {
		const d = new Date(t);
		return d.toString();
	}
	function parseDuration(d) {
		let duration = { days: 0, hours: 0, minutes: 0, seconds: 0 };
		const days = Math.floor(d / 86400);
		if (days >= 1) {
			duration.days = days;
			d -= days * 86400;
		}
		const hours = Math.floor(d / 3600);
		if (hours >= 1) {
			duration.hours = hours;
			d -= hours * 3600;
		}
		const minutes = Math.floor(d / 60);
		if (minutes >= 1) {
			duration.minutes = minutes;
			d -= minutes * 60;
		}
		duration.seconds = d;
		// browsers don't support Intl.DurationFormat() yet
		return `${duration.days} days; ${duration.hours}:${duration.minutes}:${duration.seconds}`;
	}
</script>

<Heading tag="h2">Blower {data.ID}: {data.Name}</Heading>
<Table>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell>ID</TableBodyCell>
			<TableBodyCell><A href="/blower/{data.ID}">{data.ID}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Name</TableBodyCell>
			<TableBodyCell>{data.Name}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Hot Loop</TableBodyCell>
			<TableBodyCell><A href="/loop/{data.HotLoop}">{data.HotLoop}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Cold Loop</TableBodyCell>
			<TableBodyCell><A href="/loop/{data.ColdLoop}">{data.ColdLoop}</A></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Running</TableBodyCell>
			<TableBodyCell>{data.Running}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Current Start Time</TableBodyCell>
			<TableBodyCell>{parseTime(data.CurrentStartTime)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Last Start Time</TableBodyCell>
			<TableBodyCell>{parseTime(data.LastStartTime)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Last Stop Time</TableBodyCell>
			<TableBodyCell>{parseTime(data.LastStopTime)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Total Run Time</TableBodyCell>
			<TableBodyCell>{parseDuration(data.Runtime)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Filter Time</TableBodyCell>
			<TableBodyCell>{parseDuration(data.FilterTime)}</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
