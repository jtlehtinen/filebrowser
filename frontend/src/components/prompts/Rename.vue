<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t("prompts.rename") }}</h2>
    </div>

    <div class="card-content">
      <p>
        {{ $t("prompts.renameMessage") }} <code>{{ oldName() }}</code
        >:
      </p>
      <input
        id="focus-prompt"
        class="input input--block"
        type="text"
        @keyup.enter="submit"
        v-model.trim="name"
      />
    </div>

    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="closeHovers"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')"
      >
        {{ $t("buttons.cancel") }}
      </button>
      <button
        @click="submit"
        class="button button--flat"
        type="submit"
        :aria-label="$t('buttons.rename')"
        :title="$t('buttons.rename')"
      >
        {{ $t("buttons.rename") }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import url from "@/utils/url";
import { files as api } from "@/api";
import { removePrefix } from "@/api/utils";

const $showError = inject<(e: unknown) => void>("$showError");

const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const { req, selected, selectedCount, isListing } = storeToRefs(fileStore);
const { reload, preselect } = storeToRefs(fileStore);
const { closeHovers } = layoutStore;

const name = ref("");

/**
 * Returns the old name for the item being renamed.
 */
const oldName = computed(() => {
  if (!isListing.value) {
    return req.value!.name;
  }
  if (selectedCount.value === 0 || selectedCount.value > 1) {
    return;
  }
  return req.value!.items[selected.value[0]].name;
});

onMounted(() => {
  // @ts-expect-error Deal with this later.
  name.value = oldName.value;
});

/**
 * Submits the rename request.
 */
async function submit() {
  let oldLink = "";
  let newLink = "";

  if (!isListing.value) {
    oldLink = req.value!.url;
  } else {
    oldLink = req.value!.items[selected.value[0]].url;
  }

  newLink = url.removeLastDir(oldLink) + "/" + encodeURIComponent(name.value);

  try {
    await api.move([{ from: oldLink, to: newLink }]);
    if (!isListing.value) {
      // @ts-expect-error: $router is injected by Vue Router

      (getCurrentInstance()?.proxy as any).$router.push({ path: newLink });
      return;
    }
    preselect.value = removePrefix(newLink);
    reload.value = true;
  } catch (e) {
    $showError?.(e);
  }

  closeHovers();
}
</script>
