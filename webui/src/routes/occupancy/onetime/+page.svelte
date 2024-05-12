<script>
	import { goto } from '$app/navigation';
	import Flatpickr from 'svelte-flatpickr';
	import 'flatpickr/dist/flatpickr.css';
	import { onMount } from 'svelte';
	import { postRecurringOccupancy } from '$lib/occupancy';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Input,
		Button,
		Dropdown,
		DropdownItem,
		Checkbox,
		Radio,
		Heading,
		Hr,
		A
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	const weekdays = ['Sun', 'M', 'T', 'W', 'Th', 'F', 'Sat'];
	const selectedwd = weekdays.map(() => false);
	const options = {
		enableTime: true
	};

	export let data;
	data.Rooms.forEach((r) => {
		r.selected = false;
	});

	let id = data.Recurring.length + 1; // get highest number and increment
	let name = 'not set';
	let start;
	let end;

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	function parseRooms(r) {
		const u = new TextEncoder().encode(atob(r));
		return data.Rooms.filter((r) => u.includes(r.ID)).map((s) => s.Name);
	}

	async function doAddOneTime() {
		let c = {
			ID: Number(id),
			Name: name,
			Start: start,
			End: end,
			Rooms: data.Rooms.filter((r) => r.selected).map((r) => z.ID)
		};
		// await postOccupancyOneTime(c);
		// goto('/occupancy');
		console.log(c);
	}
</script>

<form>
	<Table>
		<TableHead>
			<TableHeadCell>Name</TableHeadCell>
			<TableHeadCell>Value</TableHeadCell>
		</TableHead>
		<TableBody>
			<TableBodyRow>
				<TableBodyCell>ID</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={id} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Name</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={name} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Start</TableBodyCell>
				<TableBodyCell>
					<Flatpickr {options} bind:start name="startdate" />
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>End</TableBodyCell>
				<TableBodyCell>
					<Flatpickr {options} bind:end name="enddate" />
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Rooms</TableBodyCell>
				<TableBodyCell>
					{#each data.Rooms as r}
						<Checkbox value={r.ID} bind:checked={r.selected}>{r.Name}</Checkbox>
					{/each}
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>&nbsp;</TableBodyCell>
				<TableBodyCell><Button on:click={() => doAddOneTime()}>Add</Button></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
