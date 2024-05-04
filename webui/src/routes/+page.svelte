<script>
	import { onMount } from 'svelte';
	import { invalidateAll } from '$app/navigation';
	import {
		Badge,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		P,
		A,
		Hr
	} from 'flowbite-svelte';

	export let data;

	// refresh every 30 seconds
	onMount(() => {
		const interval = setInterval(() => {
			invalidateAll();
		}, 30000);

		return () => {
			clearInterval(interval);
		};
	});

	function systemModeLabel(sm) {
		if (sm == 0) return 'heating';
		return 'cooling';
	}

	function controlModeLabel(scm) {
		if (scm == 0) return 'manual';
		if (scm == 1) return 'schedule';
		if (scm == 2) return 'temp';
		if (scm == 3) return 'off';
		return 'manual';
	}

	function inMode(pump) {
		return data.SystemMode == pump.SystemMode;
	}

	function zoneName(zoneID) {
		const z = data.Zones.filter((z) => z.ID == zoneID);
		return z[0].Name;
	}

	function loopName(loopID) {
		const l = data.Loops.filter((l) => l.ID == loopID);
		return l[0].Name;
	}
</script>

<P><A href="/systemmode">System Mode</A> {systemModeLabel(data.SystemMode)}</P>
<P><A href="/controlmode">Control Mode</A> {controlModeLabel(data.ControlMode)}</P>
<Hr />
<Heading tag="h2">Blowers</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
		<TableHeadCell>Zone</TableHeadCell>
		<TableHeadCell>Loop</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Blowers as blower}
			<TableBodyRow>
				<TableBodyCell><A href="/blower/{blower.ID}">{blower.Name}</A></TableBodyCell>
				{#if blower.Running}
					<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
				{/if}
				{#if !blower.Running}
					<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
				{/if}
				<TableBodyCell><A href="/zone/{blower.Zone}">{zoneName(blower.Zone)}</A></TableBodyCell>
				{#if data.SystemMode == 0}
					<TableBodyCell
						><A href="/loop/{blower.HotLoop}">{loopName(blower.HotLoop)}</A></TableBodyCell
					>
				{/if}
				{#if data.SystemMode == 1}
					<TableBodyCell
						><A href="/loop/{blower.ColdLoop}">{loopName(blower.ColdLoop)}</A></TableBodyCell
					>
				{/if}
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
<Hr />
<Heading tag="h2">Pumps ({systemModeLabel(data.SystemMode)})</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
		<TableHeadCell>Loop</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Pumps as pump}
			{#if inMode(pump)}
				<TableBodyRow>
					<TableBodyCell><A href="/pump/{pump.ID}">{pump.Name}</A></TableBodyCell>
					{#if pump.Running}
						<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
					{/if}
					{#if !pump.Running}
						<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
					{/if}
					<TableBodyCell><A href="/loop/{pump.Loop}">{loopName(pump.Loop)}</A></TableBodyCell>
				</TableBodyRow>
			{/if}
		{/each}
	</TableBody>
</Table>
{#if data.SystemMode == 1}
	<Hr />
	<Heading tag="h2">Chillers</Heading>
	<Table>
		<TableHead>
			<TableHeadCell>Name</TableHeadCell>
			<TableHeadCell>Running</TableHeadCell>
		</TableHead>
		<TableBody>
			{#each data.Chillers as chiller}
				<TableBodyRow>
					<TableBodyCell><A href="/chiller/{chiller.ID}">{chiller.Name}</A></TableBodyCell>
					{#if chiller.Running}
						<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
					{/if}
					{#if !chiller.Running}
						<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
					{/if}
				</TableBodyRow>
			{/each}
		</TableBody>
	</Table>
{/if}
