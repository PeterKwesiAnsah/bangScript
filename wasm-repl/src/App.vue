<script setup lang="ts">
import BangScriptREPL from "./components/BangScriptREPL.vue";
</script>

<template>
	<div class="app">
		<BangScriptREPL />
	</div>
</template>

<style>
* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

body {
	font-family: "Inter", sans-serif;
	background: #0a0a0a;
	color: #ffffff;
	overflow: hidden;
}

.app {
	width: 100vw;
	height: 100vh;
	background: linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 100%);
}
</style>
