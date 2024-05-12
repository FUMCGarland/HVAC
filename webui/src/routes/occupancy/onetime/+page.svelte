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
	export let data;
	data.Rooms.forEach((r) => {
		r.selected = false;
	});

	let id = data.OneTime.length + 1; // TODO: find lowest unused
	let name = 'not set';
	let start = '2000-01-01 11:00';
	let end = '2000-01-01 13:00';

	const optionsStart = {
		enableTime: true,
		minDate: 'today',
		onChange(selectedDates, dateStr) {
			start = dateStr;
			optionsEnd.enable = [dateStr.split(' ')[0]];
		}
	};

	const optionsEnd = {
		enableTime: true,
		enable: [],
		onChange(selectedDates, dateStr) {
			end = dateStr;
		}
	};

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
			Rooms: data.Rooms.filter((r) => r.selected).map((r) => r.ID)
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
					<Flatpickr options={optionsStart} name="startdate" />
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>End</TableBodyCell>
				<TableBodyCell>
					<Flatpickr options={optionsEnd} name="enddate" />
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
