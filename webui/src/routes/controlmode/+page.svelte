<script>
	import { P, List, Li } from 'flowbite-svelte';
	import { Button, Dropdown, DropdownItem } from 'flowbite-svelte';
	import { setSystemControlMode } from '$lib/hvac.js';

	export let data;
	let scmstring = systemControlModeLabel(data.SystemControlMode);
	export let dropdownOpen = false;

	function systemControlModeLabel(scm) {
		if (scm == 0) return 'manual';
		if (scm == 1) return 'schedule';
		if (scm == 2) return 'temp';
		if (scm == 3) return 'off';
		return 'manual';
	}

	async function setSCM(mode) {
		scmstring = systemControlModeLabel(mode);
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
			temperature.</Li
		>
		<Li
			>Off - nothing will run, this setting is only used when changing the system from heating to
			cooling.</Li
		>
	</List>
</P>

<Button>{scmstring}</Button>
<Dropdown scmstring>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setSCM(0);
		}}>Manual</DropdownItem
	>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setSCM(1);
		}}>Schedule</DropdownItem
	>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setSCM(2);
		}}>Temperature Based</DropdownItem
	>
	<DropdownItem
		on:click={() => {
			dropdownOpen = false;
			setSCM(3);
		}}>Off</DropdownItem
	>
</Dropdown>
