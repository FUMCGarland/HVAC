<script>
	import { P, List, Li, Button, Dropdown, DropdownItem } from 'flowbite-svelte';
	import { setSystemControlMode } from '$lib/hvac.js';

	export let data;
	let cmstring = controlModeLabel(data.ControlMode);
	export let dropdownOpen = false;

	function controlModeLabel(cm) {
		if (cm == 0) return 'manual';
		if (cm == 1) return 'schedule';
		if (cm == 2) return 'temp';
		if (cm == 3) return 'off';
		return 'unknown';
	}

	async function setCM(mode) {
		cmstring = controlModeLabel(mode);
		await setSystemControlMode(mode);
	}
</script>

<P>
	This button changes the controller's operational mode.
	<List>
		<Li
			>Manual is "full manual mode" - the system will only make changed based on user action in this
			UI.</Li
		>
		<Li>Schedule is "schedule mode" - the system will run based on the scheduled entries.</Li>
		<Li
			>Temperature - the system will run when sensors in the rooms indicate an out-of-range
			temperature. (NOT BUILT YET)</Li
		>
		<Li
			>Off - nothing will run. This setting is only used when changing the system from heating to
			cooling.</Li
		>
	</List>
</P>

<Button>{cmstring}</Button>
<Dropdown cmstring>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setCM(0);
		}}>Manual</DropdownItem
	>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setCM(1);
		}}>Schedule</DropdownItem
	>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setCM(2);
		}}>Temperature Based</DropdownItem
	>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setCM(3);
		}}>Off</DropdownItem
	>
</Dropdown>
