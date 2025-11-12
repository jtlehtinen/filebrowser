<template>
  <div
    class="shell"
    :class="{ ['shell--hidden']: !showShell }"
    :style="{ height: `${shellHeight}em`, direction: 'ltr' }"
  >
    <div
      @pointerdown="startDrag()"
      @pointerup="stopDrag()"
      class="shell__divider"
      :style="shellDrag ? { background: `${checkTheme()}` } : ''"
    ></div>
    <div @click="focus" class="shell__content" ref="scrollable">
      <div v-for="(c, index) in content" :key="index" class="shell__result">
        <div class="shell__prompt">
          <i class="material-icons">chevron_right</i>
        </div>
        <pre class="shell__text">{{ c.text }}</pre>
      </div>

      <div
        class="shell__result"
        :class="{ 'shell__result--hidden': !canInput }"
      >
        <div class="shell__prompt">
          <i class="material-icons">chevron_right</i>
        </div>
        <pre
          tabindex="0"
          ref="input"
          class="shell__text"
          :contenteditable="true"
          @keydown.prevent.arrow-up="historyUp"
          @keydown.prevent.arrow-down="historyDown"
          @keypress.prevent.enter="submit"
        />
      </div>
    </div>
    <div
      @pointerup="stopDrag()"
      class="shell__overlay"
      v-show="shellDrag"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { commands } from "@/api";
import { throttle } from "lodash-es";
import { theme } from "@/utils/constants";
import { useRoute } from "vue-router";

const layoutStore = useLayoutStore();
const fileStore = useFileStore();
const route = useRoute();

const showShell = computed(() => layoutStore.showShell);
const isFiles = computed(() => fileStore.isFiles);
const path = computed(() => (isFiles.value ? route.path : ""));

const content = ref<{ text: string }[]>([]);
const history = ref<string[]>([]);
const historyPos = ref(0);
const canInput = ref(true);
const shellDrag = ref(false);
const shellHeight = ref(25);
const fontsize = ref(
  parseFloat(getComputedStyle(document.documentElement).fontSize)
);

const scrollable = ref<HTMLElement | null>(null);
const input = ref<HTMLElement | null>(null);

function checkTheme(): string {
  return theme === "dark"
    ? "rgba(255, 255, 255, 0.4)"
    : "rgba(127, 127, 127, 0.4)";
}

function startDrag() {
  document.addEventListener("pointermove", handleDrag);
  shellDrag.value = true;
}

function stopDrag() {
  document.removeEventListener("pointermove", handleDrag);
  shellDrag.value = false;
}

const handleDrag = throttle(function (event: PointerEvent) {
  const top = window.innerHeight / fontsize.value - 4;
  const divider = document.querySelector(".shell__divider") as HTMLElement;
  const userPos = (window.innerHeight - event.clientY) / fontsize.value;
  const bottom = 2.25 + (divider?.offsetHeight ?? 0) / fontsize.value;

  if (userPos <= top && userPos >= bottom) {
    shellHeight.value = Number(userPos.toFixed(2));
  }
}, 32);

const resize = throttle(function () {
  const top = window.innerHeight / fontsize.value - 4;
  const divider = document.querySelector(".shell__divider") as HTMLElement;
  const bottom = 2.25 + (divider?.offsetHeight ?? 0) / fontsize.value;

  if (shellHeight.value > top) {
    shellHeight.value = top;
  } else if (shellHeight.value < bottom) {
    shellHeight.value = bottom;
  }
}, 32);

function scroll() {
  if (scrollable.value) {
    scrollable.value.scrollTop = scrollable.value.scrollHeight;
  }
}

function focus() {
  input.value?.focus();
}

function historyUp() {
  if (historyPos.value > 0) {
    if (input.value) input.value.innerText = history.value[--historyPos.value];
    focus();
  }
}

function historyDown() {
  if (historyPos.value >= 0 && historyPos.value < history.value.length - 1) {
    if (input.value) input.value.innerText = history.value[++historyPos.value];
    focus();
  } else {
    historyPos.value = history.value.length;
    if (input.value) input.value.innerText = "";
  }
}

function submit(event: KeyboardEvent) {
  const target = event.target as HTMLElement;
  const cmd = target.innerText.trim();

  if (cmd === "") return;

  if (cmd === "clear") {
    content.value = [];
    target.innerHTML = "";
    return;
  }

  if (cmd === "exit") {
    target.innerHTML = "";
    layoutStore.toggleShell();
    return;
  }

  canInput.value = false;
  target.innerHTML = "";

  const results = { text: `${cmd}\n\n` };
  history.value.push(cmd);
  historyPos.value = history.value.length;
  content.value.push(results);

  commands(
    path.value,
    cmd,
    (event: MessageEvent) => {
      results.text += `${event.data}\n`;
      scroll();
    },
    () => {
      results.text = results.text.replace(/\u001b\[[0-9;]+m/g, "").trimEnd();
      canInput.value = true;
      input.value?.focus();
      scroll();
    }
  );
}

onMounted(() => {
  window.addEventListener("resize", resize);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", resize);
});
</script>
