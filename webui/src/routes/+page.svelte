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
</script>

<P><A href="/systemmode">System Mode</A> {systemModeLabel(data.SystemMode)}</P>
<P><A href="/controlmode">Control Mode</A> {controlModeLabel(data.ControlMode)}</P>
<Hr />
<Heading tag="h2">Blowers</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Hot Loop</TableHeadCell>
		<TableHeadCell>Cold Loop</TableHeadCell>
		<TableHeadCell>Zone</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Blowers as blower}
			<TableBodyRow>
				<TableBodyCell><A href="/blower/{blower.ID}">{blower.ID}</A></TableBodyCell>
				<TableBodyCell>{blower.Name}</TableBodyCell>
				<TableBodyCell><A href="/loop/{blower.HotLoop}">{blower.HotLoop}</A></TableBodyCell>
				<TableBodyCell><A href="/loop/{blower.ColdLoop}">{blower.ColdLoop}</A></TableBodyCell>
				<TableBodyCell><A href="/zone/{blower.Zone}">{blower.Zone}</A></TableBodyCell>
				{#if blower.Running}
					<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
				{/if}
				{#if !blower.Running}
					<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
				{/if}
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
<Hr />
<Heading tag="h2">Pumps ({systemModeLabel(data.SystemMode)})</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Loop</TableHeadCell>
		<TableHeadCell>System Mode</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Pumps as pump}
			{#if inMode(pump)}
				<TableBodyRow>
					<TableBodyCell><A href="/pump/{pump.ID}">{pump.ID}</A></TableBodyCell>
					<TableBodyCell>{pump.Name}</TableBodyCell>
					<TableBodyCell><A href="/loop/{pump.Loop}">{pump.Loop}</A></TableBodyCell>
					<TableBodyCell>{systemModeLabel(pump.SystemMode)}</TableBodyCell>
					{#if pump.Running}
						<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
					{/if}
					{#if !pump.Running}
						<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
					{/if}
				</TableBodyRow>
			{/if}
		{/each}
	</TableBody>
</Table>
