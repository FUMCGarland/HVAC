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
		Input,
		Label,
		Helper,
		Button
	} from 'flowbite-svelte';
	import { zoneStart, zoneStop } from '$lib/hvac.js';

	export let data;

	// add a private value for to track display state
	data.Zones.forEach((z) => {
		z.newRunTime = 0;
	});

	// refresh every 30 seconds
	onMount(() => {
		const interval = setInterval(() => {
			invalidateAll();
		}, 30000);

		return () => {
			clearInterval(interval);
		};
	});

	// a zone is running if all the blowers/radiant loops in the zone are running
	// TODO: check radiant heat loops
	// ignore pumps and chillers (?)
	function zoneRunning(zone) {
		const blowers = data.Blowers.filter((b) => b.Zone == zone);
		const running = blowers.filter((b) => b.Running);
		return blowers.length == running.length && blowers.length != 0;
	}
</script>

<Heading tag="h2">Zones</Heading>
<Table>
	<TableHead>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
		<TableHeadCell>Command</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Zones as zone}
			<TableBodyRow>
				<TableBodyCell><A href="/zone/{zone.ID}">{zone.Name}</A></TableBodyCell>
				{#if zoneRunning(zone.ID)}
					<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
				{/if}
				{#if !zoneRunning(zone.ID)}
					<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
				{/if}
				<TableBodyCell>
					{#if zoneRunning(zone.ID)}
						<form>
							<div class="mb-6 grid gap-6">
								<div>
									<Button
										type="submit"
										on:click={(x) => {
											// console.log(x.srcElement);
											x.srcElement.textContent = 'Processing...';
											x.srcElement.disable = true;
											zoneStop(zone.ID);
										}}>Stop</Button
									>
								</div>
							</div>
						</form>
					{/if}
					{#if !zoneRunning(zone.ID)}
						<form>
							<div class="mb-6 grid gap-6">
								<div>
									<Label for="run_time${zone.ID}" class="mb-2">Run Time (minutes)</Label>
									<Input
										type="text"
										id="run_time{zone.ID}"
										placeholder="60"
										required
										bind:value={zone.newRunTime}
									/>
									<Button
										type="submit"
										on:click={(x) => {
											// console.log(x.srcElement);
											x.srcElement.textContent = 'Processing...';
											x.srcElement.disable = true;
											zoneStart(zone.ID, zone.newRunTime);
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
