<template>
  <v-row class="mb-5" align-content="center" justify="center" :style="{ height }" no-gutters>
    <template v-if="loading">
      <v-col class="text-subtitle-1 text-center mb-4" cols="12"> 生成中 </v-col>
      <v-col cols="6">
        <v-progress-linear color="primary" indeterminate rounded height="6"></v-progress-linear>
      </v-col>
    </template>
    <template v-else>
      <v-col cols="12">
        <v-img v-if="url" :src="url" :height="height" alt="previewImageUrl" cover class="rounded-md align-end text-right">
          <template #placeholder>
            <div class="d-flex align-center justify-center fill-height">
              <v-progress-circular color="primary" indeterminate></v-progress-circular>
            </div>
          </template>
        </v-img>
      </v-col>
    </template>
  </v-row>
</template>
<script setup lang="ts">
interface IProps {
  loading?: boolean;
  url: string;
  height?: string;
}

withDefaults(defineProps<IProps>(), {
  loading: false,
  url: "",
  height: "350px"
});
</script>
