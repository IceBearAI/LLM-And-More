import{r as x,i as A,z as b,x as w,U,k as h,l as $,m as a,j as e,_ as m,n as t,c as M,W as N,Z as B,$ as q,Y as j,X as z,ac as E,au as G,D as H,o as O,S as Q,Q as T,a3 as W,L as X,F as Y,av as Z,a8 as P,a1 as J}from"./utils-7a539270.js";import{_ as D}from"./UiParentCard.vue_vue_type_script_setup_true_lang-f6f68ebc.js";import{_ as K}from"./UiChildCard.vue_vue_type_script_setup_true_lang-f8c4b585.js";import{_ as ee}from"./NavBack.vue_vue_type_script_setup_true_lang-a80a61f1.js";import{S as ae}from"./SpeakerSelector-0f578bdc.js";import{V as C}from"./VRow-220780f3.js";import{V as u}from"./VCol-b16873ba.js";import{a4 as y,ae as te,a5 as R,_ as le,O as se,ag as oe}from"./index-de28ea06.js";import{V as I}from"./VForm-3af80d80.js";import{_ as re}from"./Explain.vue_vue_type_style_index_0_lang-e3bfca38.js";import{V as ie}from"./VTextarea-b952e231.js";import{V as ne}from"./VSlider-1987474d.js";import{_ as de}from"./AiAudio.vue_vue_type_style_index_0_lang-0cb907a9.js";import"./IconCircleCheckFilled-215ef6ba.js";import"./IconInfoCircle-86362c10.js";const ue={__name:"Left",setup(F,{expose:k}){const _=x(),v=x(),i=A("provideSynthesisVoice"),{formData:n}=b(i),d=w({selectedSpeaker:{}}),r=w({provider:[o=>!!o||"请选择供应"],lang:[o=>!!o||"请选择语言"],speakName:[o=>!!o||"请选择发声人"]}),f=()=>{v.value.reload({lang:n.value.lang,provider:n.value.provider})};return k({validate(){return _.value.validate()}}),(o,l)=>{const p=U("Select");return h(),$(I,{ref_key:"refForm",ref:_,class:"my-form"},{default:a(()=>[e(C,null,{default:a(()=>[e(u,{cols:"12"},{default:a(()=>[e(y,{class:"mb-2"},{default:a(()=>[m("供应")]),_:1}),e(p,{mapDictionary:{code:"speak_provider"},placeholder:"请选择供应",modelValue:t(n).provider,"onUpdate:modelValue":l[0]||(l[0]=s=>t(n).provider=s),onChange:f},null,8,["modelValue"])]),_:1}),e(u,{cols:"12"},{default:a(()=>[e(y,{class:"mb-2"},{default:a(()=>[m("语言")]),_:1}),e(p,{mapDictionary:{code:"speak_lang"},placeholder:"请选择语言",modelValue:t(n).lang,"onUpdate:modelValue":l[1]||(l[1]=s=>t(n).lang=s),infos:t(i).selectedLanguage,"onUpdate:infos":l[2]||(l[2]=s=>t(i).selectedLanguage=s),onChange:f},null,8,["modelValue","infos"])]),_:1}),e(u,{cols:"12"},{default:a(()=>[e(y,{class:"mb-2 required"},{default:a(()=>[m("请选择需要合成的发声人")]),_:1}),e(te,{rules:r.speakName,modelValue:t(n).speakName,"onUpdate:modelValue":l[5]||(l[5]=s=>t(n).speakName=s)},{default:a(()=>[e(ae,{ref_key:"refSpeakerSelector",ref:v,modelValue:t(n).speakName,"onUpdate:modelValue":l[3]||(l[3]=s=>t(n).speakName=s),infos:d.selectedSpeaker,"onUpdate:infos":l[4]||(l[4]=s=>d.selectedSpeaker=s)},null,8,["modelValue","infos"])]),_:1},8,["rules","modelValue"])]),_:1})]),_:1})]),_:1},512)}}},me={key:0,class:"text-primary font-weight-bold"},pe={style:{width:"300px"}},ce={__name:"Right",setup(F,{expose:k}){const _=x(),v=A("provideSynthesisVoice"),{formData:i}=b(v),n=w({speedConfig:{min:.5,max:2,step:.1}}),{speedConfig:d}=b(n),r=w({title:[o=>!!o||"请输入标题"],text:[o=>!!o||"请输入语音播放文本"],speed:[o=>{if(o){let{min:l,max:p}=n.speedConfig;return o<l?"语速不能低于"+l+"倍":o>p?"语速不能高于"+p+"倍":!0}else return"请设置语速"}]}),f=M(()=>{let{label:o}=v.selectedLanguage;return o?o.trim().replace(/([^(（]+)(.*)$/,"$1"):""});return k({validate(){return _.value.validate()}}),(o,l)=>{const p=U("Select");return h(),$(I,{ref_key:"refForm",ref:_,class:"my-form"},{default:a(()=>[e(C,null,{default:a(()=>[e(u,{cols:"12"},{default:a(()=>[e(y,{class:"mb-2 required"},{default:a(()=>[m("标题"),e(re,{class:"ml-2"},{default:a(()=>[m("用于列表展示和搜索，能够快速了解基本信息")]),_:1})]),_:1}),e(R,{density:"compact",variant:"outlined",placeholder:"请输入标题","hide-details":"auto",rules:r.title,modelValue:t(i).title,"onUpdate:modelValue":l[0]||(l[0]=s=>t(i).title=s),modelModifiers:{trim:!0},clearable:""},null,8,["rules","modelValue"])]),_:1}),e(u,{cols:"12"},{default:a(()=>[e(y,{class:"mb-2 required",style:{"white-space":"inherit"}},{default:a(()=>[N("div",null,[m(" 请输入"),f.value?(h(),B("span",me,"「 "+q(f.value)+" 」",1)):j("",!0),m("语音播放文本，文本内容小于200个字(包括标点符号)。 ")])]),_:1}),e(ie,{modelValue:t(i).text,"onUpdate:modelValue":l[1]||(l[1]=s=>t(i).text=s),modelModifiers:{trim:!0},rules:r.text,placeholder:"语音播放文本",counter:"",rows:"5",maxlength:"200"},null,8,["modelValue","rules"])]),_:1}),e(u,{cols:"12"},{default:a(()=>[e(y,{class:"mb-2"},{default:a(()=>[m("语气")]),_:1}),e(p,{mapDictionary:{code:"speak_tone"},placeholder:"请选择语气",modelValue:t(i).tone,"onUpdate:modelValue":l[2]||(l[2]=s=>t(i).tone=s)},null,8,["modelValue"])]),_:1}),e(u,{cols:"12"},{default:a(()=>[e(y,{class:"mb-2 required"},{default:a(()=>[m("语速")]),_:1}),N("div",pe,[e(ne,{density:"compact",modelValue:t(i).speed,"onUpdate:modelValue":l[4]||(l[4]=s=>t(i).speed=s),color:"primary",max:t(d).max,min:t(d).min,step:t(d).step,"thumb-label":"",rules:r.speed,"hide-details":"auto"},{append:a(()=>[e(R,{modelValue:t(i).speed,"onUpdate:modelValue":l[3]||(l[3]=s=>t(i).speed=s),modelModifiers:{number:!0},type:"number",density:"compact",max:t(d).max,min:t(d).min,step:t(d).step,style:{width:"80px"}},null,8,["modelValue","max","min","step"])]),_:1},8,["modelValue","max","min","step","rules"])])]),_:1}),z(o.$slots,"default")]),_:3})]),_:3},512)}}};const fe={class:"hv-center py-3"},_e={__name:"synthesisVoice",setup(F){const{mappings:k,loadDictTree:_}=E(),v=G(),i=x(),n=x(),d=x(),r=w({style:{showPreview:!1,loadingPreview:!1},formData:{provider:"",lang:"",speakName:"",text:"",title:"",speed:1,tone:null},selectedLanguage:{label:"",value:""},selectedSpeaker:{speakName:"",speakCname:"",headImg:"",speakDemo:"",gender:"",ageGroup:"",subTitle:""},selectedDigitalHuman:{name:"",cname:"",cover:"",video:""},audioInfo:{src:"",gender:"",type:"complex"}}),{style:f,formData:o,selectedSpeaker:l}=b(r);H("provideSynthesisVoice",r);const p=c=>{r.style.loadingPreview=!0,r.style.showPreview=!0,setTimeout(()=>{oe.scrollTo(d.value.$el,500),P.assign(r.audioInfo,{src:c.s3Url,gender:c.gender})},100),setTimeout(()=>{r.style.loadingPreview=!1},500)},s=async()=>{const{valid:c}=await i.value.validate(),{valid:S}=await n.value.validate();if(c&&S){const[g,V]=await Z.post({showLoading:"btn#btnSubmit",showSuccess:!0,url:"/api/voice/tts",data:P.omit(r.formData,["provider","lang"])});V?p(V):r.style.showPreview=!1}else{let g="请处理页面标错的地方后，再尝试提交";J.warning(g)}};return O(async()=>{var V,L;await _(["speak_provider","speak_lang"]);let{provider:c,lang:S,speakName:g}=v.query;c&&((V=k.speak_provider)!=null&&V[c])&&(r.formData.provider=c),S&&((L=k.speak_lang)!=null&&L[S])&&(r.formData.lang=S),g&&(r.formData.speakName=g)}),(c,S)=>{const g=U("AiBtn"),V=Q("loading");return h(),B(Y,null,[e(ee,{backUrl:"/voice-print/synthesis/voice-list"},{default:a(()=>[m("创建TTS")]),_:1}),e(C,{class:"mt-1"},{default:a(()=>[e(u,{cols:"12",md:"6"},{default:a(()=>[e(D,null,{default:a(()=>[e(ue,{ref_key:"refFormLeft",ref:i},null,512)]),_:1})]),_:1}),e(u,{cols:"12",md:"6"},{default:a(()=>[e(D,null,{default:a(()=>[e(ce,{ref_key:"refFormRight",ref:n},{default:a(()=>[e(u,{cols:"12",class:"mt-4"},{default:a(()=>[N("div",fe,[e(g,{id:"btnSubmit",color:"secondary",width:"200",height:"48",size:"large",onClick:s},{default:a(()=>[m("开始合成")]),_:1})])]),_:1})]),_:1},512)]),_:1})]),_:1}),e(se,null,{default:a(()=>[T(e(u,{cols:"12"},{default:a(()=>[e(D,null,{default:a(()=>[e(u,{cols:"12",class:"pt-10"},{default:a(()=>[T((h(),$(K,{title:"合成音频"},{default:a(()=>[e(de,X({...r.audioInfo},{ref_key:"refAiAudio",ref:d}),null,16)]),_:1})),[[V,t(f).loadingPreview]])]),_:1})]),_:1})]),_:1},512),[[W,t(f).showPreview]])]),_:1})]),_:1})],64)}}},Fe=le(_e,[["__scopeId","data-v-665157a5"]]);export{Fe as default};
