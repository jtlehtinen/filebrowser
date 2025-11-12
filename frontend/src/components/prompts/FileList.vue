<template>
  <div>
    <ul class="file-list">
      <li
        @click="itemClick"
        @touchstart="touchstart"
        @dblclick="next"
        role="button"
        tabindex="0"
        :aria-label="item.name"
        :aria-selected="selected == item.url"
        :key="item.name"
        v-for="item in items"
        :data-url="item.url"
      >
        {{ item.name }}
      </li>
    </ul>

    <p>
      {{ $t("prompts.currentlyNavigating") }} <code>{{ nav }}</code
      >.
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, inject } from "vue";
import { storeToRefs } from "pinia";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import url from "@/utils/url";
import { files } from "@/api";
import { StatusError } from "@/api/utils.js";

const props = defineProps<{ exclude?: string[] }>();
const emit = defineEmits<{ "update:selected": [value: string] }>();
const $showError = inject<(e: unknown) => void>("$showError");

const { user } = storeToRefs(useAuthStore());
const { req } = storeToRefs(useFileStore());
const { showHover } = useLayoutStore();

const items = ref<{ name: string; url: string }[]>([]);
const touches = ref({ id: "", count: 0 });
const selected = ref<string | null>(null);
const current = ref(window.location.pathname);
let nextAbortController = new AbortController();

const nav = computed(() => decodeURIComponent(current.value));

function abortOngoingNext() {
  nextAbortController.abort();
}

function fillOptions(request: {
  url: string;
  items: null | { isDir: boolean; url: string; name: string }[];
}) {
  current.value = request.url;
  items.value = [];

  emit("update:selected", current.value);

  if (request.url !== "/files/") {
    items.value.push({
      name: "..",
      url: url.removeLastDir(request.url) + "/",
    });
  }

  if (request.items === null) return;

  for (const item of request.items) {
    if (!item.isDir) continue;
    if (props.exclude?.includes(item.url)) continue;

    items.value.push({
      name: item.name,
      url: item.url,
    });
  }
}

async function next(event: Event) {
  const target = event.currentTarget as HTMLElement;
  const uri = target.dataset.url;
  if (!uri) return;
  abortOngoingNext();
  nextAbortController = new AbortController();
  try {
    const response = await files.fetch(uri, nextAbortController.signal);
    fillOptions(response);
  } catch (e) {
    if (e instanceof StatusError && e.is_canceled) {
      return;
    }
    $showError?.(e);
  }
}

function touchstart(event: TouchEvent | MouseEvent) {
  const target = event.currentTarget as HTMLElement;
  const url = target.dataset.url;
  if (!url) return;

  setTimeout(() => {
    touches.value.count = 0;
  }, 300);

  if (touches.value.id !== url) {
    touches.value.id = url;
    touches.value.count = 1;
    return;
  }

  touches.value.count++;

  if (touches.value.count > 1) {
    next(event);
  }
}

function itemClick(event: MouseEvent) {
  if (user.value?.singleClick) next(event);
  else select(event);
}

function select(event: MouseEvent) {
  const target = event.currentTarget as HTMLElement;
  const url = target.dataset.url;
  if (!url) return;

  if (selected.value === url) {
    selected.value = null;
    emit("update:selected", current.value);
    return;
  }

  selected.value = url;
  emit("update:selected", selected.value);
}

/**
 * Shows the new directory creation prompt.
 */
async function createDir() {
  try {
    showHover({
      prompt: "newDir",
      // @ts-expect-error Deal with this later.
      action: null,
      confirm: null,
      props: {
        redirect: false,
        base:
          current.value === (window as any).$route?.path ? null : current.value,
      },
    });
  } catch (e) {
    console.error("createDir error:", e);
  }
}

onMounted(() => {
  fillOptions(req.value!);
});

onUnmounted(() => {
  abortOngoingNext();
});

defineExpose({ createDir });
</script>
