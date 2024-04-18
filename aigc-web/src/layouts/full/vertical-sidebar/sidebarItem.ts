import { useUserStore } from "@/stores";

import {
  IconHome,
  IconVolume2,
  IconRegistered,
  IconSearch,
  IconReplace,
  IconList,
  IconColumns2,
  IconHexagons,
  IconCircleDot,
  IconAffiliate,
  IconMicroscope,
  IconBrandBinance,
  IconDeviceAirtag,
  IconSettings,
  IconColorPicker,
  IconTextSpellcheck,
  IconAdjustmentsDown,
  IconPhotoUp,
  IconPhotoX,
  IconPhotoCheck,
  IconTools
} from "@tabler/icons-vue";

export interface menu {
  header?: string;
  title?: string;
  icon?: any;
  iconSize?: "small";
  to?: string;
  chip?: string;
  chipBgColor?: string;
  chipColor?: string;
  chipVariant?: string;
  chipIcon?: string;
  children?: menu[];
  disabled?: boolean;
  type?: string;
  subCaption?: string;
}

let sidebarItem: menu[] = [
  {
    title: "homePage",
    icon: IconHome,
    to: "/dashboards/index"
  },
  { header: "aiAssistant" },
  {
    title: "aiAssistant",
    icon: IconTools,
    to: "/",
    children: [
      {
        title: "assistantList",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/ai-assistant/assistants/list"
      },
      {
        title: "toolsList",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/ai-assistant/tools-list"
      }
    ]
  },

  { header: "sampleLibrary" },
  // {  delete
  //   title: "soundSample",
  //   icon: IconColorPicker,
  //   to: "/",
  //   children: [
  //     {
  //       title: "audioAnnotation",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/audio-manage/audio-mark"
  //     }
  //   ]
  // },
  {
    title: "textSample",
    icon: IconTextSpellcheck,
    to: "/",
    children: [
      // {
      //   title: "fine-tuneSampleList",
      //   icon: IconCircleDot,
      //   iconSize: "small",
      //   to: "/sample-library/text-sample/list"
      // },
      {
        title: "dataSampleList",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/sample-library/mgr-datasets/list"
      },
      {
        title: "textAnnotationList",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/sample-library/text-mark/list"
      }
      // {
      //   title: "意图模型标注列表",
      //   icon: IconCircleDot,
      //   iconSize: "small",
      //   to: "/sample-library/intention-mark/list"
      // }
    ]
  },
  // {
  //   title: "voiceprintRecognition",
  //   icon: IconVolume2,
  //   to: "/voice-print/analysis"
  // },
  // {
  //   title: "voiceprintRegistration",
  //   icon: IconRegistered,
  //   to: "/voice-print/register"
  // },
  // {
  //   title: "voiceprintQuery",
  //   icon: IconSearch,
  //   to: "/voice-print/search"
  // },
  // {
  //   title: "voiceprintLibrary",
  //   icon: IconList,
  //   to: "/voice-print/library-list"
  // },
  // // {
  // //   title: "audioAnnotation",
  // //   icon: IconCircleDot,
  // //   iconSize: "small",
  // //   to: "/audio-manage/audio-mark"
  // // },

  // { header: "voiceServices" },
  // {
  //   title: "voiceToText",
  //   icon: IconReplace,
  //   to: "/voice-print/translation"
  // },
  // {
  //   title: "voiceSynthesis",
  //   icon: IconHexagons,
  //   to: "/",
  //   children: [
  //     {
  //       title: "speakerManagement",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/voice-print/synthesis/speaker"
  //     },
  //     {
  //       title: "voiceList",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/voice-print/synthesis/voice-list"
  //     },
  //     {
  //       title: "syntheticSound",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/voice-print/synthesis/synthesis-voice"
  //     }
  //   ]
  // },
  // {
  //   title: "voiceDenoiseList",
  //   icon: IconAdjustmentsDown,
  //   to: "/voice-print/denoise-list"
  // },

  // { header: "voiceServices" },
  // {
  //   title: "voiceToText",
  //   icon: IconReplace,
  //   to: "/voice-print/translation"
  // },
  // {
  //   title: "ttsSynthesis",
  //   icon: IconHexagons,
  //   to: "/",
  //   children: [
  //     {
  //       title: "speakerManagement",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/voice-print/synthesis/speaker"
  //     },
  //     {
  //       title: "ttsSynthesisList",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/voice-print/synthesis/voice-list"
  //     }
  //   ]
  // },
  // {
  //   title: "intelligentVoiceprint",
  //   icon: IconColumns2,
  //   to: "/",
  //   children: [
  //     {
  //       title: "voiceprintLibrary",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/voice-print/library-list"
  //     },
  //     {
  //       title: "voiceprintCompareList",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/voice-print/compare-list"
  //     }
  //   ]
  // },
  // {
  //   title: "voiceDenoiseList",
  //   icon: IconAdjustmentsDown,
  //   to: "/voice-print/denoise-list"
  // },

  // { header: "imageServices" },
  // {
  //   title: "imageSuperResolution",
  //   icon: IconPhotoUp,
  //   to: "/image-services/super-resolution/list"
  // },
  // {
  //   title: "imageMatting",
  //   icon: IconPhotoCheck,
  //   to: "/image-services/image-matting/list"
  // },
  // {
  //   title: "faceRecognition",
  //   icon: IconPhotoX,
  //   to: "/image-services/face-recognition/list"
  // },

  // { header: "videoService" },
  // {
  //   title: "digitalHuman",
  //   icon: IconAffiliate,
  //   to: "/",
  //   children: [
  //     {
  //       title: "characterList",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/digital-human/character-list"
  //     },
  //     {
  //       title: "videoList",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/digital-human/video-list/list"
  //     }
  //   ]
  // },
  { header: "textModel" },
  {
    title: "modelManagement",
    icon: IconHexagons,
    to: "/",
    children: [
      {
        title: "modelList",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/model/model-list"
      },
      {
        title: "modelChatPlayground",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/model/chat-playground"
      }
    ]
  },
  {
    title: "modelFine-tuning",
    icon: IconMicroscope,
    to: "/",
    children: [
      {
        title: "fine-tuneTasks",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/model/fine-tuning"
      }
    ]
  },
  // {
  //   title: "sceneManagement",
  //   icon: IconBrandBinance,
  //   to: "/",
  //   children: [
  //     {
  //       title: "sceneList",
  //       icon: IconCircleDot,
  //       iconSize: "small",
  //       to: "/scene/scene-list"
  //     }
  //   ]
  // },
  { header: "systemManagement" },
  {
    title: "systemManagement",
    icon: IconSettings,
    to: "/",
    children: [
      {
        title: "systemDictionary",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/system/dict"
      },
      {
        title: "templateManagement",
        icon: IconCircleDot,
        iconSize: "small",
        to: "/system/template"
      }
    ]
  }
];

const userStore = useUserStore();
if (userStore.userInfo.tenantId == "5f9b3b3d-9b9c-4e1a-8e1a-5a4b4b4b4b43") {
  sidebarItem = [
    {
      title: "audioAnnotation",
      icon: IconCircleDot,
      iconSize: "small",
      to: "/audio-manage/audio-mark"
    }
  ];
}

export default sidebarItem;
