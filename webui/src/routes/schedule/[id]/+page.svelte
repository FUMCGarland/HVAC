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

	$: visiblepumps = data.System.Pumps.filter((p) => {
		return data.Mode == p.SystemMode;
	});
	$: activeloops = visiblepumps.map((p) => {
		if (p.selected === true) return p.Loop;
	});
	$: visibleblowers = data.System.Blowers.filter((b) => {
		if (activeloops.includes(b.HotLoop) || activeloops.includes(b.ColdLoop)) return b;
	});

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
			Pumps: visiblepumps
				.filter((p) => {
					if (p.selected === true) return p;
				})
				.map((p) => p.ID),
			Blowers: visibleblowers
				.filter((b) => {
					if (b.selected === true) return b;
				})
				.map((b) => b.ID)
		};
		await putSchedule(c);
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
				{parseWeekdays(data.Weekdays)}<br />
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
			<TableBodyCell>Pumps</TableBodyCell>
			<TableBodyCell
				>{parsePumps(data.Pumps)}<br />
				<Button>Pumps<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" /></Button>
				<Dropdown class="w-44 space-y-3 p-3 text-sm">
					{#each visiblepumps as pump, index}
						<li>
							<Checkbox value={pump.ID} bind:checked={pump.selected}>{pump.Name}</Checkbox>
						</li>
					{/each}
				</Dropdown>
			</TableBodyCell>
		</TableBodyRow>
		<TableBodyRow>
			<TableBodyCell>Blowers</TableBodyCell>
			<TableBodyCell
				>{parseBlowers(data.Blowers)}<br />
				<Button
					>Blowers<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" /></Button
				>
				<Dropdown class="w-44 space-y-3 p-3 text-sm">
					{#each visibleblowers as blower, index}
						<li>
							<Checkbox value={blower.ID} bind:checked={blower.selected}>{blower.Name}</Checkbox>
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
						throw redirect(303, '/schedule');
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
