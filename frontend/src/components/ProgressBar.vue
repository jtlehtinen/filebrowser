<!-- This component taken directly from vue-simple-progress
since it didnt support Vue 3 but the component itself does
https://raw.githubusercontent.com/dzwillia/vue-simple-progress/master/src/components/Progress.vue -->
<template>
  <div>
    <div
      class="vue-simple-progress-text"
      :style="textStyle"
      v-if="text.length > 0 && textPosition == 'top'"
    >
      {{ text }}
    </div>
    <div class="vue-simple-progress" :style="progressStyle">
      <div
        class="vue-simple-progress-text"
        :style="textStyle"
        v-if="text.length > 0 && textPosition == 'middle'"
      >
        {{ text }}
      </div>
      <div
        style="position: relative; left: -9999px"
        :style="textStyle"
        v-if="text.length > 0 && textPosition == 'inside'"
      >
        {{ text }}
      </div>
      <div class="vue-simple-progress-bar" :style="barStyle">
        <div
          :style="textStyle"
          v-if="text.length > 0 && textPosition == 'inside'"
        >
          {{ text }}
        </div>
      </div>
    </div>
    <div
      class="vue-simple-progress-text"
      :style="textStyle"
      v-if="text.length > 0 && textPosition == 'bottom'"
    >
      {{ text }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { isNumber } from "lodash-es";
import { computed } from "vue";

const props = defineProps<{
  val?: number;
  max?: number;
  size?: number | string;
  bgColor?: string;
  barColor?: string;
  barTransition?: string;
  barBorderRadius?: number;
  spacing?: number;
  text?: string;
  textAlign?: string;
  textPosition?: string;
  fontSize?: number;
  textFgColor?: string;
}>();

const val = computed(() => props.val ?? 0);
const max = computed(() => props.max ?? 100);
const size = computed(() => props.size ?? 3);
const bgColor = computed(() => props.bgColor ?? "#eee");
const barColor = computed(() => props.barColor ?? "#2196f3");
const barTransition = computed(() => props.barTransition ?? "all 0.5s ease");
const barBorderRadius = computed(() => props.barBorderRadius ?? 0);
const spacing = computed(() => props.spacing ?? 4);
const text = computed(() => props.text ?? "");
const textAlign = computed(() => props.textAlign ?? "center");
const textPosition = computed(() => props.textPosition ?? "bottom");
const fontSize = computed(() => props.fontSize ?? 13);
const textFgColor = computed(() => props.textFgColor ?? "#222");

const pct = computed(() => {
  let percent = (val.value / max.value) * 100;
  percent = parseFloat(percent.toFixed(2));
  return Math.min(percent, 100);
});

const sizePx = computed(() => {
  switch (size.value) {
    case "tiny":
      return 2;
    case "small":
      return 4;
    case "medium":
      return 8;
    case "large":
      return 12;
    case "big":
      return 16;
    case "huge":
      return 32;
    case "massive":
      return 64;
    default:
      return isNumber(size.value) ? Number(size.value) : 32;
  }
});

const textPadding = computed(() => {
  switch (size.value) {
    case "tiny":
    case "small":
    case "medium":
    case "large":
    case "big":
    case "huge":
    case "massive":
      return Math.min(Math.max(Math.ceil(sizePx.value / 8), 3), 12);
    default:
      return isNumber(spacing.value) ? spacing.value : 4;
  }
});

const textFontSize = computed(() => {
  switch (size.value) {
    case "tiny":
    case "small":
    case "medium":
    case "large":
    case "big":
    case "huge":
    case "massive":
      return Math.min(Math.max(Math.ceil(sizePx.value * 1.4), 11), 32);
    default:
      return isNumber(fontSize.value) ? fontSize.value : 13;
  }
});

const progressStyle = computed(() => {
  const style: Record<string, string | number> = {
    background: bgColor.value,
  };
  if (textPosition.value === "middle" || textPosition.value === "inside") {
    style.position = "relative";
    style["min-height"] = `${sizePx.value}px`;
    style["z-index"] = "-2";
  }
  if (barBorderRadius.value > 0) {
    style["border-radius"] = `${barBorderRadius.value}px`;
  }
  return style;
});

const barStyle = computed(() => {
  const style: Record<string, string | number> = {
    background: barColor.value,
    width: `${pct.value}%`,
    height: `${sizePx.value}px`,
    transition: barTransition.value,
  };
  if (barBorderRadius.value > 0) {
    style["border-radius"] = `${barBorderRadius.value}px`;
  }
  if (textPosition.value === "middle" || textPosition.value === "inside") {
    style.position = "absolute";
    style.top = "0";
    style.height = "100%";
    style["min-height"] = `${sizePx.value}px`;
    style["z-index"] = "-1";
  }
  return style;
});

const textStyle = computed(() => {
  const style: Record<string, string | number> = {
    color: textFgColor.value,
    "font-size": `${textFontSize.value}px`,
    "text-align": textAlign.value,
  };
  if (
    textPosition.value === "top" ||
    textPosition.value === "middle" ||
    textPosition.value === "inside"
  ) {
    style["padding-bottom"] = `${textPadding.value}px`;
  }
  if (
    textPosition.value === "bottom" ||
    textPosition.value === "middle" ||
    textPosition.value === "inside"
  ) {
    style["padding-top"] = `${textPadding.value}px`;
  }
  return style;
});
</script>
