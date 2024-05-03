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
		TableHeadCell
	} from 'flowbite-svelte';
	import { Heading, P, A, Input, Label, Helper, Button } from 'flowbite-svelte';
	import { zoneStart, zoneStop } from '$lib/hvac.js';

	export let data;
	data.Zones.forEach((z) => {
		z.newRunTime = 0;
	});

	onMount(() => {
		const interval = setInterval(() => {
			invalidateAll();
		}, 30000);

		return () => {
			clearInterval(interval);
		};
	});

	function zoneRunning(zone) {
		return false;
	}
</script>

<Heading tag="h2">Zones</Heading>
<Table>
	<TableHead>
		<TableHeadCell>ID</TableHeadCell>
		<TableHeadCell>Name</TableHeadCell>
		<TableHeadCell>Running</TableHeadCell>
		<TableHeadCell>Command</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Zones as zone}
			<TableBodyRow>
				<TableBodyCell><A href="/zone/{zone.ID}">{zone.ID}</A></TableBodyCell>
				<TableBodyCell>{zone.Name}</TableBodyCell>
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
