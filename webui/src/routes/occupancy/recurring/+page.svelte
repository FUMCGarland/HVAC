<script>
	import { goto } from '$app/navigation';
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

	let id = data.Recurring.length + 1; // get highest number and increment
	let name = 'not set';
	$: mode = 0;
	let starttime = '09:30';
	let endtime = '16:30';

	function modeString(mode) {
		if (mode == 0) return 'heating';
		return 'cooling';
	}

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	function parseRooms(r) {
		const u = new TextEncoder().encode(atob(r));
		return data.Rooms.filter((r) => u.includes(r.ID)).map((s) => s.Name);
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
			StartTime: starttime,
			EndTime: endtime,
			Rooms: data.Rooms.filter((r) => r.selected).map((r) => r.ID)
		};
		await postRecurringOccupancy(c);
		goto('/occupancy');
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
				<TableBodyCell>Weekdays</TableBodyCell>
				<TableBodyCell>
					<Button
						>Weekdays<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" /></Button
					>
					<Dropdown class="w-44 space-y-3 p-3 text-sm">
						{#each weekdays as wd, index}
							<li><Checkbox value={wd} bind:checked={selectedwd[index]}>{wd}</Checkbox></li>
						{/each}
					</Dropdown>
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Start Time (24-hour hh:mm format)</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={starttime} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>End Time (24-hour hh:mm format)</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={endtime} /></TableBodyCell>
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
				<TableBodyCell><Button on:click={() => doAddRecurring()}>Add</Button></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
