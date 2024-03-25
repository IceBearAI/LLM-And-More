<script setup>
import NavItem from "../NavItem/index.vue";

const props = defineProps({ item: Object, level: Number });
</script>

<template>
  <!-- ---------------------------------------------- -->
  <!---Item Childern -->
  <!-- ---------------------------------------------- -->
  <v-list-group no-action>
    <!-- ---------------------------------------------- -->
    <!---Dropdown  -->
    <!-- ---------------------------------------------- -->
    <template v-slot:activator="{ props }">
      <v-list-item v-bind="props" :value="item.title" rounded class="mb-1">
        <!---Icon  -->
        <template v-slot:prepend>
          <component :is="item.icon" :level="level" stroke-width="1" :size="item.iconSize == 'small' ? 12 : 20" />
        </template>
        <!---Title  -->
        <v-list-item-title class="mr-auto">{{ $t(item.title) }}</v-list-item-title>
        <!---If Caption-->
        <v-list-item-subtitle v-if="item.subCaption" class="text-caption mt-n1 hide-menu">
          {{ item.subCaption }}
        </v-list-item-subtitle>
      </v-list-item>
    </template>
    <!-- ---------------------------------------------- -->
    <!---Sub Item-->
    <!-- ---------------------------------------------- -->
    <template v-for="(subitem, i) in item.children" :key="i" v-if="item.children">
      <NavCollapse :item="subitem" v-if="subitem.children" :level="level + 1" />
      <NavItem :item="subitem" :level="level + 1" v-else></NavItem>
    </template>
  </v-list-group>

  <!-- ---------------------------------------------- -->
  <!---End Item Sub Header -->
  <!-- ---------------------------------------------- -->
</template>
