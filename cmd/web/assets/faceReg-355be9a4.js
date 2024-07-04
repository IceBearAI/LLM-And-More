import{d as A,x as R,r as c,c as z,U as K,k as _,l as S,m as t,W as h,j as e,Y as le,b9 as D,am as H,an as Q,V as ae,_ as u,Z as T,a6 as oe,$ as V,F as $,o as se,S as re,bf as B,Q as ne,n as W,bg as q}from"./utils-5e900427.js";import{_ as ue}from"./BaseBreadcrumb.vue_vue_type_style_index_0_lang-6c45f033.js";import{_ as Y}from"./UiParentCard.vue_vue_type_script_setup_true_lang-80e0344c.js";import{_ as ie}from"./ConfirmByInput.vue_vue_type_style_index_0_lang-41959266.js";import{V as j}from"./VFileInput-c67ccfde.js";import{L as O,a4 as F,_ as Z,V as L,I as de}from"./index-81f01cfc.js";import{V as E}from"./VForm-b09769b4.js";import{V as k}from"./VCol-8780c689.js";import{V as N}from"./VRow-1ab8c2e6.js";import"./Confirm-298ef43f.js";const M=y=>(H("data-v-296d1857"),y=y(),Q(),y),ce={class:"mx-auto mt-3",style:{width:"540px"}},me=M(()=>h("label",{class:"required"},"人脸照片",-1)),pe=M(()=>h("label",{class:"required"},"用户标识",-1)),fe=M(()=>h("label",null,"用户名称",-1)),_e=A({__name:"CreateRegPane",emits:["submit"],setup(y,{expose:U,emit:I}){const m=I,o=R({operateType:"add"}),s=R({file:null,userKey:"",userName:""}),w=c(),g=c(),i=R({userKey:[n=>/^[a-zA-Z0-9-_]+$/.test(n)||"只允许字母、数字、“-” 、“_”"]}),b=z(()=>s.file&&s.file.length>0?URL.createObjectURL(s.file[0]):""),C=async({valid:n,showLoading:r})=>{if(n){const[d,a]=await D.upload({showLoading:r,showSuccess:!0,url:"/esrgan/face/reg",data:s});a&&(w.value.hide(),m("submit"))}};return U({show({title:n,operateType:r,infos:d}){w.value.show({title:n,refForm:g}),o.operateType=r,o.operateType==="add"&&(s.file=null,s.userKey="",s.userName="")}}),(n,r)=>{const d=K("Pane");return _(),S(d,{ref_key:"refPane",ref:w,onSubmit:C},{default:t(()=>[h("div",ce,[e(E,{ref_key:"refForm",ref:g,class:"my-form"},{default:t(()=>[e(j,{modelValue:s.file,"onUpdate:modelValue":r[0]||(r[0]=a=>s.file=a),"prepend-icon":null,accept:"image/*",label:"请上传人脸照片","hide-details":"auto",variant:"outlined",rules:[a=>a&&a.length>0||"请上传人脸照片"]},{prepend:t(()=>[me]),append:t(()=>[b.value?(_(),S(O,{key:0,src:b.value,width:"80px",alt:"previewImageUrl",cover:"",class:"rounded-md align-end text-right"},null,8,["src"])):le("",!0)]),_:1},8,["modelValue","rules"]),e(F,{type:"text",placeholder:"只允许字母、数字、“-” 、“_”","hide-details":"auto",clearable:"",rules:i.userKey,modelValue:s.userKey,"onUpdate:modelValue":r[1]||(r[1]=a=>s.userKey=a)},{prepend:t(()=>[pe]),_:1},8,["rules","modelValue"]),e(F,{type:"text",placeholder:"请输入用户名称","hide-details":"auto",clearable:"",modelValue:s.userName,"onUpdate:modelValue":r[2]||(r[2]=a=>s.userName=a)},{prepend:t(()=>[fe]),_:1},8,["modelValue"])]),_:1},512)])]),_:1},512)}}});const he=Z(_e,[["__scopeId","data-v-296d1857"]]),G=y=>(H("data-v-a9e52492"),y=y(),Q(),y),ye=G(()=>h("label",{class:"required"},"人脸照片",-1)),ge=G(()=>h("label",null,"用户标识",-1)),be=["src"],ve={class:"text-center text-body-1 font-weight-medium mt-2"},we=A({__name:"FaceSearchPane",emits:["search-success"],setup(y,{expose:U,emit:I}){const m=I,o=R({file:null,userKey:""}),s=c(),w=c(),g=c(!1),i=c([]),b=z(()=>o.file&&o.file.length>0?URL.createObjectURL(o.file[0]):""),C=async()=>{let{valid:n}=await w.value.validate();if(n){g.value=!0;const[r,d]=await D.upload({url:"/esrgan/face/search",data:o});d&&(i.value=d.list||[],m("search-success",i.value)),g.value=!1}};return U({show({title:n}){s.value.show({title:n,hasSubmitBtn:!1,width:"1000px"}),o.file=null,o.userKey="",i.value=[]}}),(n,r)=>{const d=K("el-empty"),a=K("Pane");return _(),S(a,{ref_key:"refPane",ref:s},{default:t(()=>[e(N,null,{default:t(()=>[e(k,{cols:"6"},{default:t(()=>[e(Y,{title:"输入"},{default:t(()=>[e(E,{ref_key:"refForm",ref:w,class:"my-form"},{default:t(()=>[e(j,{modelValue:o.file,"onUpdate:modelValue":r[0]||(r[0]=p=>o.file=p),"prepend-icon":null,accept:"image/*",label:"请上传人脸照片","hide-details":"auto",variant:"outlined",rules:[p=>p&&p.length>0||"请上传人脸照片"]},ae({prepend:t(()=>[ye]),_:2},[b.value?{name:"append",fn:t(()=>[e(O,{src:b.value,width:"80px",alt:"previewImageUrl",cover:"",class:"rounded-md align-end text-right"},null,8,["src"])]),key:"0"}:void 0]),1032,["modelValue","rules"]),e(F,{type:"text",placeholder:"只允许字母、数字、“-” 、“_”","hide-details":"auto",clearable:"",modelValue:o.userKey,"onUpdate:modelValue":r[1]||(r[1]=p=>o.userKey=p)},{prepend:t(()=>[ge]),_:1},8,["modelValue"])]),_:1},512),e(L,{color:"primary",block:"",size:"large",flat:"",loading:g.value,onClick:C},{default:t(()=>[u("开始搜索")]),_:1},8,["loading"])]),_:1})]),_:1}),e(k,{cols:"6"},{default:t(()=>[e(Y,{title:"输出"},{default:t(()=>[i.value.length>0?(_(),S(N,{key:0},{default:t(()=>[(_(!0),T($,null,oe(i.value,p=>(_(),S(k,{cols:i.value.length>1?6:12},{default:t(()=>[h("img",{class:"w-full rounded-md align-top",src:p.imgUrl,alt:"人脸图片"},null,8,be),h("p",ve,V(p.dist.toFixed(4)),1)]),_:2},1032,["cols"]))),256))]),_:1})):(_(),S(d,{key:1,"image-size":42}))]),_:1})]),_:1})]),_:1})]),_:1},512)}}});const xe=Z(we,[["__scopeId","data-v-a9e52492"]]),Ve=["src"],ke={class:"text-primary font-weight-black"},Pe=h("br",null,null,-1),Be=A({__name:"faceReg",setup(y){const U=c({title:"人脸注册"}),I=c([{text:"人脸服务",disabled:!1,href:"#"},{text:"人脸注册",disabled:!0,href:"#"}]),m=R({userKey:"",userType:"",userRegion:""}),o=R({list:[],total:0}),s=c(),w=c(),g=c(),i=c(),b=R({uuid:""}),C=x=>{let f=[];return f.push({text:"删除",color:"error",click(){n(x)}}),f},n=x=>{b.uuid=x.uuid,i.value.show({width:"500px",confirmText:b.uuid})},r=async(x={})=>{const[f,P]=await D.delete({...x,showSuccess:!0,url:`/esrgan/face/reg/${b.uuid}`});P&&(i.value.hide(),d())},d=async(x={})=>{const[f,P]=await D.get({url:"/esrgan/face/reg",showLoading:g.value.el,data:{...m,...x}});P?(o.list=P.list||[],o.total=P.total):(o.list=[],o.total=0)},a=()=>{g.value.query({page:1})},p=()=>{s.value.show({title:"注册人脸",operateType:"add"})},J=()=>{w.value.show({title:"人脸搜索"})};return se(()=>{d()}),(x,f)=>{const P=K("ButtonsInForm"),v=K("el-table-column"),X=K("ButtonsInTable"),ee=K("TableWithPager"),te=re("copy");return _(),T($,null,[e(ue,{title:U.value.title,breadcrumbs:I.value},null,8,["title","breadcrumbs"]),e(Y,null,{default:t(()=>[e(N,null,{default:t(()=>[e(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(F,{modelValue:m.userKey,"onUpdate:modelValue":f[0]||(f[0]=l=>m.userKey=l),label:"请输入用户标识","hide-details":"",clearable:"",onKeyup:B(a,["enter"]),"onClick:clear":a},null,8,["modelValue","onKeyup"])]),_:1}),e(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(F,{modelValue:m.userType,"onUpdate:modelValue":f[1]||(f[1]=l=>m.userType=l),label:"请输入用户类型","hide-details":"",clearable:"",onKeyup:B(a,["enter"]),"onClick:clear":a},null,8,["modelValue","onKeyup"])]),_:1}),e(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(F,{modelValue:m.userRegion,"onUpdate:modelValue":f[2]||(f[2]=l=>m.userRegion=l),label:"请输入用户渠道","hide-details":"",clearable:"",onKeyup:B(a,["enter"]),"onClick:clear":a},null,8,["modelValue","onKeyup"])]),_:1}),e(k,{cols:"12",lg:"3",md:"4",sm:"6"},{default:t(()=>[e(P,null,{default:t(()=>[e(L,{color:"primary",onClick:p},{default:t(()=>[u("注册人脸")]),_:1}),e(L,{color:"primary",onClick:J},{default:t(()=>[u("人脸搜索")]),_:1})]),_:1})]),_:1}),e(k,{cols:"12"},{default:t(()=>[e(ee,{onQuery:d,ref_key:"tableWithPagerRef",ref:g,infos:o},{default:t(()=>[e(v,{label:"人脸图片","min-width":"150px"},{default:t(({row:l})=>[e(de,{size:"80",rounded:"md"},{default:t(()=>[h("img",{src:l.inputS3Url,alt:"人脸图片",height:"80"},null,8,Ve)]),_:2},1024)]),_:1}),e(v,{label:"向量库uuid","min-width":"200px","show-overflow-tooltip":""},{default:t(({row:l})=>[ne((_(),T("span",null,[u(V(l.uuid),1)])),[[te,l.uuid]])]),_:1}),e(v,{label:"用户标识","min-width":"100px","show-overflow-tooltip":""},{default:t(({row:l})=>[l.userKey?(_(),T($,{key:0},[u(V(l.userKey),1)],64)):(_(),T($,{key:1},[u(" 未指定 ")],64))]),_:1}),e(v,{label:"用户类型",prop:"userType","min-width":"100px","show-overflow-tooltip":""}),e(v,{label:"用户渠道",prop:"userRegion","min-width":"100px","show-overflow-tooltip":""}),e(v,{label:"py耗时",prop:"durationPy","min-width":"90px"},{default:t(({row:l})=>[u(V(l.durationPy.toFixed(4)),1)]),_:1}),e(v,{label:"接口耗时",prop:"duration","min-width":"90px"},{default:t(({row:l})=>[u(V(l.duration.toFixed(4)),1)]),_:1}),e(v,{label:"创建时间","min-width":"165px"},{default:t(({row:l})=>[u(V(W(q).dateFormat(l.createdAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),e(v,{label:"更新时间","min-width":"165px"},{default:t(({row:l})=>[u(V(W(q).dateFormat(l.updatedAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),e(v,{label:"操作",width:"80px",fixed:"right"},{default:t(({row:l})=>[e(X,{buttons:C(l),onlyOne:""},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),_:1}),e(ie,{ref_key:"refConfirmDelete",ref:i,onSubmit:r},{text:t(()=>[u(" 您将要删除 "),h("span",ke,V(b.uuid),1),u(" 这条记录"),Pe,u(" 确定要继续吗？ ")]),_:1},512),e(he,{ref_key:"createRegPaneRef",ref:s,onSubmit:a},null,512),e(xe,{ref_key:"faceSearchPaneRef",ref:w},null,512)],64)}}});export{Be as default};