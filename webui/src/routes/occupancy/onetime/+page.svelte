<script>
	import { goto } from '$app/navigation';
	// import Flatpickr from 'svelte-flatpickr';
	// import 'flatpickr/dist/flatpickr.css';
	import svlatepickr from 'svelte-flatpickr-plus';
	import { postOneTimeOccupancy } from '$lib/occupancy';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Input,
		Button,
		Checkbox
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import { lowestFreeID } from '$lib/util';

	export let data;

	data.Rooms.forEach((r) => {
		r.selected = false;
	});

	let id = lowestFreeID(data.OneTime);
	let name = 'not set';
	let start = '2000-01-01 11:00'; // 2006-01-02T15:04:05.999999999 -0700
	let end = '2000-01-01 13:00'; // 2006-01-02T15:04:05.999999999 -0700

	const optionsStart = {
		enableTime: true,
		enableSeconds: true,
		minDate: 'today',
		dateFormat: 'Z',
		onChange(selectedDates, dateStr) {
			start = dateStr;
			optionsEnd.enable = [dateStr.split(' ')[0]];
		}
	};

	const optionsEnd = {
		enableTime: true,
		enableSeconds: true,
		dateFormat: 'Z',
		enable: [],
		onChange(selectedDates, dateStr) {
			end = dateStr;
		}
	};

	async function doAddOneTime() {
		let c = {
			ID: Number(id),
			Name: name,
			Start: start,
			End: end,
			Rooms: data.Rooms.filter((r) => r.selected).map((r) => r.ID)
		};
		console.log(c);
		try {
			await postOneTimeOccupancy(c);
			goto('/occupancy');
		} catch (e) {
			console.log(e);
		}
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
					<!-- <Flatpickr options={optionsStart} name="startdate" /> -->
					<input name="startdate" use:svlatepickr />
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>End</TableBodyCell>
				<TableBodyCell>
					<!-- <Flatpickr options={optionsEnd} name="enddate" /> -->
					<input name="enddate" use:svlatepickr />
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
