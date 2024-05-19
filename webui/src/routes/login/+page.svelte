<script>
	import { goto } from '$app/navigation';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { page } from '$app/stores';
	import { hvaccontroller } from '$lib/hvac';
	import { toast } from '@zerodevx/svelte-toast';

	// export let data;
	let username;
	let password;

	async function doLogin() {
		console.log(username, password);
		try {
			await getJWT(username, password);
			goto('/');
		} catch (e) {
			console.log(e);
		}
	}

	async function getJWT(username, password) {
		const dataArray = new FormData();
		dataArray.append('username', username);
		dataArray.append('password', password);

		const request = {
			method: 'POST',
			mode: 'cors',
			credentials: 'include',
			redirect: 'manual',
			referrerPolicy: 'origin',
			body: dataArray
		};

		const response = await fetch(`${hvaccontroller}/api/v1/getJWT`, request);
		const payload = await response.text();

		if (response.status != 200) {
			console.log('server returned ', response.status);
			toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
			return;
		}
		// stuff it away for no particular reason, we will be using the cookies because svelte and safari and request.header == badness
		localStorage.setItem('jwt', payload);
	}
</script>

<form
	on:submit={() => {
		doLogin();
	}}
>
	<Table>
		<TableBody>
			<TableBodyRow>
				<TableBodyCell>Username:</TableBodyCell>
				<TableBodyCell><Input type="text" name="username" bind:value={username} /></TableBodyCell>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>Password:</TableBodyCell>
				<TableBodyCell
					><Input type="password" name="password" bind:value={password} /></TableBodyCell
				>
			</TableBodyRow>
			<TableBodyRow>
				<TableBodyCell>&nbsp;</TableBodyCell>
				<TableBodyCell><Button type="submit">Login</Button></TableBodyCell>
			</TableBodyRow>
		</TableBody>
	</Table>
</form>
