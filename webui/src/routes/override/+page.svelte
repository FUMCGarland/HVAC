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
		A,
		Input,
		Label,
		Button
	} from 'flowbite-svelte';
	import { blowerStart, blowerStop, pumpStart, pumpStop } from '$lib/hvac.js';

	export let data;
	data.Pumps.forEach((p) => {
		p.newRunTime = 0;
	});
	data.Blowers.forEach((b) => {
		b.newRunTime = 0;
	});

	onMount(() => {
		const interval = setInterval(() => {
			invalidateAll();
		}, 30000);

		return () => {
			clearInterval(interval);
		};
	});

	function inMode(pump) {
		return data.SystemMode == pump.SystemMode;
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
				{#if blower.Running}
					<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
				{/if}
				{#if !blower.Running}
					<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
				{/if}
				<TableBodyCell>
					{#if blower.Running}
						<form>
							<div class="mb-6 grid gap-6">
								<div>
									<Button
										type="submit"
										on:click={(x) => {
											x.srcElement.textContent = 'Processing...';
											x.srcElement.disable = true;
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
										bind:value={blower.newRunTime}
									/>
									<Button
										type="submit"
										on:click={(x) => {
											x.srcElement.textContent = 'Processing...';
											x.srcElement.disable = true;
											blowerStart(blower.ID, blower.newRunTime);
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
					{#if pump.Running}
						<TableBodyCell><Badge color="green">Running</Badge></TableBodyCell>
					{/if}
					{#if !pump.Running}
						<TableBodyCell><Badge color="red">Stopped</Badge></TableBodyCell>
					{/if}
					<TableBodyCell>
						{#if pump.Running}
							<form>
								<div class="mb-6 grid gap-6">
									<div>
										<Button
											type="submit"
											on:click={(x) => {
												x.srcElement.textContent = 'Processing...';
												x.srcElement.disable = true;
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
											bind:value={pump.newRunTime}
										/>
										<Button
											type="submit"
											on:click={(x) => {
												x.srcElement.textContent = 'Processing...';
												x.srcElement.disable = true;
												pumpStart(pump.ID, pump.newRunTime);
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
