<script>
	import {
		Heading,
	} from 'flowbite-svelte';

	export let data;
</script>

<Heading tag="h2">Run Log</Heading>
<pre>{data.data}</pre>
