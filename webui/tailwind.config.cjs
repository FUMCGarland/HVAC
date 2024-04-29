/** @type {import('tailwindcss').Config}*/
// const config = { content: ['./src/**/*.{html,js,svelte,ts}'], theme: { extend: {} }, plugins: [] };
const config = {
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		'./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'
	],

	plugins: [require('flowbite/plugin')],

	darkMode: 'class',

	theme: {
		extend: {
			colors: {
				// flowbite-svelte
				primary: {
					50: '#f2f2ff',
					100: '#eef5ff',
					200: '#dee4ff',
					300: '#cccdff',
					400: '#c8adff',
					500: '#5d98fe',
					600: '#552fef',
					700: '#275beb',
					800: '#223bcc',
					900: '#1b32a5'
				}
			}
		}
	}
};

module.exports = config;
