<script setup lang="ts">
import { RouterLink } from "vue-router";
import logo from "@/assets/images/logos/logo-dark.png";
</script>
<template>
  <div class="logo">
    <RouterLink to="/">
      <img style="width: 106px; height: 38px" :src="logo" alt="home" />
    </RouterLink>
  </div>
</template>
