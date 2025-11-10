<template>
  <div>
    <h3>{{ $t("settings.userCommands") }}</h3>
    <p class="small">
      {{ $t("settings.userCommandsHelp") }} <i>git svn hg</i>.
    </p>
    <input class="input input--block" type="text" v-model.trim="raw" />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{ commands: string[] }>();
const emit = defineEmits<{ (e: "update:commands", value: string[]): void }>();

const raw = computed({
  get() {
    return props.commands.join(" ");
  },
  set(value: string) {
    emit("update:commands", value !== "" ? value.split(" ") : []);
  },
});
</script>
