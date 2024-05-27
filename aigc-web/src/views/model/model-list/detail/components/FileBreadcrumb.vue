<template>
  <v-breadcrumbs :items="rootPath.concat(items)">
    <template #title="{ item }">
      <a class="hover:underline cursor-pointer" @click="navigateTo(item)">{{ item.title }}</a>
    </template>
  </v-breadcrumbs>
</template>
<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const { path } = route;

const jobId = route.query.jobId;
const modelName = route.query.modelName as string;

const rootPath = [
  {
    title: modelName
  }
];

const items = computed(() => {
  const filePath = route.query.filePath as string;
  if (!filePath) return [];
  const result = [];
  const pathArr = filePath.split("/");
  let levelFilePath = "";
  for (let i = 1; i < pathArr.length; i++) {
    if (!pathArr[i]) continue;
    levelFilePath += `/${pathArr[i]}`;
    const item = {
      title: pathArr[i],
      filePath: levelFilePath
    };
    result.push(item);
  }
  return result;
});

const navigateTo = item => {
  router.push({
    path,
    query: { ...route.query, filePath: item.filePath }
  });
};
</script>
