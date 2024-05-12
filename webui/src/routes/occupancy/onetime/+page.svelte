<script>
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	const weekdays = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
	const selectedwd = weekdays.map(() => false);
	export let data;

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	async function doAddRecurring() {
		let c = {
			ID: Number(id),
			Name: name,
			Weekdays: selectedwd
				.map((n, i) => {
					if (n) {
						return i;
					}
				})
				.filter((o) => {
					return o !== undefined;
				}),
			Starttime: starttime,
			Runtime: Number(runtime),
			Rooms: data.Rooms.filter((r) => r.selected).map((r) => z.ID)
		};
		await postOccupancyRecurring(c);
	}

	async function doAddOneTime() {
		let c = {
			ID: Number(id),
			Name: name,
			Start: start,
			End: end,
			Rooms: data.Rooms.filter((r) => r.selected).map((r) => z.ID)
		};
		await postOccupancyOneTime(c);
	}
</script>

<Heading tag="h2">Recurring Occupancy Events</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Rooms</TableHeadCell>
		<TableHeadCell>Weekdays</TableHeadCell>
		<TableHeadCell>Start Time</TableHeadCell>
		<TableHeadCell>Duration</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Recurring as r}
			<TableBodyRow>
				<TableBodyCell>{r.Name}</TableBodyCell>
				<TableBodyCell>{r.Rooms}</TableBodyCell>
				<TableBodyCell>{r.Weekdays}</TableBodyCell>
				<TableBodyCell>{r.StartTime}</TableBodyCell>
				<TableBodyCell>{r.Duration}</TableBodyCell>
			</TableBodyRow>
		{/each}
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
		{#each data.Recurring as r}
			<TableBodyRow>
				<TableBodyCell>{r.Name}</TableBodyCell>
				<TableBodyCell>{r.Rooms}</TableBodyCell>
				<TableBodyCell>{r.Start}</TableBodyCell>
				<TableBodyCell>{r.End}</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
