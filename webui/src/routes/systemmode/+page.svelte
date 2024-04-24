<script>
	import { P, Button, Dropdown, DropdownItem } from 'flowbite-svelte';
	import { setSystemMode } from '$lib/hvac.js';

	export let data;
	let smstring = systemModeLabel(data.SystemMode);
	export let dropdownOpen = false;

	function systemModeLabel(sm) {
		if (sm == 0) return 'heat';
		if (sm == 1) return 'cool';
		return 'heat';
	}

	async function setSM(mode) {
		smstring = systemModeLabel(mode);
		await setSystemMode(mode);
	}
</script>

<P>
	This button changes the systems's operational mode. The control mode must be "off" to change the
	system mode. Set the control mode off. Go spend several hours changing valves and such, then come
	back and set the new system mode. Then set the control mode once things are working.
</P>

<Button>{smstring}</Button>
<Dropdown smstring>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setSM(0);
		}}>Heat</DropdownItem
	>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setSM(1);
		}}>Cool</DropdownItem
	>
</Dropdown>
