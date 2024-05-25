<script>
	import {
		Navbar,
		NavBrand,
		NavLi,
		NavUl,
		NavHamburger,
		Toast,
		Dropdown,
		DropdownItem,
		DropdownDivider
	} from 'flowbite-svelte';
	import { ChevronDownOutline } from 'flowbite-svelte-icons';
	import { SvelteToast } from '@zerodevx/svelte-toast';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import '../app.pcss';
	import { getAccessLevel } from '$lib/util';

	const level = getAccessLevel();
</script>

<svelte:head>
	<title>FUMC Garland HVAC Controller</title>
</svelte:head>

<svelte:window />
<header class="mx-auto w-full flex-none bg-white dark:bg-slate-950">
	<Navbar>
		<NavBrand href="/">
			<img src="/fumcg.jpg" class="me-3 h-6 sm:h-9" alt="FUMCG Logo" />
			<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
				>FUMC Garland HVAC</span
			>
		</NavBrand>
		<NavHamburger />
		<NavUl>
			<NavLi href="/rooms">Room Status</NavLi>
			<NavLi href="/occupancy">Occupancy Schedule</NavLi>
			<NavLi class="cursor-pointer">
				Advanced<ChevronDownOutline class="ms-2 inline h-6 w-6 text-primary-800 dark:text-white" />
			</NavLi>

			<Dropdown class="z-20 w-44">
				<DropdownItem href="/zones">Zones Temp Settings</DropdownItem>
				<DropdownItem href="/shelly">Sensor Status</DropdownItem>
				<DropdownItem href="/api/v1/datalog">Download datalog</DropdownItem>
				<DropdownItem href="/settings">System Settings</DropdownItem>
				<DropdownDivider />
				{#if level > 0}
					<DropdownItem href="/schedule">Zone Schedule</DropdownItem>
					<DropdownItem href="/manual">Manual Zone Control</DropdownItem>
					<DropdownItem href="/override">Manual Device Override</DropdownItem>
				{/if}
				{#if level > 1}
					<DropdownDivider />
					<DropdownItem href="/systemmode">System Mode</DropdownItem>
					<DropdownItem href="/controlmode">Control Mode</DropdownItem>
				{/if}
				<DropdownItem href="/login">Logout</DropdownItem>
			</Dropdown>
		</NavUl>
	</Navbar>
	<SvelteToast />
	<div class="flex gap-10"></div>
</header>

<div class="mx-auto flex w-full px-4">
	<main class="mx-auto w-full">
		<slot></slot>
	</main>
</div>
