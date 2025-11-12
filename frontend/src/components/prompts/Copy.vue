<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t("prompts.copy") }}</h2>
    </div>

    <div class="card-content">
      <p>{{ $t("prompts.copyMessage") }}</p>
      <FileList
        ref="fileList"
        @update:selected="(val) => (dest = val)"
        tabindex="1"
      />
    </div>

    <div
      class="card-action"
      style="display: flex; align-items: center; justify-content: space-between"
    >
      <template v-if="user?.perm.create">
        <button
          class="button button--flat"
          @click="fileList.createDir()"
          :aria-label="$t('sidebar.newFolder')"
          :title="$t('sidebar.newFolder')"
          style="justify-self: left"
        >
          <span>{{ $t("sidebar.newFolder") }}</span>
        </button>
      </template>
      <div>
        <button
          class="button button--flat button--grey"
          @click="closeHovers"
          :aria-label="$t('buttons.cancel')"
          :title="$t('buttons.cancel')"
          tabindex="3"
        >
          {{ $t("buttons.cancel") }}
        </button>
        <button
          id="focus-prompt"
          class="button button--flat"
          @click="copy"
          :aria-label="$t('buttons.copy')"
          :title="$t('buttons.copy')"
          tabindex="2"
        >
          {{ $t("buttons.copy") }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, inject } from "vue";
import { useRouter, useRoute } from "vue-router";
import { storeToRefs } from "pinia";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { useAuthStore } from "@/stores/auth";
import FileList from "./FileList.vue";
import { files as api } from "@/api";
import buttons from "@/utils/buttons";
import * as upload from "@/utils/upload";
import { removePrefix } from "@/api/utils";

const fileList = ref<InstanceType<typeof FileList> | null>(null);
const dest = ref<string | null>(null);

const { req, selected, reload, preselect } = storeToRefs(useFileStore());
const { user } = storeToRefs(useAuthStore());
const layoutStore = useLayoutStore();
const router = useRouter();
const route = useRoute();

const showHover = layoutStore.showHover;
const closeHovers = layoutStore.closeHovers;

const $showError = inject<(e: unknown) => void>("$showError")!;

/**
 * Copies selected files to the destination directory.
 * Handles conflicts and user prompts for overwrite/rename.
 */
async function copy(event: Event) {
  event.preventDefault();
  if (!dest.value) return;

  const items = selected.value.map((item) => ({
    from: req.value!.items[item].url,
    to: dest.value + encodeURIComponent(req.value!.items[item].name),
    name: req.value!.items[item].name,
  }));

  async function action(overwrite: boolean, rename: boolean) {
    try {
      buttons.loading("copy");
      await api.copy(items, overwrite, rename);
      buttons.success("copy");
      preselect.value = removePrefix(items[0].to);
      if (route.path === dest.value) {
        reload.value = true;
        return;
      }
      // @ts-expect-error Deal with this later.
      router.push({ path: dest.value });
    } catch (e) {
      buttons.done("copy");
      $showError(e);
    }
  }

  if (route.path === dest.value) {
    closeHovers();
    action(false, true);
    return;
  }

  try {
    const dstItems = (await api.fetch(dest.value)).items;
    // @ts-expect-error Deal with this later.
    const conflict = upload.checkConflict(items, dstItems);
    let overwrite = false;
    let rename = false;

    if (conflict) {
      showHover({
        prompt: "replace-rename",
        confirm: (event: Event, option: string) => {
          overwrite = option === "overwrite";
          rename = option === "rename";
          event.preventDefault();
          closeHovers();
          action(overwrite, rename);
        },
      });
      return;
    }
    action(overwrite, rename);
  } catch (e) {
    $showError(e);
  }
}
</script>
