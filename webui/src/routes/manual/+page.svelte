<script>
	import { onMount } from 'svelte';
	import { invalidateAll } from '$app/navigation';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell
	} from 'flowbite-svelte';
	import { Heading, P, A, Input, Label, Helper, Button } from 'flowbite-svelte';
	import { blowerStart, blowerStop, pumpStart, pumpStop } from '$lib/hvac.js';

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
		if (sm == 1) return 'cooling';
		return 'unknown';
	}

	function systemControlModeLabel(scm) {
		if (scm == 0) return 'manual';
		if (scm == 1) return 'schedule';
		if (scm == 2) return 'temp';
		if (scm == 3) return 'off';
		return 'manual';
	}

	function inMode(pump) {
		if (data.SystemMode == 0 && pump.Hot) return true;
		if (data.SystemMode == 1 && !pump.Hot) return true;
		return false;
	}
</script>

<Heading tag="h2">Blowers</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Hot Loop</TableHeadCell>
		<TableHeadCell>Cold Loop</TableHeadCell>
		<TableHeadCell>Zone</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
		<TableHeadCell>Command</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Blowers as blower}
			<TableBodyRow>
				<TableBodyCell><A href="/blower/{blower.ID}">{blower.ID}</A></TableBodyCell>
				<TableBodyCell>{blower.Name}</TableBodyCell>
				<TableBodyCell><A href="/loop/{blower.HotLoop}">{blower.HotLoop}</A></TableBodyCell>
				<TableBodyCell><A href="/loop/{blower.ColdLoop}">{blower.ColdLoop}</A></TableBodyCell>
				<TableBodyCell><A href="/zone/{blower.Zone}">{blower.Zone}</A></TableBodyCell>
				<TableBodyCell>{blower.Running}</TableBodyCell>
				<TableBodyCell>
					{#if blower.Running}
						<form>
							<div class="mb-6 grid gap-6">
								<div>
									<Button
										type="submit"
										on:click={() => {
											blowerStop(blower.ID);
										}}>Stop</Button
									>
								</div>
							</div>
						</form>
					{/if}
					{#if !blower.Running}
						<form>
							<div class="mb-6 grid gap-6">
								<div>
									<Label for="run_time${blower.ID}" class="mb-2">Run Time (minutes)</Label>
									<Input
										type="text"
										id="run_time{blower.ID}"
										placeholder="60"
										required
										bind:this={blower.newRunTime}
									/>
									<Button
										type="submit"
										on:click={() => {
											blowerStart(blower.ID, blower.newRunTime.$capture_state().value);
										}}>Start</Button
									>
								</div>
							</div>
						</form>
					{/if}
				</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
<Heading tag="h2">Pumps</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Loop</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
		<TableHeadCell>Command</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Pumps as pump}
			{#if inMode(pump)}
				<TableBodyRow>
					<TableBodyCell><A href="/pump/{pump.ID}">{pump.ID}</A></TableBodyCell>
					<TableBodyCell>{pump.Name}</TableBodyCell>
					<TableBodyCell><A href="/loop/{pump.Loop}">{pump.Loop}</A></TableBodyCell>
					<TableBodyCell>{pump.Running}</TableBodyCell>
					<TableBodyCell>
						{#if pump.Running}
							<form>
								<div class="mb-6 grid gap-6">
									<div>
										<Button
											type="submit"
											on:click={() => {
												pumpStop(pump.ID);
											}}>Stop</Button
										>
									</div>
								</div>
							</form>
						{/if}
						{#if !pump.Running}
							<form>
								<div class="mb-6 grid gap-6">
									<div>
										<Label for="run_time${pump.ID}" class="mb-2">Run Time (minutes)</Label>
										<Input
											type="text"
											id="run_time{pump.ID}"
											placeholder="60"
											required
											bind:this={pump.newRunTime}
										/>
										<Button
											type="submit"
											on:click={() => {
												pumpStart(pump.ID, pump.newRunTime.$capture_state().value);
											}}>Start</Button
										>
									</div>
								</div>
							</form>
						{/if}
					</TableBodyCell>
				</TableBodyRow>
			{/if}
		{/each}
	</TableBody>
</Table>
