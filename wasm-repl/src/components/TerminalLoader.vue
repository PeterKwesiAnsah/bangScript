<template>
  <div class="loading-screen">
    <div class="terminal-loading">
      <div class="terminal-header">
        <div class="terminal-buttons">
          <div class="button close"></div>
          <div class="button minimize"></div>
          <div class="button maximize"></div>
        </div>
      </div>
      <div class="terminal-content">
        <div class="loading-line">
          <span class="prompt">$</span>
          <span class="command">Just a moment...</span>
          <span class="cursor">█</span>
        </div>
        <div class="loading-output">
          <div class="output-line">Loading WASM module...</div>
          <div class="output-line">Initializing runtime...</div>
          <div class="output-line">Preparing interpreter...</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.loading-screen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.terminal-loading {
  width: 600px;
  max-width: 90vw;
  background: rgba(10, 10, 10, 0.9);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 0.5rem;
  overflow: hidden;
  box-shadow: 0 0 30px rgba(0, 255, 255, 0.2);
}

.terminal-header {
  background: rgba(0, 255, 255, 0.1);
  padding: 0.75rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.terminal-buttons {
  display: flex;
  gap: 0.5rem;
}

.button {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.button.close {
  background: #ff5f56;
}

.button.minimize {
  background: #ffbd2e;
}

.button.maximize {
  background: #27c93f;
}

.terminal-content {
  padding: 1.5rem;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.95rem;
  line-height: 1.6;
}

.loading-line {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.prompt {
  color: #00ffff;
  font-weight: 600;
  margin-right: 0.5rem;
  text-shadow: 0 0 10px rgba(0, 255, 255, 0.5);
}

.command {
  color: #ffffff;
  margin-right: 0.5rem;
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

.loading-output {
  color: #cccccc;
}

.output-line {
  margin-bottom: 0.5rem;
  opacity: 0;
  animation: fadeInUp 0.5s ease-out forwards;
}

.output-line:nth-child(1) {
  animation-delay: 0.5s;
}

.output-line:nth-child(2) {
  animation-delay: 1s;
}

.output-line:nth-child(3) {
  animation-delay: 1.5s;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
