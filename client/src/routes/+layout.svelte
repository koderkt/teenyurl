<script>
	import NavBar from '$lib/NavBar.svelte';
	import '../app.css';
	export let data;
</script>

<NavBar {data} />
<slot></slot>
