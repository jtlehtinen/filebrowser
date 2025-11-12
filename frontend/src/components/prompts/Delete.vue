<template>
  <div class="card floating">
    <div class="card-content">
      <p v-if="!isListing || selectedCount === 1">
        {{ $t("prompts.deleteMessageSingle") }}
      </p>
      <p v-else>
        {{ $t("prompts.deleteMessageMultiple", { count: selectedCount }) }}
      </p>
    </div>
    <div class="card-action">
      <button
        @click="closeHovers"
        class="button button--flat button--grey"
        :aria-label="$t('buttons.cancel')"
        :title="$t('buttons.cancel')"
        tabindex="2"
      >
        {{ $t("buttons.cancel") }}
      </button>
      <button
        id="focus-prompt"
        @click="submit"
        class="button button--flat button--red"
        :aria-label="$t('buttons.delete')"
        :title="$t('buttons.delete')"
        tabindex="1"
      >
        {{ $t("buttons.delete") }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from "vue";
import { storeToRefs } from "pinia";
import { files as api } from "@/api";
import buttons from "@/utils/buttons";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

const $showError = inject<(e: unknown) => void>("$showError");

const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const { isListing, selectedCount, req, selected } = storeToRefs(fileStore);
const { currentPrompt } = storeToRefs(layoutStore);
const { reload, preselect } = storeToRefs(fileStore);

const closeHovers = layoutStore.closeHovers;

/**
 * Handles the delete action for single or multiple items.
 */
async function submit() {
  buttons.loading("delete");
  try {
    if (!isListing.value) {
      await api.remove((<any>window).$router?.currentRoute?.value?.path ?? "");
      buttons.success("delete");
      currentPrompt.value?.confirm();
      closeHovers();
      return;
    }

    closeHovers();

    if (selectedCount.value === 0) {
      return;
    }

    const promises: Promise<unknown>[] = [];
    for (const index of selected.value) {
      promises.push(api.remove(req.value!.items[index].url));
    }

    await Promise.all(promises);
    buttons.success("delete");

    const nearbyItem =
      req.value!.items[Math.max(0, Math.min(...selected.value) - 1)];
    preselect.value = nearbyItem?.path;
    reload.value = true;
  } catch (e) {
    buttons.done("delete");
    $showError?.(e);
    if (isListing.value) reload.value = true;
  }
}
</script>
