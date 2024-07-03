<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let location = '/';
	if ($page.url.searchParams.has('u')) {
		location = $page.url.searchParams.get('u');
	}

	goto(location);
</script>
