<template>
  <div v-show="active" @click="closeHovers" class="overlay"></div>
  <nav :class="{ active }">
    <template v-if="isLoggedIn">
      <button @click="toAccountSettings" class="action">
        <i class="material-icons">person</i>
        <span>{{ user?.username }}</span>
      </button>
      <button
        class="action"
        @click="toRoot"
        :aria-label="$t('sidebar.myFiles')"
        :title="$t('sidebar.myFiles')"
      >
        <i class="material-icons">folder</i>
        <span>{{ $t("sidebar.myFiles") }}</span>
      </button>

      <div v-if="user?.perm.create">
        <button
          @click="showHover('newDir')"
          class="action"
          :aria-label="$t('sidebar.newFolder')"
          :title="$t('sidebar.newFolder')"
        >
          <i class="material-icons">create_new_folder</i>
          <span>{{ $t("sidebar.newFolder") }}</span>
        </button>

        <button
          @click="showHover('newFile')"
          class="action"
          :aria-label="$t('sidebar.newFile')"
          :title="$t('sidebar.newFile')"
        >
          <i class="material-icons">note_add</i>
          <span>{{ $t("sidebar.newFile") }}</span>
        </button>
      </div>

      <div v-if="user?.perm.admin">
        <button
          class="action"
          @click="toGlobalSettings"
          :aria-label="$t('sidebar.settings')"
          :title="$t('sidebar.settings')"
        >
          <i class="material-icons">settings_applications</i>
          <span>{{ $t("sidebar.settings") }}</span>
        </button>
      </div>
      <button
        v-if="canLogout"
        @click="auth.logout()"
        class="action"
        id="logout"
        :aria-label="$t('sidebar.logout')"
        :title="$t('sidebar.logout')"
      >
        <i class="material-icons">exit_to_app</i>
        <span>{{ $t("sidebar.logout") }}</span>
      </button>
    </template>
    <template v-else>
      <router-link
        class="action"
        to="/login"
        :aria-label="$t('sidebar.login')"
        :title="$t('sidebar.login')"
      >
        <i class="material-icons">exit_to_app</i>
        <span>{{ $t("sidebar.login") }}</span>
      </router-link>

      <router-link
        v-if="signup"
        class="action"
        to="/login"
        :aria-label="$t('sidebar.signup')"
        :title="$t('sidebar.signup')"
      >
        <i class="material-icons">person_add</i>
        <span>{{ $t("sidebar.signup") }}</span>
      </router-link>
    </template>

    <div
      class="credits"
      v-if="isFiles && !disableUsedPercentage"
      style="width: 90%; margin: 2em 2.5em 3em 2.5em"
    >
      <ProgressBar :val="usage.usedPercentage" size="small"></ProgressBar>
      <br />
      {{ usage.used }} of {{ usage.total }} used
    </div>

    <p class="credits">
      <span>
        <span v-if="disableExternal">File Browser</span>
        <a
          v-else
          rel="noopener noreferrer"
          target="_blank"
          href="https://github.com/filebrowser/filebrowser"
          >File Browser</a
        >
        <span> {{ " " }} {{ version }}</span>
      </span>
      <span>
        <a @click="help">{{ $t("sidebar.help") }}</a>
      </span>
    </p>
  </nav>
</template>

<script setup lang="ts">
import { computed, reactive, inject, watch, onUnmounted } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import * as auth from "@/utils/auth";
import {
  version as VERSION,
  signup as SIGNUP,
  disableExternal as DISABLE_EXTERNAL,
  disableUsedPercentage as DISABLE_USED_PERCENTAGE,
  noAuth as NO_AUTH,
  loginPage as LOGIN_PAGE,
} from "@/utils/constants";
import { files as api } from "@/api";
import ProgressBar from "@/components/ProgressBar.vue";
import prettyBytes from "pretty-bytes";

const USAGE_DEFAULT = { used: "0 B", total: "0 B", usedPercentage: 0 };

const usage = reactive({ ...USAGE_DEFAULT });
let usageAbortController = new AbortController();

const $showError = inject<(msg: string) => void>("$showError");

const authStore = useAuthStore();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const user = computed(() => authStore.user);
const isLoggedIn = computed(() => authStore.isLoggedIn);
const isFiles = computed(() => fileStore.isFiles);
const currentPromptName = computed(() => layoutStore.currentPromptName);

const active = computed(() => currentPromptName.value === "sidebar");
const signup = computed(() => SIGNUP);
const version = computed(() => VERSION);
const disableExternal = computed(() => DISABLE_EXTERNAL);
const disableUsedPercentage = computed(() => DISABLE_USED_PERCENTAGE);
const canLogout = computed(() => !NO_AUTH && LOGIN_PAGE);

const closeHovers = layoutStore.closeHovers;
const showHover = layoutStore.showHover;

/**
 * Abort any ongoing usage fetch.
 */
function abortOngoingFetchUsage() {
  usageAbortController.abort();
}

/**
 * Fetch usage statistics for the current path.
 */
async function fetchUsage() {
  const path = window.location.pathname.endsWith("/")
    ? window.location.pathname
    : window.location.pathname + "/";
  let usageStats = USAGE_DEFAULT;
  if (disableUsedPercentage.value) {
    Object.assign(usage, usageStats);
    return;
  }
  try {
    abortOngoingFetchUsage();
    usageAbortController = new AbortController();
    const result = await api.usage(path, usageAbortController.signal);
    usageStats = {
      used: prettyBytes(result.used, { binary: true }),
      total: prettyBytes(result.total, { binary: true }),
      usedPercentage: Math.round((result.used / result.total) * 100),
    };
  } catch (err) {
    if ($showError)
      $showError("Failed to fetch usage: " + (err as Error).message);
  } finally {
    Object.assign(usage, usageStats);
  }
}

/**
 * Navigate to the root files page.
 */
function toRoot() {
  layoutStore.closeHovers();
  window.location.assign("/files");
}

/**
 * Navigate to the account settings page.
 */
function toAccountSettings() {
  layoutStore.closeHovers();
  window.location.assign("/settings/profile");
}

/**
 * Navigate to the global settings page.
 */
function toGlobalSettings() {
  layoutStore.closeHovers();
  window.location.assign("/settings/global");
}

/**
 * Show help hover.
 */
function help() {
  showHover("help");
}

watch(
  () => window.location.pathname,
  (to) => {
    if (to.includes("/files")) {
      fetchUsage();
    }
  },
  { immediate: true }
);

onUnmounted(() => {
  abortOngoingFetchUsage();
});
</script>
