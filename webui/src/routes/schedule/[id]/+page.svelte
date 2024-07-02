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
		Button,
		Dropdown,
		Checkbox,
		Input
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import { putSchedule, deleteSchedule } from '$lib/hvac.js';
	import { redirect } from '@sveltejs/kit';

	export let data;

	const weekdays = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
	const selectedwd = weekdays.map(() => false);
	data.Weekdays.forEach((d) => {
		selectedwd[d] = true;
	});
	data.System.Zones.forEach((z) => {
		z.selected = false;
	});
	parseZones(data.Zones).forEach((zid) => {
		data.System.Zones.forEach((z) => {
			if (zid == z.ID) z.selected = true;
		});
	});

	function modeString(mode) {
		if (mode == 0) return 'heating';
		return 'cooling';
	}

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	// comes across as a byte string of UTF8...
	function parseZones(p) {
		const e = new TextEncoder();
		const u = e.encode(atob(p));
		return u;
	}

	async function doEdit() {
		const c = {
			ID: Number(data.ID),
			Name: data.Name,
			Mode: Number(data.Mode),
			Weekdays: selectedwd
				.map((n, i) => {
					if (n) {
						return i;
					}
				})
				.filter((o) => {
					return o !== undefined;
				}),
			Starttime: data.StartTime,
			Runtime: Number(data.RunTime),
			Zones: data.System.Zones.filter((z) => z.selected).map((z) => z.ID)
		};
		await putSchedule(c);
		data.RunTime = data.RunTime;
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
			<TableBodyCell><Input type="text" bind:value={data.Name} /></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Mode</TableBodyCell>
			<TableBodyCell>{modeString(data.Mode)}</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Start Times</TableBodyCell>
			<TableBodyCell><Input type="text" bind:value={data.StartTime} /></TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Run Time (minutes)</TableBodyCell>
			<TableBodyCell><Input type="text" bind:value={data.RunTime} /></TableBodyCell>
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
			<TableBodyCell>Zones</TableBodyCell>
			<TableBodyCell>
				<Button>Zones<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" /></Button>
				<Dropdown class="w-44 space-y-3 p-3 text-sm">
					{#each data.System.Zones as z}
						<li>
							<Checkbox value={z.ID} bind:checked={z.selected}>{z.Name}</Checkbox>
						</li>
					{/each}
				</Dropdown>
			</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>
				<Button
					on:click={() => {
						deleteSchedule(data.ID);
						redirect(303, '/schedule');
					}}>Delete</Button
				>
			</TableBodyCell>
			<TableBodyCell>
				<Button
					on:click={() => {
						doEdit();
					}}>Update</Button
				>
			</TableBodyCell>
		</TableBodyRow>
	</TableBody>
</Table>
