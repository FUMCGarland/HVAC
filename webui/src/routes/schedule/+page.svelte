<script>
	import { onMount } from 'svelte';
	import { postSchedule } from '$lib/hvac';
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
	import { lowestFreeID } from '$lib/util';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	const weekdays = ['Sun', 'M', 'T', 'W', 'Th', 'F', 'Sat'];
	const selectedwd = weekdays.map(() => false);

	export let data;
	data.Zones.forEach((z) => {
		z.selected = false;
	});

	let id = lowestFreeID(data.Schedule);
	let name = 'not set';
	$: mode = 0;
	let starttime = '05:00';
	let runtime = 60;

	function modeString(mode) {
		if (mode == 0) return 'heating';
		return 'cooling';
	}

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	// shoot me, I've become one of those obnoxious js devs who do stuff in non-obvious one-liners
	function parseZones(z) {
		const u = new TextEncoder().encode(atob(z));
		return data.Zones.filter((z) => u.includes(z.ID)).map((s) => s.Name);
	}

	async function doAdd() {
		let c = {
			ID: Number(id),
			Name: name,
			Mode: Number(mode),
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
			Zones: data.Zones.filter((z) => z.selected).map((z) => z.ID)
		};
		await postSchedule(c);
	}
</script>

<Heading tag="h2">Schedule</Heading>
<form>
	<Table>
		<TableHead>
			<TableHeadCell colspan="2">Name</TableHeadCell>
			<TableHeadCell>System Mode</TableHeadCell>
			<TableHeadCell>Weekdays</TableHeadCell>
			<TableHeadCell>Start Time(s)</TableHeadCell>
			<TableHeadCell>Run Duration</TableHeadCell>
			<TableHeadCell>Zones</TableHeadCell>
		</TableHead>
		<TableBody>
			{#each data.Schedule as sched}
				<TableBodyRow>
					<TableBodyCell colspan="2"
						><A href="/schedule/{sched.ID}">{sched.ID}: {sched.Name}</A></TableBodyCell
					>
					<TableBodyCell>{modeString(sched.Mode)}</TableBodyCell>
					<TableBodyCell>{parseWeekdays(sched.Weekdays)}</TableBodyCell>
					<TableBodyCell>{sched.StartTime}</TableBodyCell>
					<TableBodyCell>{sched.RunTime}</TableBodyCell>
					<TableBodyCell>{parseZones(sched.Zones)}</TableBodyCell>
				</TableBodyRow>
			{/each}
		</TableBody>
	</Table>
	<Hr />
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
				<TableBodyCell>Mode</TableBodyCell>
				<TableBodyCell>
					<Button>
						Mode<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" />
					</Button>
					<Dropdown class="w-44 space-y-3 p-3 text-sm">
						<li><Radio name="mode" bind:group={mode} value={0}>Heating</Radio></li>
						<li><Radio name="mode" bind:group={mode} value={1}>Cooling</Radio></li>
					</Dropdown>
				</TableBodyCell>
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
				<TableBodyCell>Run Time (minutes)</TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={runtime} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Zones</TableBodyCell>
				<TableBodyCell>
					<Button
						>Zones<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" /></Button
					>
					<Dropdown class="w-44 space-y-3 p-3 text-sm">
						{#each data.Zones as z}
							<li>
								<Checkbox value={z.ID} bind:checked={z.selected}>{z.Name}</Checkbox>
							</li>
						{/each}
					</Dropdown>
				</TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>&nbsp;</TableBodyCell>
				<TableBodyCell><Button on:click={doAdd}>Add</Button></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
