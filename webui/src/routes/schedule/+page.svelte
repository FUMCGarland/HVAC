<script>
	import { onMount } from 'svelte';
	import { invalidateAll } from '$app/navigation';
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
		P,
		A
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	const weekdays = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
	const selectedwd = weekdays.map(() => false);

	export let data;
	let id = data.Schedule.length; // get highest number and increment
	let name = 'not set';
	$: mode = 0;
	let starttime = '05:00';
	let runtime = 60;
	$: visiblepumps = data.Pumps.filter((p) => {
		return mode == p.SystemMode;
	});
	$: activeloops = visiblepumps.map((p) => {
		if (p.selected === true) return p.Loop;
	});
	$: visibleblowers = data.Blowers.filter((b) => {
		if (activeloops.includes(b.HotLoop) || activeloops.includes(b.ColdLoop)) return b;
	});

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	function parsePumps(p) {
		console.log(p);
		return 'tbp';
	}

	function parseBlowers(b) {
		console.log(b);
		return 'tbp';
	}

	async function doAdd() {
		const c = {
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
		await postSchedule(JSON.stringify(c));
	}
</script>

<Heading tag="h2">Schedule</Heading>
<form>
	<Table>
		<TableHead>
			<TableHeadCell>ID</TableHeadCell>
			<TableHeadCell>Name</TableHeadCell>
			<TableHeadCell>System Mode</TableHeadCell>
			<TableHeadCell>Weekdays</TableHeadCell>
			<TableHeadCell>Start Time(s)</TableHeadCell>
			<TableHeadCell>Run Duration</TableHeadCell>
			<TableHeadCell>Pumps</TableHeadCell>
			<TableHeadCell>Blowers</TableHeadCell>
		</TableHead>
		<TableBody>
			{#each data.Schedule as sched}
				<TableBodyRow>
					<TableBodyCell><A href="/schedule/{sched.ID}">{sched.ID}</A></TableBodyCell>
					<TableBodyCell>{sched.Name}</TableBodyCell>
					<TableBodyCell>{sched.Mode}</TableBodyCell>
					<TableBodyCell>{parseWeekdays(sched.Weekdays)}</TableBodyCell>
					<TableBodyCell>{sched.StartTime}</TableBodyCell>
					<TableBodyCell>{sched.RunTime}</TableBodyCell>
					<TableBodyCell>{parsePumps(sched.Pumps)}</TableBodyCell>
					<TableBodyCell>{parseBlowers(sched.Blowers)}</TableBodyCell>
				</TableBodyRow>
			{/each}
			<TableBodyRow>
				<TableBodyCell><Input type="text" bind:value={id} /></TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={name} /></TableBodyCell>
				<TableBodyCell>
					<Button>
						Mode<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" />
					</Button>
					<Dropdown class="w-44 space-y-3 p-3 text-sm">
						<li><Radio name="mode" bind:group={mode} value={0}>Heating</Radio></li>
						<li><Radio name="mode" bind:group={mode} value={1}>Cooling</Radio></li>
					</Dropdown>
				</TableBodyCell>
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
				<TableBodyCell><Input type="text" bind:value={starttime} /></TableBodyCell>
				<TableBodyCell><Input type="text" bind:value={runtime} /></TableBodyCell>
				<TableBodyCell>
					<Button
						>Pumps<ChevronDownOutline class="ms-2 h-6 w-6 text-white dark:text-white" /></Button
					>
					<Dropdown class="w-44 space-y-3 p-3 text-sm">
						{#each visiblepumps as pump, index}
							<li>
								<Checkbox value={pump.ID} bind:checked={pump.selected}>{pump.Name}</Checkbox>
							</li>
						{/each}
					</Dropdown>
				</TableBodyCell>
				<TableBodyCell>
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
				<TableBodyCell colspan="7">&nbsp;</TableBodyCell>
				<TableBodyCell><Button on:click={doAdd}>Add</Button></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
