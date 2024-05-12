<script>
	import { goto } from '$app/navigation';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		Button,
		A
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	const weekdays = ['Sun', 'M', 'T', 'W', 'Th', 'F', 'Sat'];
	export let data;

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	function formatDate(d) {
		return new Date(Date.parse(d)).toLocaleString();
	}
</script>

<Heading tag="h2">Recurring Occupancy Events</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Rooms</TableHeadCell>
		<TableHeadCell>Weekdays</TableHeadCell>
		<TableHeadCell>Start Time</TableHeadCell>
		<TableHeadCell>End Time</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Recurring as r}
			<TableBodyRow>
				<TableBodyCell><A href="/occupancy/recurring/{r.ID}">{r.Name}</A></TableBodyCell>
				<TableBodyCell>{r.Rooms}</TableBodyCell>
				<TableBodyCell>{parseWeekdays(r.Weekdays)}</TableBodyCell>
				<TableBodyCell>{r.StartTime}</TableBodyCell>
				<TableBodyCell>{r.EndTime}</TableBodyCell>
			</TableBodyRow>
		{/each}
		<TableBodyRow>
			<TableBodyCell colspan="4">&nbsp;</TableBodyCell>
			<TableBodyCell
				><Button on:click={() => goto('/occupancy/recurring')}>Add Recurring</Button></TableBodyCell
			>
		</TableBodyRow>
	</TableBody>
</Table>

<Heading tag="h2">One-Time Occupancy Events</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Rooms</TableHeadCell>
		<TableHeadCell>Start</TableHeadCell>
		<TableHeadCell>End</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.OneTime as o}
			<TableBodyRow>
				<TableBodyCell><A href="/occupancy/onetime/{o.ID}">{o.Name}</A></TableBodyCell>
				<TableBodyCell>{o.Rooms}</TableBodyCell>
				<TableBodyCell>{formatDate(o.Start)}</TableBodyCell>
				<TableBodyCell>{formatDate(o.End)}</TableBodyCell>
			</TableBodyRow>
		{/each}
		<TableBodyRow>
			<TableBodyCell colspan="3">&nbsp;</TableBodyCell>
			<TableBodyCell
				><Button on:click={() => goto('/occupancy/onetime')}>Add One-Time</Button></TableBodyCell
			>
		</TableBodyRow>
	</TableBody>
</Table>
