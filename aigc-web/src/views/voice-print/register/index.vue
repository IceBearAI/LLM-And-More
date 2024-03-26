<template>
  <NavBack v-if="route.path === '/voice-print/library-list/register'" backUrl="/voice-print/library-list/list">
    {{ page.title }}
  </NavBack>
  <BaseBreadcrumb v-else :title="page.title" :breadcrumbs="breadcrumbs"></BaseBreadcrumb>
  <UiParentCard title="注册">
    <v-row justify="center">
      <v-col cols="12" md="6">
        <v-form ref="formRef">
          <v-row>
            <v-col cols="12">
              <div class="d-flex flex-column flex-sm-row gap-3 justify-center">
                <AiAudio :src="audioUrl" />
              </div>
              <v-label class="mb-2 font-weight-medium">音频文件</v-label>
              <v-file-input
                v-model="formData.file"
                accept="audio/*, .mp3, .wma, .amr, .wav, .m4a"
                prepend-icon="mdi-volume-high"
                :rules="formRules.file"
                hide-details="auto"
                variant="outlined"
                required
                multiple
              ></v-file-input>
            </v-col>
            <v-col cols="12">
              <v-label class="mb-2 font-weight-medium">用户姓名</v-label>
              <v-text-field v-model="formData.userName" hide-details="auto" :rules="formRules.userName" required></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-label class="mb-2 font-weight-medium">用户标识</v-label>
              <v-text-field v-model="formData.userKey" hide-details="auto" :rules="formRules.userKey" required></v-text-field>
            </v-col>
          </v-row>

          <div class="d-flex flex-column">
            <v-btn :loading="submitLoading" color="primary" class="mt-4" block @click="validate"> 提交 </v-btn>
          </div>
        </v-form>
      </v-col>
    </v-row>
  </UiParentCard>
</template>
<script setup lang="ts">
import { ref, reactive, computed, inject } from "vue";
import { useRoute, useRouter } from "vue-router";
import NavBack from "@/components/business/NavBack.vue";
import BaseBreadcrumb from "@/components/shared/BaseBreadcrumb.vue";
import UiParentCard from "@/components/shared/UiParentCard.vue";
import AiAudio from "@/components/business/AiAudio.vue";
import { http } from "@/utils";
import { type ItfAspectPageState } from "@/types/AspectPageTypes.ts";

const provideAspectPage = inject("provideAspectPage") as ItfAspectPageState;
const route = useRoute();
const router = useRouter();
const page = ref({ title: "声纹注册" });
const breadcrumbs = ref([
  {
    text: "智能声纹",
    disabled: false,
    href: "#"
  },
  {
    text: "声纹注册",
    disabled: true,
    href: "#"
  }
]);

const formRef = ref();
const submitLoading = ref(false);
const formData = reactive({
  file: [],
  userName: "",
  userKey: ""
});
const formRules = reactive({
  file: [(v: string) => (v && v.length !== 0) || "请上传音频文件"],
  userName: [v => !!v || "请输入用户姓名"],
  userKey: [v => !!v || "请输入用户标识"]
});

const audioUrl = computed(() => {
  if (formData.file && formData.file.length > 0) {
    return URL.createObjectURL(formData.file[0]);
  }
  return "";
});

const validate = async () => {
  const { valid } = await formRef.value.validate();
  if (valid) {
    submitLoading.value = true;
    const [err, res] = await http.upload({
      url: "/voice/register",
      showSuccess: "注册成功",
      data: formData
    });
    if (res) {
      if (route.path !== "/voice-print/library-list/register") {
        formRef.value.reset();
      } else {
        // 如果是从声纹库跳转过来的, 注册成功后跳转回声纹库页面
        provideAspectPage.methods.refreshListPage();
        router.replace("/voice-print/library-list/list");
      }
    }
    submitLoading.value = false;
  }
};
</script>
