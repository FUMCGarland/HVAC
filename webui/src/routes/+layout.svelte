<script>
	import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Toast } from 'flowbite-svelte';
	import { SvelteToast } from '@zerodevx/svelte-toast';
	import '../app.pcss';
</script>

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
			<NavLi href="/occupancy">Occupancy</NavLi>
			<NavLi href="/rooms">Rooms</NavLi>
			<NavLi href="/zones">Zones</NavLi>
			<NavLi href="/schedule">Schedule</NavLi>
			<NavLi href="/manual">Manual</NavLi>
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
