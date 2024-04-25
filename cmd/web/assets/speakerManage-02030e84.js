import{d as L,ac as Q,x as N,r as b,c as W,U as V,k as c,l as C,m as a,W as d,j as e,_ as u,Y as J,Z as U,a5 as O,F as q,av as G,al as Z,am as X,o as ee,S as ae,aA as le,Q as te,$ as S,n as P,aB as oe}from"./utils-fc8ebe1f.js";import{_ as se}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-02f17eb8.js";import{_ as re}from"./UiParentCard.vue_vue_type_script_setup_true_lang-bc14ee23.js";import{A as ne}from"./AlertBlock-c7626136.js";import{_ as de}from"./ConfirmByInput.vue_vue_type_style_index_0_lang-10388eaa.js";import{_ as ue}from"./AiAudio.vue_vue_type_style_index_0_lang-0c48023d.js";import{_ as Y}from"./Explain.vue_vue_type_style_index_0_lang-2aca43e9.js";import{_ as ie}from"./UploadFile.vue_vue_type_script_setup_true_lang-01a60301.js";import{a5 as M,ae as pe,af as me,J as ce,M as _e,_ as fe,V as ve}from"./index-b84f97ad.js";import{V as ke,a as be}from"./VRadioGroup-d5988ff5.js";import{V as ge}from"./VAlert-2b69a982.js";import{V as he}from"./VSwitch-a8718d97.js";import{V as Ve}from"./VTextarea-e2668a01.js";import{V as ye}from"./VForm-cdbc2f03.js";import{C as xe}from"./ChipStatus-84a878a5.js";import{V as Ce}from"./VRow-4fd670bb.js";import{V as I}from"./VCol-f5dbb3fc.js";import"./IconInfoCircle-d065b385.js";import"./Confirm-dfed4e97.js";import"./VFileInput-6c8fbbc5.js";const f=D=>(Z("data-v-9865707a"),D=D(),X(),D),Se={class:"mx-auto mt-3",style:{width:"540px"}},Ie={class:"required"},De={class:"required"},we=f(()=>d("a",{class:"link",href:"http://docs.paas.creditease.corp/paas-gpt/voice/voice_personal.html",target:"_blank"},"声音克隆（个人版）",-1)),Ue=f(()=>d("label",{class:"required"},"姓名",-1)),Ne=f(()=>d("label",{class:"required"},"性别",-1)),Te=f(()=>d("label",{class:"required"},"年龄段",-1)),Fe=f(()=>d("label",{class:"required"},"语言",-1)),Pe=f(()=>d("label",{class:"required"},"风格",-1)),qe=f(()=>d("label",{class:"required"},"适应范围",-1)),Be=f(()=>d("label",null,"头像",-1)),$e=f(()=>d("label",{class:"required"},"启用",-1)),Re=f(()=>d("label",null,"备注",-1)),Ae=L({__name:"CreateSpeakerPane",emits:["submit"],setup(D,{expose:B,emit:w}){const T={provider:null,speakName:"",speakCname:"",gender:1,ageGroup:null,lang:null,speakStyle:null,area:null,headImgFileId:"",enabled:!1,remark:""},$=w,{options:v}=Q(),_=N({operateType:"add"}),o=b({...T}),p=b(null),y=b(),g=b(),k=N({provider:[n=>!!n||"请选择供应商"],speakName:[n=>!!n||"请输入标识"],speakCname:[n=>!!n||"请输入姓名"],ageGroup:[n=>!!n||"请选择年龄段"],lang:[n=>!!n||"请选择语言"],speakStyle:[n=>!!n||"请选择风格"],area:[n=>!!n||"请选择适应范围"]}),h=W(()=>_.operateType==="edit"),x=W(()=>v.speak_gender),R=()=>{o.value.headImgFileId="",p.value=null},A=async({valid:n,showLoading:s})=>{if(n){const r={url:"",method:""};_.operateType=="add"?(r.url="/voice/speak",r.method="post"):(r.url=`/voice/speak/${o.value.id}`,r.method="put");const[i,l]=await G[r.method]({showLoading:s,showSuccess:!0,url:r.url,data:o.value});l&&(y.value.hide(),$("submit"))}};return B({show({title:n,operateType:s,infos:r}){y.value.show({title:n,refForm:g}),_.operateType=s,_.operateType==="add"?(o.value={...T},p.value=null):(o.value={...r},p.value={s3Url:r.headImg})}}),(n,s)=>{const r=V("Select"),i=V("Pane");return c(),C(i,{ref_key:"refPane",ref:y,onSubmit:A},{default:a(()=>[d("div",Se,[e(ye,{ref_key:"refForm",ref:g,class:"my-form"},{default:a(()=>[e(r,{placeholder:"请选择供应商",rules:k.provider,mapDictionary:{code:"speak_provider"},modelValue:o.value.provider,"onUpdate:modelValue":s[0]||(s[0]=l=>o.value.provider=l),disabled:h.value},{prepend:a(()=>[d("label",Ie,[u("供应 "),e(Y,null,{default:a(()=>[u("供应商指的是外部服务提供，自己有服务请选择Local")]),_:1})])]),_:1},8,["rules","modelValue","disabled"]),e(M,{type:"text",placeholder:"请输入标识","hide-details":"auto",clearable:"",rules:k.speakName,modelValue:o.value.speakName,"onUpdate:modelValue":s[1]||(s[1]=l=>o.value.speakName=l),disabled:h.value},{prepend:a(()=>[d("label",De,[u("标识 "),o.value.provider==="azure-personal"?(c(),C(Y,{key:0},{default:a(()=>[u("zh-CN-yxh-bozhang75-adrxf")]),_:1})):J("",!0)])]),append:a(()=>[we]),_:1},8,["rules","modelValue","disabled"]),e(M,{type:"text",placeholder:"请输入姓名","hide-details":"auto",clearable:"",rules:k.speakCname,modelValue:o.value.speakCname,"onUpdate:modelValue":s[2]||(s[2]=l=>o.value.speakCname=l)},{prepend:a(()=>[Ue]),_:1},8,["rules","modelValue"]),e(ke,{"hide-details":"auto",modelValue:o.value.gender,"onUpdate:modelValue":s[3]||(s[3]=l=>o.value.gender=l),inline:"",disabled:h.value},{prepend:a(()=>[Ne]),default:a(()=>[(c(!0),U(q,null,O(x.value,l=>(c(),C(be,{label:l.label,color:"primary",value:l.value},null,8,["label","value"]))),256))]),_:1},8,["modelValue","disabled"]),e(r,{placeholder:"请选择年龄段",rules:k.ageGroup,mapDictionary:{code:"speak_age_group"},modelValue:o.value.ageGroup,"onUpdate:modelValue":s[4]||(s[4]=l=>o.value.ageGroup=l)},{prepend:a(()=>[Te]),_:1},8,["rules","modelValue"]),e(r,{placeholder:"请选择语言",rules:k.lang,mapDictionary:{code:"speak_lang"},modelValue:o.value.lang,"onUpdate:modelValue":s[5]||(s[5]=l=>o.value.lang=l),disabled:h.value},{prepend:a(()=>[Fe]),_:1},8,["rules","modelValue","disabled"]),e(r,{placeholder:"请选择风格",rules:k.speakStyle,mapDictionary:{code:"speak_style"},modelValue:o.value.speakStyle,"onUpdate:modelValue":s[6]||(s[6]=l=>o.value.speakStyle=l)},{prepend:a(()=>[Pe]),_:1},8,["rules","modelValue"]),e(r,{placeholder:"请选择适应范围",rules:k.area,mapDictionary:{code:"speak_area"},modelValue:o.value.area,"onUpdate:modelValue":s[7]||(s[7]=l=>o.value.area=l)},{prepend:a(()=>[qe]),_:1},8,["rules","modelValue"]),e(pe,{"hide-details":"auto"},{prepend:a(()=>[Be]),default:a(()=>[p.value&&p.value.s3Url?(c(),C(ge,{key:0,color:"borderColor",variant:"outlined",density:"compact"},{close:a(()=>[e(me,{class:"text-24 opacity-50 cursor-pointer",color:"textPrimary",onClick:R},{default:a(()=>[u("mdi-close-circle")]),_:1})]),default:a(()=>[e(ce,{size:"60"},{default:a(()=>[e(_e,{transition:!1,src:p.value.s3Url,alt:"上传成功后的头像",cover:""},null,8,["src"])]),_:1})]),_:1})):(c(),C(ie,{key:1,accept:"image/*",modelValue:o.value.headImgFileId,"onUpdate:modelValue":s[8]||(s[8]=l=>o.value.headImgFileId=l),infos:p.value,"onUpdate:infos":s[9]||(s[9]=l=>p.value=l),"prepend-icon":null,"prepend-inner-icon":"mdi-camera"},null,8,["modelValue","infos"]))]),_:1}),e(he,{modelValue:o.value.enabled,"onUpdate:modelValue":s[10]||(s[10]=l=>o.value.enabled=l),color:"primary","hide-details":"auto"},{prepend:a(()=>[$e]),_:1},8,["modelValue"]),e(Ve,{modelValue:o.value.remark,"onUpdate:modelValue":s[11]||(s[11]=l=>o.value.remark=l),modelModifiers:{trim:!0},placeholder:"请输入备注",clearable:""},{prepend:a(()=>[Re]),_:1},8,["modelValue"])]),_:1},512)])]),_:1},512)}}});const Ge=fe(Ae,[["__scopeId","data-v-9865707a"]]),Me={class:"text-primary font-weight-black"},ze=d("br",null,null,-1),da=L({__name:"speakerManage",setup(D){const{loadDictTree:B,getLabels:w}=Q();B(["speak_age_group","speak_gender","speak_provider","speak_lang"]);const T=b({title:"发声人管理"}),$=b([{text:"声音合成",disabled:!1,href:"#"},{text:"发声人管理",disabled:!0,href:"#"}]),v=N({speakName:"",provider:null,lang:null}),_=N({list:[],total:0}),o=b(),p=b(),y=b(),g=N({id:"",name:""}),k=r=>{let i=[];return i.push({text:"删除",color:"error",click(){R(r)}}),i.push({text:"编辑",color:"info",click(){s(r)}}),i},h=async(r={})=>{const[i,l]=await G.get({url:"/voice/speak",showLoading:p.value.el,data:{...v,...r}});l?(_.list=l.list||[],_.total=l.total):(_.list=[],_.total=0)},x=()=>{p.value.query({page:1})},R=r=>{g.name=r.speakCname,g.id=r.id,y.value.show({width:"400px",confirmText:g.name})},A=async(r={})=>{const[i,l]=await G.delete({...r,showSuccess:!0,url:`/voice/speak/${g.id}`});l&&(y.value.hide(),h())},n=()=>{o.value.show({title:"添加发声人",operateType:"add"})},s=r=>{o.value.show({title:"编辑发声人",infos:r,operateType:"edit"})};return ee(()=>{h()}),(r,i)=>{const l=V("Select"),E=V("ButtonsInForm"),m=V("el-table-column"),z=V("router-link"),K=V("ButtonsInTable"),j=V("TableWithPager"),H=ae("copy");return c(),U(q,null,[e(se,{title:T.value.title,breadcrumbs:$.value},null,8,["title","breadcrumbs"]),e(re,null,{default:a(()=>[e(Ce,null,{default:a(()=>[e(I,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(M,{modelValue:v.speakName,"onUpdate:modelValue":i[0]||(i[0]=t=>v.speakName=t),label:"请输入标识","hide-details":"",clearable:"",onKeyup:le(x,["enter"]),"onClick:clear":x},null,8,["modelValue","onKeyup"])]),_:1}),e(I,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(l,{modelValue:v.provider,"onUpdate:modelValue":i[1]||(i[1]=t=>v.provider=t),mapDictionary:{code:"speak_provider"},label:"请选择供应商","hide-details":"",onChange:x},null,8,["modelValue"])]),_:1}),e(I,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(l,{modelValue:v.lang,"onUpdate:modelValue":i[2]||(i[2]=t=>v.lang=t),mapDictionary:{code:"speak_lang"},label:"请选择语言","hide-details":"",onChange:x},null,8,["modelValue"])]),_:1}),e(I,{cols:"12",lg:"3",md:"4",sm:"6"},{default:a(()=>[e(E,null,{default:a(()=>[e(ve,{color:"primary",onClick:n},{default:a(()=>[u("添加发声人")]),_:1})]),_:1})]),_:1}),e(I,{cols:"12"},{default:a(()=>[e(ne,null,{default:a(()=>[u("修改之后将实时生效，请谨慎操作！")]),_:1})]),_:1}),e(I,{cols:"12"},{default:a(()=>[e(j,{onQuery:h,ref_key:"tableWithPagerRef",ref:p,infos:_},{default:a(()=>[e(m,{label:"标识",width:"150px","show-overflow-tooltip":""},{default:a(({row:t})=>[te((c(),U("span",null,[u(S(t.speakName),1)])),[[H,t.speakName]])]),_:1}),e(m,{label:"姓名",prop:"speakCname",width:"100px"}),e(m,{label:"供应",width:"100px"},{default:a(({row:t})=>[d("span",null,S(P(w)([["speak_provider",t.provider]])),1)]),_:1}),e(m,{label:"状态",width:"100px"},{default:a(({row:t})=>[e(xe,{modelValue:t.enabled,"onUpdate:modelValue":F=>t.enabled=F},null,8,["modelValue","onUpdate:modelValue"])]),_:1}),e(m,{label:"语音合成",width:"100px"},{default:a(({row:t})=>[t.enabled?(c(),C(z,{key:0,to:{path:"/voice-print/synthesis/synthesis-voice",query:{provider:t.provider,lang:t.lang,speakName:t.speakName}},class:"text-info"},{default:a(()=>[u("合成")]),_:2},1032,["to"])):(c(),U(q,{key:1},[u(" -- ")],64))]),_:1}),e(m,{label:"克隆",width:"100px"},{default:a(({row:t})=>[t.provider==="azure-personal"?(c(),C(z,{key:0,to:{path:"/voice-print/synthesis/speaker/clone-detail",query:{speakName:t.speakName}},class:"text-info"},{default:a(()=>[u("克隆")]),_:2},1032,["to"])):(c(),U(q,{key:1},[u(" -- ")],64))]),_:1}),e(m,{label:"语言",width:"160px"},{default:a(({row:t})=>[d("span",null,S(P(w)([["speak_lang",t.lang]])),1)]),_:1}),e(m,{label:"音色",width:"120px"},{default:a(({row:t})=>[d("div",null,S(P(w)([["speak_age_group",t.ageGroup],["speak_gender",t.gender]],F=>F.length?F.join("")+"声":"未知")),1)]),_:1}),e(m,{label:"试听","min-width":"330px"},{default:a(({row:t})=>[e(ue,{src:t==null?void 0:t.speakDemo},null,8,["src"])]),_:1}),e(m,{label:"备注",prop:"remark","min-width":"200px"}),e(m,{label:"更新时间","min-width":"160px"},{default:a(({row:t})=>[u(S(P(oe).dateFormat(t.updatedAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),e(m,{label:"操作",width:"120px",fixed:"right"},{default:a(({row:t})=>[e(K,{buttons:k(t)},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1}),e(de,{ref_key:"refConfirmDelete",ref:y,onSubmit:A},{text:a(()=>[u(" 您将要删除"),d("span",Me,S(g.name),1),u("发声人，删除之后该声音将无法继续合成新的声音。"),ze,u(" 确定要继续吗？ ")]),_:1},512),e(Ge,{ref_key:"createSpeakerPaneRef",ref:o,onSubmit:x},null,512)],64)}}});export{da as default};
