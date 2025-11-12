<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t("prompts.move") }}</h2>
    </div>

    <div class="card-content">
      <file-list
        ref="fileList"
        @update:selected="(val) => (dest = val)"
        :exclude="excludedFolders"
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
          @click="move"
          :disabled="$route.path === dest"
          :aria-label="$t('buttons.move')"
          :title="$t('buttons.move')"
          tabindex="2"
        >
          {{ $t("buttons.move") }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, toRefs } from "vue";
import { useRouter } from "vue-router";
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

const fileStore = useFileStore();
const layoutStore = useLayoutStore();
const authStore = useAuthStore();
const router = useRouter();

const { req, selected } = toRefs(fileStore);
const { user } = toRefs(authStore);
const showHover = layoutStore.showHover;
const closeHovers = layoutStore.closeHovers;

const $showError = inject<(e: unknown) => void>("$showError");

const excludedFolders = computed(() =>
  selected.value
    .filter((idx) => req.value!.items[idx].isDir)
    .map((idx) => req.value!.items[idx].url)
);

/**
 * Move selected files/folders to destination.
 * Handles conflicts and user confirmation.
 */
async function move(event: Event) {
  event.preventDefault();
  try {
    const items = selected.value.map((item) => ({
      from: req.value!.items[item].url,
      to: dest.value + encodeURIComponent(req.value!.items[item].name),
      name: req.value!.items[item].name,
    }));

    async function action(overwrite: boolean, rename: boolean) {
      buttons.loading("move");
      try {
        await api.move(items, overwrite, rename);
        buttons.success("move");
        fileStore.preselect = removePrefix(items[0].to);
        await router.push({ path: dest.value! });
      } catch (e) {
        buttons.done("move");
        $showError?.(e);
      }
    }

    const dstItems = (await api.fetch(dest.value!)).items;
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

    await action(overwrite, rename);
  } catch (e) {
    $showError?.(e);
  }
}
</script>
