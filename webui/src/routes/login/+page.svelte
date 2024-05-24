<script>
	import { goto } from '$app/navigation';
	import { Table, TableBody, TableBodyCell, TableBodyRow, Button, Input } from 'flowbite-svelte';
	import { page } from '$app/stores';
	import { hvaccontroller } from '$lib/hvac';
	import { toast } from '@zerodevx/svelte-toast';

	// export let data;
	let username;
	let password;

	// if you come to this page and are logged in, log out...
	const jwt = localStorage.getItem('jwt');
	if (jwt) {
		localStorage.removeItem('jwt');
	}

	async function doLogin() {
		try {
			await getJWT(username, password);
			goto('/');
		} catch (e) {
			console.log(e);
			toast.push(e);
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
