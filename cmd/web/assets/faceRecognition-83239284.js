import{_ as L}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-02f17eb8.js";import{_ as P}from"./UiParentCard.vue_vue_type_script_setup_true_lang-bc14ee23.js";import{C as A}from"./ChipBoolean-45755643.js";import{d as W,x as B,r as h,c as Q,U as x,k as c,l as f,m as e,j as t,W as r,_ as u,Z as M,$ as g,F,Y as v,a0 as j,av as Y,al as H,am as O,o as J,n as Z,aB as G}from"./utils-fc8ebe1f.js";import{_ as K}from"./Explain.vue_vue_type_style_index_0_lang-2aca43e9.js";import{V as k}from"./VCol-f5dbb3fc.js";import{V as X}from"./VForm-cdbc2f03.js";import{V as q}from"./VFileInput-6c8fbbc5.js";import{M as N,a5 as ee,V as z,ae as V,_ as te,J as D}from"./index-b84f97ad.js";import{V as ae}from"./VSlider-bf749123.js";import{V as $}from"./VRow-4fd670bb.js";import"./IconInfoCircle-d065b385.js";const p=I=>(H("data-v-b78e16db"),I=I(),O(),I),le=p(()=>r("label",{class:"required"},"检测类型",-1)),oe=p(()=>r("label",{class:"required"},"检测图片",-1)),ne=p(()=>r("label",{class:"required"},"比对图片",-1)),se={class:"text-center",style:{width:"28px"}},re=p(()=>r("label",{class:"required"},"比对阈值",-1)),de=p(()=>r("label",null,"检测图片：",-1)),ie=["src"],ue=p(()=>r("label",null,"人脸边距：",-1)),ce=p(()=>r("label",null,"比对图片：",-1)),pe=["src"],me=p(()=>r("label",null,"人脸个数：",-1)),_e=p(()=>r("label",null,"比对阈值：",-1)),fe=p(()=>r("label",null,"同一个人：",-1)),ge=W({__name:"CreateRecognitionPane",emits:["submit"],setup(I,{expose:T,emit:R}){const b=R,m=B({operateType:"add"}),y=h(),S=h(),l=B({checkType:null,checkImage:[],baseImage:[],tolerance:.6,faceMargin:""}),n=h({}),w=h(!1),U=Q(()=>{const i={checkImage:"",baseImage:""};return["checkImage","baseImage"].forEach(a=>{l[a]&&l[a].length>0&&(i[a]=URL.createObjectURL(l[a][0]))}),i}),_=async()=>{let{valid:i}=await S.value.validate();if(i){w.value=!0;const[a,s]=await Y.upload({url:"/esrgan/face/recognition",data:l});s&&(n.value=s,b("submit")),w.value=!1}};return T({show({title:i,operateType:a,infos:s}){y.value.show({title:i,width:a==="add"?"1000px":"600px",hasSubmitBtn:!1}),m.operateType=a,m.operateType==="add"?(l.checkType=null,l.checkImage=[],l.baseImage=[],l.tolerance=.6,n.value={}):n.value=s}}),(i,a)=>{const s=x("Select"),C=x("Pane");return c(),f(C,{ref_key:"refPane",ref:y},{default:e(()=>[t($,null,{default:e(()=>[m.operateType==="add"?(c(),f(k,{key:0,cols:"12",md:"7"},{default:e(()=>[t(P,{title:"输入"},{default:e(()=>[t(X,{ref_key:"refForm",ref:S,class:"my-form"},{default:e(()=>[t(s,{placeholder:"请选择检测类型",rules:[o=>!!o||"请选择检测类型"],mapDictionary:{code:"face_check_type"},modelValue:l.checkType,"onUpdate:modelValue":a[0]||(a[0]=o=>l.checkType=o)},{prepend:e(()=>[le]),_:1},8,["rules","modelValue"]),t(q,{modelValue:l.checkImage,"onUpdate:modelValue":a[1]||(a[1]=o=>l.checkImage=o),"prepend-icon":null,accept:"image/*",label:"请上传检测图片","hide-details":"auto",variant:"outlined",rules:[o=>o.length>0||"请上传检测图片"]},{prepend:e(()=>[oe]),append:e(()=>[t(N,{src:U.value.checkImage,width:"80px",alt:"previewImageUrl",cover:"",class:"rounded-md align-end text-right"},null,8,["src"])]),_:1},8,["modelValue","rules"]),t(ee,{type:"text",placeholder:"格式：top, right, bottom, left","hide-details":"auto",clearable:"",modelValue:l.faceMargin,"onUpdate:modelValue":a[2]||(a[2]=o=>l.faceMargin=o),rules:[o=>!o||/^(\d{1,4},\d{1,4},\d{1,4},\d{1,4})$/.test(o)||"请输入正确格式：top, right, bottom, left"]},{prepend:e(()=>[r("label",null,[u("人脸边距 "),t(K,null,{default:e(()=>[u("人脸边距约束，格式：top, right, bottom, left")]),_:1})])]),_:1},8,["modelValue","rules"]),l.checkType===2?(c(),M(F,{key:0},[t(q,{modelValue:l.baseImage,"onUpdate:modelValue":a[3]||(a[3]=o=>l.baseImage=o),"prepend-icon":null,accept:"image/*",label:"请上传比对图片","hide-details":"auto",variant:"outlined",rules:[o=>o.length>0||"请上传比对图片"]},{prepend:e(()=>[ne]),append:e(()=>[t(N,{src:U.value.baseImage,width:"80px",alt:"previewImageUrl",cover:"",class:"rounded-md align-end text-right"},null,8,["src"])]),_:1},8,["modelValue","rules"]),t(ae,{class:"mx-0",modelValue:l.tolerance,"onUpdate:modelValue":a[4]||(a[4]=o=>l.tolerance=o),color:"primary",max:1,min:0,step:.1,"hide-details":"auto","thumb-label":""},{append:e(()=>[r("div",se,g(l.tolerance),1)]),prepend:e(()=>[re]),_:1},8,["modelValue"])],64)):v("",!0)]),_:1},512),t(z,{color:"primary",block:"",size:"large",flat:"",loading:w.value,onClick:_},{default:e(()=>[u("开始检测")]),_:1},8,["loading"])]),_:1})]),_:1})):v("",!0),t(k,{class:j({"result-right":m.operateType==="edit"}),cols:"12",md:m.operateType==="add"?5:0},{default:e(()=>[t(P,{title:"输出"},{default:e(()=>[n.value.inputS3Url?(c(),f(V,{key:0},{prepend:e(()=>[de]),default:e(()=>[r("img",{src:n.value.inputS3Url,width:"200",alt:"检测图片",class:"rounded-md align-end text-right"},null,8,ie)]),_:1})):v("",!0),n.value.faceMargin?(c(),f(V,{key:1},{prepend:e(()=>[ue]),default:e(()=>[u(" "+g(n.value.faceMargin),1)]),_:1})):v("",!0),n.value.outputS3Url?(c(),f(V,{key:2},{prepend:e(()=>[ce]),default:e(()=>[r("img",{src:n.value.outputS3Url,width:"200",alt:"比对图片",class:"rounded-md align-end text-right"},null,8,pe)]),_:1})):v("",!0),t(V,null,{prepend:e(()=>[me]),default:e(()=>[u(" "+g(n.value.faceNum),1)]),_:1}),n.value.outputS3Url?(c(),f(V,{key:3},{prepend:e(()=>[_e]),default:e(()=>[u(" "+g(n.value.denoiseStrength),1)]),_:1})):v("",!0),l.checkType===2||n.value.outputS3Url?(c(),f(V,{key:4},{prepend:e(()=>[fe]),default:e(()=>[u(" "+g(n.value.isSame===void 0?"":n.value.isSame?"是":"否"),1)]),_:1})):v("",!0)]),_:1})]),_:1},8,["class","md"])]),_:1})]),_:1},512)}}});const he=te(ge,[["__scopeId","data-v-b78e16db"]]),be=["src"],ve=["src"],Be=W({__name:"faceRecognition",setup(I){const T=h({title:"人脸检测"}),R=h([{text:"图像服务",disabled:!1,href:"#"},{text:"人脸检测",disabled:!0,href:"#"}]),b=B({list:[],total:0}),m=h(),y=h(),S=_=>{let i=[];return i.push({text:"查看",color:"info",click(){U(_)}}),i},l=async(_={})=>{const[i,a]=await Y.get({url:"/esrgan/list",showLoading:y.value.el,data:{modelType:"faceRecognition",..._}});a?(b.list=a.list||[],b.total=a.total):(b.list=[],b.total=0)},n=()=>{y.value.query({page:1})},w=()=>{m.value.show({title:"创建检测",operateType:"add"})},U=_=>{m.value.show({title:"查看",infos:_,operateType:"edit"})};return J(()=>{l()}),(_,i)=>{const a=x("ButtonsInForm"),s=x("el-table-column"),C=x("ButtonsInTable"),o=x("TableWithPager");return c(),M(F,null,[t(L,{title:T.value.title,breadcrumbs:R.value},null,8,["title","breadcrumbs"]),t(P,null,{default:e(()=>[t($,null,{default:e(()=>[t(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:e(()=>[t(a,null,{default:e(()=>[t(z,{color:"primary",onClick:w},{default:e(()=>[u("创建检测")]),_:1})]),_:1})]),_:1}),t(k,{cols:"12"},{default:e(()=>[t(o,{onQuery:l,ref_key:"tableWithPagerRef",ref:y,infos:b},{default:e(()=>[t(s,{label:"检测图片","min-width":"150px"},{default:e(({row:d})=>[t(D,{size:"80",rounded:"md"},{default:e(()=>[r("img",{src:d.inputS3Url,alt:"检测图片",height:"80"},null,8,be)]),_:2},1024)]),_:1}),t(s,{label:"比对图片","min-width":"150px"},{default:e(({row:d})=>[d.outputS3Url?(c(),f(D,{key:0,size:"80",rounded:"md"},{default:e(()=>[r("img",{src:d.outputS3Url,alt:"比对图片",height:"80"},null,8,ve)]),_:2},1024)):(c(),M(F,{key:1},[u(" -- ")],64))]),_:1}),t(s,{label:"人脸个数",prop:"faceNum","min-width":"160px"}),t(s,{label:"比对阈值",prop:"denoiseStrength","min-width":"160px"}),t(s,{label:"人脸边距",prop:"faceMargin","min-width":"160px"},{default:e(({row:d})=>[t($,{class:"text-left",dense:""},{default:e(()=>[t(k,{cols:"12"},{default:e(()=>[u("req: "+g(d.faceMargin),1)]),_:2},1024),t(k,{cols:"12"},{default:e(()=>[u("real: "+g(d.faceMarginReal),1)]),_:2},1024)]),_:2},1024)]),_:1}),t(s,{label:"是否同一个人","min-width":"120px"},{default:e(({row:d})=>[t(A,{modelValue:d.isSame,"onUpdate:modelValue":E=>d.isSame=E},null,8,["modelValue","onUpdate:modelValue"])]),_:1}),t(s,{label:"操作人",prop:"operatorEmail","min-width":"150px","show-overflow-tooltip":""}),t(s,{label:"创建时间","min-width":"165px"},{default:e(({row:d})=>[u(g(Z(G).dateFormat(d.createdAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),t(s,{label:"操作",width:"80px",fixed:"right"},{default:e(({row:d})=>[t(C,{buttons:S(d),onlyOne:""},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1}),t(he,{ref_key:"createRecognitionPaneRef",ref:m,onSubmit:n},null,512)],64)}}});export{Be as default};
