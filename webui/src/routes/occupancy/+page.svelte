<script>
	import { goto } from '$app/navigation';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		Button,
		A,
		Banner
	} from 'flowbite-svelte';
	import { BullhornSolid } from 'flowbite-svelte-icons';
	import { toLocaleTimestring } from '$lib/util.js';

	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	const weekdays = ['Sun', 'M', 'T', 'W', 'Th', 'F', 'Sat'];
	export let data;

	function parseWeekdays(w) {
		return w.map((p) => weekdays[p]);
	}

	function formatDate(d) {
		return new Date(Date.parse(d)).toLocaleString();
	}
</script>

{#if data.ControlMode != 2}
	<Banner id="default-banner">
		<!-- position="absolute" -->
		<p class="flex items-center text-sm font-normal text-gray-500 dark:text-gray-400">
			<span class="me-3 inline-flex rounded-full bg-gray-200 p-1 dark:bg-gray-600">
				<BullhornSolid class="h-3 w-3 text-gray-500 dark:text-gray-400" />
				<span class="sr-only">Light bulb</span>
			</span>
			<span>
				The system is not in "temp" control mode; you can configure the occupancy schedule, but it
				will not be active until the control mode is switched.
			</span>
		</p>
	</Banner>
{/if}
<Heading tag="h2">Recurring Occupancy Events</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Rooms</TableHeadCell>
		<TableHeadCell>Weekdays</TableHeadCell>
		<TableHeadCell>Start Time (local timezone)</TableHeadCell>
		<TableHeadCell>End Time</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Recurring as r}
			<TableBodyRow>
				<TableBodyCell><A href="/occupancy/recurring/{r.ID}">{r.Name}</A></TableBodyCell>
				<TableBodyCell>
					{#each r.Rooms as room}
						<A href="/room/{room}">{room}</A>&nbsp;
					{/each}
				</TableBodyCell>
				<TableBodyCell>{parseWeekdays(r.Weekdays)}</TableBodyCell>
				<TableBodyCell>{toLocaleTimestring(r.StartTime)}</TableBodyCell>
				<TableBodyCell>{toLocaleTimestring(r.EndTime)}</TableBodyCell>
			</TableBodyRow>
		{/each}
		<TableBodyRow>
			<TableBodyCell colspan="4">&nbsp;</TableBodyCell>
			<TableBodyCell
				><Button on:click={() => goto('/occupancy/recurring')}>Add Recurring</Button></TableBodyCell
			>
		</TableBodyRow>
	</TableBody>
</Table>

<Heading tag="h2">One-Time Occupancy Events</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Rooms</TableHeadCell>
		<TableHeadCell>Start</TableHeadCell>
		<TableHeadCell>End</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.OneTime as o}
			<TableBodyRow>
				<TableBodyCell><A href="/occupancy/onetime/{o.ID}">{o.Name}</A></TableBodyCell>
				<TableBodyCell>{o.Rooms}</TableBodyCell>
				<TableBodyCell>{formatDate(o.Start)}</TableBodyCell>
				<TableBodyCell>{formatDate(o.End)}</TableBodyCell>
			</TableBodyRow>
		{/each}
		<TableBodyRow>
			<TableBodyCell colspan="3">&nbsp;</TableBodyCell>
			<TableBodyCell
				><Button on:click={() => goto('/occupancy/onetime')}>Add One-Time</Button></TableBodyCell
			>
		</TableBodyRow>
	</TableBody>
</Table>
