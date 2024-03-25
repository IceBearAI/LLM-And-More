<script setup>
import Icon from "../../vertical-sidebar/Icon.vue";
const props = defineProps({ item: Object, level: Number });
</script>

<template>
  <!---Single Item-->
  <router-link :to="`${item.to}`" class="navItemLink rounded-md" :disabled="item.disabled">
    <!---If icon-->
    <i class="navIcon"> <Icon :item="item.icon" :level="level" /></i>
    <span>{{ $t(item.title) }}</span>
    <!---If Caption-->
    <small v-if="item.subCaption" class="text-caption mt-n1 hide-menu">
      {{ item.subCaption }}
    </small>
    <!---If any chip or label-->
    <template v-if="item.chip">
      <v-chip
        :color="item.chipColor"
        class="sidebarchip hide-menu ml-auto"
        :size="item.chipIcon ? 'small' : 'small'"
        :variant="item.chipVariant"
        :prepend-icon="item.chipIcon"
      >
        {{ item.chip }}
      </v-chip>
    </template>
  </router-link>
</template>
