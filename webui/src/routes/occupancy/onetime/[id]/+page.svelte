<script>
	import { goto } from '$app/navigation';
	import { putOneTimeOccupancy, deleteOneTimeOccupancy } from '$lib/occupancy';
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

	export let data;
	console.log(data);

	data.SystemRooms.forEach((r) => {
		r.selected = false;
		if (data.Rooms.includes(r.ID)) r.selected = true;
	});

	function parseRooms(r) {
		const u = new TextEncoder().encode(atob(r));
		return data.Rooms.filter((r) => u.includes(r.ID)).map((s) => s.Name);
	}

	async function doDeleteOneTime() {
		await deleteOneTimeOccupancy(data.ID);
		goto('/occupancy');
	}

	async function doUpdateOneTime() {
		let c = {
			ID: Number(data.ID),
			Name: data.Name,
			Start: data.Start,
			End: data.End,
			Rooms: data.SystemRooms.filter((r) => r.selected).map((r) => r.ID)
		};
		await putOneTimeOccupancy(c);
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
				<TableBodyCell>{data.ID}</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Name</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={data.Name} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Start</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={data.Start} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>End</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={data.End} /></TableBodyCell>
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
				<TableBodyCell><Button on:click={() => doDeleteOneTime()}>Delete</Button></TableBodyCell>
				<TableBodyCell><Button on:click={() => doUpdateOneTime()}>Update</Button></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
