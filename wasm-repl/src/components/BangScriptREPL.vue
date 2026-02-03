<script setup lang="ts">
import { ref, onMounted } from "vue";

type HistoryPromptItem = {
	indicator: ">>> " | "... ";
	input: string;
	output: {
		result: string;
		programstatus: 0 | 1;
	};
	more?: string[];
};

const primaryIndicator = ">>> ";

const inputRef = ref(null);
const history = ref<HistoryPromptItem[]>([]);
const currentInput = ref("");
const cursorVisible = ref(true);

const inputIndicator = ref<">>> " | ">... ">(primaryIndicator);
let scopeDepth = ref(0);

function processInput(input: string) {
	if (!input.trim()) return;

	if (inputIndicator.value === primaryIndicator) {
		//input can be complete or incomplete
		for (let i = 0; i < input.length; i++) {
			if (input.charAt(i) === "{") {
				scopeDepth.value++;
			} else if (input.charAt(i) === "}") {
				scopeDepth.value--;
			}
		}
		if (scopeDepth.value == 0) {
			//completes the prompt
			currentInput.value = "";
			//we execute, get to result and....
			//add to history
		} else {
			currentInput.value = "";
			//add to history
		}
	} else {
		if (scopeDepth.value === 0) {
			throw new Error(
				"Scope Depth should not be zero when more input is expected from use",
			);
		} else if (!history.value[history.value.length - 1].more) {
			throw new Error(
				"More HistoryPromptItems should be at least one or empty",
			);
		}
		for (let i = 0; i < input.length; i++) {
			if (input.charAt(i) === "{") {
				scopeDepth.value++;
			} else if (input.charAt(i) === "}") {
				scopeDepth.value--;
			}
		}
		if (scopeDepth.value == 0) {
			//completes the prompt
			inputIndicator.value = primaryIndicator;
			currentInput.value = "";
			//add to history
			//we execute, get to result and....
			//add to history
		} else {
			currentInput.value = "";
			//add to history
		}
	}
}

// Handle keydown events
const handleKeydown = (event) => {
	if (event.key === "Enter") {
		event.preventDefault();
		//const input = currentInput.value;
		currentInput.value = "";
	} else if (event.key === "Backspace") {
		event.preventDefault();
		if (currentInput.value.length > 0) {
			currentInput.value = currentInput.value.slice(0, -1);
		}
	} else if (event.key === "Tab") {
		event.preventDefault();
		currentInput.value += "  ";
	}
};

// Handle input events
const handleInput = (event) => {
	const target = event.target;
	currentInput.value = target.value;
};

// Blink cursor
setInterval(() => {
	cursorVisible.value = !cursorVisible.value;
}, 500);

// Focus input on mount and handle clicks
onMounted(() => {
	if (inputRef.value) {
		inputRef.value.focus();
	}
});

// Handle terminal click to refocus input
const handleTerminalClick = () => {
	if (inputRef.value) {
		inputRef.value.focus();
	}
};
</script>

<template>
	<div class="repl-container">
		<div class="github-container">
			<a
				href="https://github.com"
				target="_blank"
				rel="noopener"
				class="github-link"
			>
				<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
					<path
						d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
					/>
				</svg>
			</a>
		</div>
		<div class="terminal-wrapper">
			<div class="terminal" ref="terminalRef" @click="handleTerminalClick">
				<div class="terminal-content">
					<span class="output">BangScript REPL v1.0.0</span>
					<!-- <div v-for="(line, index) in terminalLines" :key="index" class="terminal-line">
            <span v-if="line.startsWith('>>>')" class="prompt-primary">{{ line }}</span>
            <span v-else-if="line.startsWith('...')" class="prompt-secondary">{{ line }}</span>
            <span v-else class="output">{{ line }}</span>
          </div> -->
					<div class="current-input-line">
						<span class="prompt-primary">{{ inputIndicator }}</span>
						<span class="input-text">{{ currentInput }}</span>
						<span class="cursor">█</span>
					</div>
				</div>
			</div>
			<input
				ref="inputRef"
				v-model="currentInput"
				@keydown="handleKeydown"
				@input="handleInput"
				class="hidden-input"
				autocomplete="off"
				autocorrect="off"
				autocapitalize="off"
				spellcheck="false"
			/>
		</div>
	</div>
</template>

<style scoped>
.repl-container {
	display: flex;
	flex-direction: column;
	height: 100vh;
	background: linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 100%);
	font-family: "JetBrains Mono", monospace;
}

.github-container {
	position: absolute;
	top: 1rem;
	right: 1rem;
	z-index: 10;
}

.github-link {
	color: #ffffff;
	transition: all 0.3s ease;
	padding: 0.5rem;
	border-radius: 0.5rem;
	display: block;
}

.github-link:hover {
	transform: scale(1.2);
	color: #00ffff;
}

.terminal-wrapper {
	flex: 1;
	display: flex;
	flex-direction: column;
	padding: 2rem;
	position: relative;
}

.terminal {
	flex: 1;
	background: rgba(10, 10, 10, 0.9);
	border: 1px solid rgba(0, 255, 255, 0.3);
	border-radius: 1rem;
	padding: 1.5rem;
	overflow-y: auto;
	position: relative;
	box-shadow:
		0 0 50px rgba(0, 255, 255, 0.2),
		inset 0 0 20px rgba(0, 255, 255, 0.05);
}

.terminal::before {
	content: "";
	position: absolute;
	top: -2px;
	left: -2px;
	right: -2px;
	bottom: -2px;
	background: linear-gradient(45deg, #00ffff, #00cccc, #00ffff);
	border-radius: 1rem;
	z-index: -1;
	opacity: 0.6;
	animation: borderGlow 3s ease-in-out infinite;
}

@keyframes borderGlow {
	0%,
	100% {
		opacity: 0.4;
	}
	50% {
		opacity: 0.8;
	}
}

.terminal-content {
	min-height: 100%;
}

.terminal-line {
	margin-bottom: 0.5rem;
	line-height: 1.6;
	font-size: 0.95rem;
}

.current-input-line {
	display: flex;
	align-items: center;
	line-height: 1.6;
	font-size: 0.95rem;
}

.prompt-primary {
	color: #00ffff;
	font-weight: 600;
	margin-right: 0.5rem;
	text-shadow: 0 0 10px rgba(0, 255, 255, 0.5);
}

.prompt-secondary {
	color: #00cccc;
	font-weight: 600;
	margin-right: 0.5rem;
	text-shadow: 0 0 10px rgba(0, 255, 255, 0.5);
}

.output {
	color: #ffffff;
	margin-left: 0.5rem;
}

.input-text {
	color: #ffffff;
}

.cursor {
	color: #00ffff;
	font-weight: bold;
	animation: cursorBlink 1s ease-in-out infinite;
	text-shadow: 0 0 10px rgba(0, 255, 255, 0.8);
}

@keyframes cursorBlink {
	0%,
	100% {
		opacity: 1;
	}
	50% {
		opacity: 0.3;
	}
}

.hidden-input {
	position: absolute;
	left: -9999px;
	opacity: 0;
}

.terminal::-webkit-scrollbar {
	display: none;
}

.terminal {
	-ms-overflow-style: none;
	scrollbar-width: none;
}

@media (max-width: 768px) {
	.terminal-wrapper {
		padding: 1rem;
	}

	.terminal {
		padding: 1rem;
	}

	.github-container {
		top: 0.5rem;
		right: 0.5rem;
	}

	.github-link svg {
		width: 20px;
		height: 20px;
	}

	.terminal-line,
	.current-input-line {
		font-size: 0.85rem;
	}
}
</style>
