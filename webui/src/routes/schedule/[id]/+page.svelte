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
		Button
	} from 'flowbite-svelte';
	import { deleteSchedule } from '$lib/hvac.js';
	import { redirect } from '@sveltejs/kit';

	export let data;
	console.log(data);
	const weekdays = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
	const selectedwd = weekdays.map(() => false);

	function modeString(mode) {
		if (mode == 0) return 'heating';
		return 'cooling';
	}

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	// comes across as a byte string of UTF8...
	function parsePumps(p) {
		const e = new TextEncoder();
		const u = e.encode(atob(p));
		return u;
	}

	// comes across as a byte string of UTF8...
	function parseBlowers(b) {
		let e = new TextEncoder();
		let u = e.encode(atob(b));
		return u;
	}
</script>

<Heading tag="h2">{data.Name}</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Value</TableHeadCell>
	</TableHead>
	<TableBody>
		<TableBodyRow>
			<TableBodyCell>ID</TableBodyCell>
			<TableBodyCell>{data.ID}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Name</TableBodyCell>
			<TableBodyCell>{data.Name}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Mode</TableBodyCell>
			<TableBodyCell>{modeString(data.Mode)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Start Times</TableBodyCell>
			<TableBodyCell>{data.StartTime}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Weekdays</TableBodyCell>
			<TableBodyCell>{parseWeekdays(data.Weekdays)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Pumps</TableBodyCell>
			<TableBodyCell>{parsePumps(data.Pumps)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Blowers</TableBodyCell>
			<TableBodyCell>{parseBlowers(data.Blowers)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell
				><Button
					on:click={() => {
						deleteSchedule(data.ID);
						throw redirect(303, '/schedule');
					}}>Delete</Button
				></TableBodyCell
			>
			<TableBodyCell></TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
