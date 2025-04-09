<script>
	import { goto } from '$app/navigation';
	import { toast } from '@zerodevx/svelte-toast';
	import { onMount } from 'svelte';
	import { putRecurringOccupancy, deleteRecurringOccupancy } from '$lib/occupancy';
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
	import { toLocaleTimestring, toZTimestring } from '$lib/util.js';

	const weekdays = ['Sun', 'M', 'T', 'W', 'Th', 'F', 'Sat'];
	const selectedwd = weekdays.map(() => false);

	export let data;
	data.StartTime = toLocaleTimestring(data.StartTime);
	data.EndTime = toLocaleTimestring(data.EndTime);

	data.SystemRooms.forEach((r) => {
		r.selected = false;
		if (data.Rooms.includes(r.ID)) r.selected = true;
	});
	data.Weekdays.forEach((d) => {
		selectedwd[d] = true;
	});

	function parseRooms(r) {
		const u = new TextEncoder().encode(atob(r));
		return data.Rooms.filter((r) => u.includes(r.ID)).map((s) => s.Name);
	}

	async function doDeleteRecurring() {
		await deleteRecurringOccupancy(data.ID);
		goto('/occupancy');
	}

	async function doUpdateRecurring() {
		let c = {
			ID: Number(data.ID),
			Name: data.Name,
			Weekdays: selectedwd
				.map((n, i) => {
					if (n) {
						return i;
					}
				})
				.filter((o) => {
					return o !== undefined;
				}),
			StartTime: toZTimestring(data.StartTime),
			EndTime: toZTimesttring(data.EndTime),
			Rooms: data.SystemRooms.filter((r) => r.selected).map((r) => r.ID)
		};
		try {
			await putRecurringOccupancy(c);
			goto('/occupancy');
		} catch (e) {
			console.log(e);
			toast.push(e);
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
				<TableBodyCell>{data.ID}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Name</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={data.Name} /></TableBodyCell>
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
				<TableBodyCell
					>Start Time <br />either 24-hour "hh:mm" or 12-hour "hh:mm PM" format<br /> local timezone</TableBodyCell
				>
				<TableBodyCell><Input type="text" bind:value={data.StartTime} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>End Time</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={data.EndTime} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Rooms</TableBodyCell>
				<TableBodyCell>
					{#each data.SystemRooms as r}
						<Checkbox value={r.ID} bind:checked={r.selected}>{r.Name}</Checkbox>
					{/each}
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell><Button on:click={() => doDeleteRecurring()}>Delete</Button></TableBodyCell>
				<TableBodyCell><Button on:click={() => doUpdateRecurring()}>Update</Button></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
