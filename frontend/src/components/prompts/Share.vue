<template>
  <div class="card floating" id="share">
    <div class="card-title">
      <h2>{{ $t("buttons.share") }}</h2>
    </div>

    <template v-if="listing">
      <div class="card-content">
        <table>
          <tr>
            <th>#</th>
            <th>{{ $t("settings.shareDuration") }}</th>
            <th></th>
            <th></th>
          </tr>

          <tr v-for="link in links" :key="link.hash">
            <td>{{ link.hash }}</td>
            <td>
              <template v-if="link.expire !== 0">{{
                humanTime(link.expire)
              }}</template>
              <template v-else>{{ $t("permanent") }}</template>
            </td>
            <td class="small">
              <button
                class="action copy-clipboard"
                :aria-label="$t('buttons.copyToClipboard')"
                :title="$t('buttons.copyToClipboard')"
                @click="copyToClipboard(buildLink(link))"
              >
                <i class="material-icons">content_paste</i>
              </button>
            </td>
            <td class="small">
              <button
                class="action"
                @click="deleteLink($event, link)"
                :aria-label="$t('buttons.delete')"
                :title="$t('buttons.delete')"
              >
                <i class="material-icons">delete</i>
              </button>
            </td>
          </tr>
        </table>
      </div>

      <div class="card-action">
        <button
          class="button button--flat button--grey"
          @click="closeHovers"
          :aria-label="$t('buttons.close')"
          :title="$t('buttons.close')"
          tabindex="2"
        >
          {{ $t("buttons.close") }}
        </button>
        <button
          id="focus-prompt"
          class="button button--flat button--blue"
          @click="() => switchListing()"
          :aria-label="$t('buttons.new')"
          :title="$t('buttons.new')"
          tabindex="1"
        >
          {{ $t("buttons.new") }}
        </button>
      </div>
    </template>

    <template v-else>
      <div class="card-content">
        <p>{{ $t("settings.shareDuration") }}</p>
        <div class="input-group input">
          <vue-number-input
            center
            controls
            size="small"
            :max="2147483647"
            :min="0"
            @keyup.enter="submit"
            v-model="time"
            tabindex="1"
          />
          <select
            class="right"
            v-model="unit"
            :aria-label="$t('time.unit')"
            tabindex="2"
          >
            <option value="seconds">{{ $t("time.seconds") }}</option>
            <option value="minutes">{{ $t("time.minutes") }}</option>
            <option value="hours">{{ $t("time.hours") }}</option>
            <option value="days">{{ $t("time.days") }}</option>
          </select>
        </div>
        <p>{{ $t("prompts.optionalPassword") }}</p>
        <input
          class="input input--block"
          type="password"
          v-model.trim="password"
          tabindex="3"
        />
      </div>

      <div class="card-action">
        <button
          class="button button--flat button--grey"
          @click="() => switchListing()"
          :aria-label="$t('buttons.cancel')"
          :title="$t('buttons.cancel')"
          tabindex="5"
        >
          {{ $t("buttons.cancel") }}
        </button>
        <button
          id="focus-prompt"
          class="button button--flat button--blue"
          @click="submit"
          :aria-label="$t('buttons.share')"
          :title="$t('buttons.share')"
          tabindex="4"
        >
          {{ $t("buttons.share") }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router";
import { useI18n } from "vue-i18n";
import { ref, computed, inject, onBeforeMount } from "vue";
import { storeToRefs } from "pinia";
import { useFileStore } from "@/stores/file";
import { share as api } from "@/api";
import dayjs from "dayjs";
import { useLayoutStore } from "@/stores/layout";
import { copy } from "@/utils/clipboard";

const route = useRoute();
const { t } = useI18n();

const $showError = inject<(e: unknown) => void>("$showError")!;
const $showSuccess = inject<(msg: string) => void>("$showSuccess")!;

const fileStore = useFileStore();
const layoutStore = useLayoutStore();
const { req, selected, selectedCount, isListing } = storeToRefs(fileStore);
const { closeHovers } = layoutStore;

const time = ref(0);
const unit = ref("hours");
const links = ref<any[]>([]);
const password = ref("");
const listing = ref(true);

const url = computed(() => {
  if (!isListing.value) {
    return route.path;
  }
  if (selectedCount.value === 0 || selectedCount.value > 1) {
    return;
  }
  return req.value!.items[selected.value[0]].url;
});

/**
 * Sorts the links array by expiration.
 */
function sort() {
  links.value = [...links.value].sort((a, b) => {
    if (a.expire === 0) return -1;
    if (b.expire === 0) return 1;
    return new Date(a.expire).getTime() - new Date(b.expire).getTime();
  });
}

/**
 * Converts a unix timestamp to a human readable string.
 * @param time - Unix timestamp in seconds
 */
function humanTime(time: number) {
  return dayjs(time * 1000).fromNow();
}

/**
 * Builds a share URL from a share object.
 * @param share - Share object
 */
function buildLink(share: any) {
  return api.getShareURL(share);
}

/**
 * Copies text to clipboard and shows a success or error message.
 * @param text - Text to copy
 */
async function copyToClipboard(text: string) {
  try {
    await copy({ text });
    $showSuccess(t("success.linkCopied"));
  } catch {
    try {
      await copy({ text }, { permission: true });
      $showSuccess(t("success.linkCopied"));
    } catch (e) {
      $showError(e);
    }
  }
}

/**
 * Submits a new share link.
 */
async function submit() {
  try {
    let res = null;
    if (!time.value) {
      res = await api.create(url.value!, password.value);
    } else {
      res = await api.create(
        url.value!,
        password.value,
        // @ts-expect-error Deal with this later.
        time.value,
        unit.value
      );
    }
    links.value.push(res);
    sort();
    time.value = 0;
    unit.value = "hours";
    password.value = "";
    listing.value = true;
  } catch (e) {
    $showError(e);
  }
}

/**
 * Deletes a share link.
 * @param event - Mouse event
 * @param link - Link object
 */
async function deleteLink(event: Event, link: any) {
  event.preventDefault();
  try {
    await api.remove(link.hash);
    links.value = links.value.filter((item) => item.hash !== link.hash);
    if (links.value.length === 0) {
      listing.value = false;
    }
  } catch (e) {
    $showError(e);
  }
}

/**
 * Switches between listing and creation mode.
 */
function switchListing() {
  if (links.value.length === 0 && !listing.value) {
    closeHovers();
  }
  listing.value = !listing.value;
}

onBeforeMount(async () => {
  try {
    const fetchedLinks = await api.get(url.value!);
    // @ts-expect-error Deal with this later.
    links.value = fetchedLinks;
    sort();
    if (links.value.length === 0) {
      listing.value = false;
    }
  } catch (e) {
    $showError(e);
  }
});
</script>
