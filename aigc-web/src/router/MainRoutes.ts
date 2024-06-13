import type { RouteRecordRaw } from "vue-router";
const MainRoutes: RouteRecordRaw = {
  path: "/main",
  meta: {
    requiresAuth: true
  },
  redirect: "/main",
  component: () => import("@/layouts/full/FullLayout.vue"),
  children: [
    {
      path: "/",
      redirect: "/dashboards/index"
    },
    {
      name: "DashboardsIndex",
      path: "/dashboards/index",
      component: () => import("@/views/dashboards/index.vue")
    },
    {
      name: "VPAnalysis",
      path: "/voice-print/analysis",
      component: () => import("@/views/voice-print/analysis/index.vue")
    },
    {
      name: "VPRegister",
      path: "/voice-print/register",
      component: () => import("@/views/voice-print/register/index.vue")
    },
    {
      name: "VPSearch",
      path: "/voice-print/search",
      component: () => import("@/views/voice-print/search/index.vue")
    },
    {
      name: "VPTranslation",
      path: "/voice-print/translation",
      component: () => import("@/views/voice-print/translation/index.vue")
    },
    {
      name: "VPLibraryList",
      path: "/voice-print/library-list",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: { aspectPageInclude: ["VoicePrintLibraryList"] },
      redirect: "/voice-print/library-list/list",
      children: [
        {
          path: "/voice-print/library-list/list",
          component: () => import("@/views/voice-print/library-list/index.vue")
        },
        {
          path: "/voice-print/library-list/register",
          component: () => import("@/views/voice-print/register/index.vue")
        },
        {
          path: "/voice-print/library-list/search",
          component: () => import("@/views/voice-print/search/index.vue")
        }
      ]
    },
    {
      name: "VPCompareList",
      path: "/voice-print/compare-list",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: { aspectPageInclude: ["VoicePrintCompareList"] },
      redirect: "/voice-print/compare-list/list",
      children: [
        {
          path: "/voice-print/compare-list/list",
          component: () => import("@/views/voice-print/compare-list/index.vue")
        },
        {
          path: "/voice-print/compare-list/compare",
          component: () => import("@/views/voice-print/analysis/index.vue")
        }
      ]
    },
    {
      name: "VPDenoiseList",
      path: "/voice-print/denoise-list",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: { aspectPageInclude: ["denoiseList"] },
      redirect: "/voice-print/denoise-list/list",
      children: [
        {
          path: "/voice-print/denoise-list/list",
          component: () => import("@/views/voice-print/denoise-list/denoiseList.vue")
        },
        {
          path: "/voice-print/denoise-list/denoise",
          component: () => import("@/views/voice-print/denoise-list/denoise.vue")
        }
      ]
    },
    {
      path: "/image-services/super-resolution/list",
      component: () => import("@/views/image-services/super-resolution/superResolution.vue")
    },
    {
      path: "/image-services/image-matting/list",
      component: () => import("@/views/image-services/image-matting/imageMatting.vue")
    },
    // 人脸服务
    {
      path: "/face-services/face-recognition/list",
      component: () => import("@/views/image-services/face-recognition/faceRecognition.vue")
    },
    {
      path: "/face-services/live-recognition/list",
      component: () => import("@/views/face-services/live-recognition/liveRecognition.vue")
    },
    {
      path: "/face-services/face-reg/list",
      component: () => import("@/views/face-services/face-reg/faceReg.vue")
    },
    {
      path: "/face-services/face-search-log/list",
      component: () => import("@/views/face-services/face-search-log/faceSearchLog.vue")
    },
    // ocr服务
    {
      path: "/ocr-services/ocr-detection/list",
      component: () => import("@/views/ocr-services/ocr-detection/ocrDetection.vue")
    },
    {
      path: "/voice-print/synthesis/speaker",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["voiceList"]
      },
      redirect: "/voice-print/synthesis/speaker/list",
      children: [
        {
          path: "list",
          component: () => import("@/views/voice-print/synthesis/speaker-manage/speakerManage.vue")
        },
        {
          path: "clone-detail",
          component: () => import("@/views/voice-print/synthesis/speaker-manage/detail/cloneDetail.vue")
        }
      ]
    },
    {
      path: "/voice-print/synthesis",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["voiceList"]
      },
      redirect: "/voice-print/synthesis/voice-list",
      children: [
        {
          path: "voice-list",
          component: () => import("@/views/voice-print/synthesis/voice-list/voiceList.vue")
        },
        {
          path: "synthesis-voice",
          component: () => import("@/views/voice-print/synthesis/synthesis-voice/synthesisVoice.vue")
        }
      ]
    },

    // {
    //   path: "/model/model-list",
    //   component: () => import("@/views/model/model-list/modelList.vue")
    // },

    {
      path: "/model/model-list",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["modelList"]
      },
      redirect: "/model/model-list/list",
      children: [
        {
          path: "/model/model-list/list",
          component: () => import("@/views/model/model-list/modelList.vue")
        },
        {
          path: "/model/model-list/detail",
          component: () => import("@/views/model/model-list/detail/modelListDetail.vue")
        }
      ]
    },

    {
      name: "ModelChatPlayground",
      path: "/model/chat-playground",
      component: () => import("@/views/model/chat-playground/chatPlayground.vue")
    },

    {
      name: "ModelTerminal",
      path: "/model/terminal",
      component: () => import("@/views/model/terminal/terminal.vue")
    },

    {
      path: "/model/fine-tuning",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["fineTuningList"]
      },
      redirect: "/model/fine-tuning/list",
      children: [
        {
          path: "/model/fine-tuning/list",
          component: () => import("@/views/model/fine-tuning/list/fineTuningList.vue")
        },
        {
          path: "/model/fine-tuning/detail",
          component: () => import("@/views/model/fine-tuning/detail/fineTuningDetail.vue")
        }
      ]
    },
    {
      path: "/scene/scene-list",
      component: () => import("@/views/scene/scene-list/sceneList.vue")
    },
    {
      path: "/digital-human/character-list",
      component: () => import("@/views/digital-human/characterList/characterList.vue")
    },
    {
      path: "/digital-human/video-list",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["DigitalVideoList"]
      },
      redirect: "/digital-human/video-list/list",
      children: [
        {
          path: "list",
          component: () => import("@/views/digital-human/videoList/videoList.vue")
        },
        {
          path: "detail",
          component: () => import("@/views/digital-human/videoList/components/Detail.vue")
        },
        {
          path: "edit",
          component: () => import("@/views/digital-human/edit/digitalHumanEdit.vue")
        }
      ]
    },
    {
      path: "/sample-library/text-sample",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["textSampleList"]
      },
      redirect: "/sample-library/text-sample/list",
      children: [
        {
          path: "list",
          component: () => import("@/views/sample-library/text-sample/textSampleList.vue")
        },
        {
          path: "detail",
          component: () => import("@/views/sample-library/text-sample/detail/textSampleDetail.vue")
        }
      ]
    },
    {
      path: "/sample-library/mgr-datasets/list",
      component: () => import("@/views/sample-library/mgr-datasets/mgrDatasets.vue")
    },
    {
      path: "/sample-library/text-mark",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["textMarkList"]
      },
      redirect: "/sample-library/text-mark/list",
      children: [
        {
          path: "list",
          component: () => import("@/views/sample-library/text-mark/textMarkList.vue")
        },
        {
          path: "detail",
          component: () => import("@/views/sample-library/text-mark/detail/textMarkDetail.vue")
        }
      ]
    },

    {
      path: "/sample-library/intention-mark",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["intentionMarkList"]
      },
      redirect: "/sample-library/intention-mark/list",
      children: [
        {
          path: "list",
          component: () => import("@/views/sample-library/intention-mark/intentionMarkList.vue")
        },
        {
          path: "detail",
          component: () => import("@/views/sample-library/intention-mark/detail/intentionMarkDetail.vue")
        }
      ]
    },

    {
      path: "/ai-assistant/tools-list",
      component: () => import("@/views/ai-assistant/tools-list/toolsList.vue")
    },
    {
      path: "/ai-assistant/assistants",
      component: () => import("@/components/business/AspectPage.vue"),
      meta: {
        aspectPageInclude: ["assistantsList"]
      },
      redirect: "/ai-assistant/assistants/list",
      children: [
        {
          path: "list",
          component: () => import("@/views/ai-assistant/assistants-list/assistantsList.vue")
        },
        {
          path: "detail",
          component: () => import("@/views/ai-assistant/assistants-list/detail/assistantsDetail.vue")
        },
        {
          path: "playground",
          component: () => import("@/views/ai-assistant/assistants-list/playground/playground.vue")
        }
      ]
    },
    {
      path: "/audio-manage/audio-mark",
      component: () => import("@/views/audio-manage/audio-mark/audioMark.vue")
    },
    {
      path: "/system/template",
      component: () => import("@/views/system/template/template.vue")
    },
    {
      path: "/system/dict",
      component: () => import("@/views/system/dict/dict.vue")
    },
    {
      path: "/system/tenant/list",
      component: () => import("@/views/system/tenants/tenantList.vue")
    },
    {
      path: "/system/user/list",
      component: () => import("@/views/system/user-list/userList.vue")
    }
  ]
};

export default MainRoutes;
