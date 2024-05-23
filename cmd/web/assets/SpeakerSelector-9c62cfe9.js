import{ai as M,c as P,k as l,l as h,m as n,j as r,n as t,bt as j,ad as q,r as F,x as A,z as G,o as Z,U as E,Z as f,F as z,W as d,a6 as O,_ as N,$ as g,Y as b,a1 as T,b8 as U,P as W}from"./utils-8a4540e3.js";import{I as C,_ as Y,av as H,Z as $,$ as D,aw as J,a2 as K}from"./index-b4ac7dfa.js";import{_ as Q}from"./AiAudio.vue_vue_type_style_index_0_lang-5b0c997f.js";const X={__name:"IconChecked",props:{size:{type:Number,default:32}},setup(v){M(x=>({f2b5c564:y.value}));const _=v,y=P(()=>-(_.size/2)+"px");return(x,S)=>(l(),h(C,{size:_.size,class:"compo-iconChecked position-absolute flag bg-surface animate__bounceIn"},{default:n(()=>[r(t(j),{color:"#67C23A",size:_.size},null,8,["size"])]),_:1},8,["size"]))}},ee=X;const se={class:"d-flex"},ae={class:"d-flex flex-wrap voice-list flex-1 overflow-auto pt-4 scrollbar-auto"},te=["src","alt"],oe={class:"ml-3 text-body-2 text-black"},le={key:0,class:"text-light"},ne=["src","alt"],re={class:"ml-3 text-body-2"},ce={key:0,class:"text-light"},ie={class:"hv-center"},de={__name:"SpeakerSelector",props:{modelValue:String,infos:Object},emits:["update:modelValue","update:infos"],setup(v,{expose:_,emit:y}){const{loadDictTree:x,getLabels:S}=q(),I=F(),e=A({listSpeaker:[],queryParams:{},optionInfo:{speakName:"",speakCname:"",headImg:"",speakDemo:"",gender:"",ageGroup:"",subTitle:""},style:{isReady:!1}}),{optionInfo:p,style:R}=G(e),B=v,V=y,k=H(B,"modelValue",V),L=a=>S([["speak_age_group",a.ageGroup],["speak_gender",a.gender]],u=>u.length?u.join("")+"声":"未知"),w=async()=>{await x(["speak_age_group","speak_gender"]);const[a,u]=await U.get({showLoading:I.value,url:"/api/voice/speak",data:{pageSize:-1,enabledType:"enabled",...e.queryParams}});if(u){if(e.listSpeaker=u.list.map(o=>({...o,subTitle:L(o)})),e.listSpeaker.length){let o=k.value;if(o){let c=e.listSpeaker.findIndex(i=>i.speakName==o),s=e.listSpeaker[c];if(s)m(s),c>6&&W(()=>{var i;(i=K("#"+o)[0])==null||i.scrollIntoView()});else{window.errorMsg(`未找到初始值对应的数字人 ${o}`);let i=e.listSpeaker[0];m(i)}}else{let c=e.listSpeaker[0];m(c)}}else m({});e.style.isReady=!0}},m=a=>{k.value=a.speakName,e.optionInfo=a,V("update:infos",a)};return Z(()=>{w()}),_({reload(a){e.queryParams=a,w()}}),(a,u)=>{var c;const o=E("NoData");return l(),f("div",{ref_key:"refContainer",ref:I,class:T(["w-100",{"opacity-0":t(R).isReady==!1}])},[((c=e.listSpeaker)==null?void 0:c.length)>0?(l(),f(z,{key:0},[d("div",se,[d("div",ae,[(l(!0),f(z,null,O(e.listSpeaker,s=>(l(),h(D,{id:s.speakName,variant:"outlined",elevation:"0",class:T(["voice-item my-1 mr-5 bg-hover-secondary d-flex align-items",{active:s.speakName===t(k)}]),rounded:"sm",pointer:""},{default:n(()=>[r($,{onClick:i=>m(s),class:"d-flex align-center py-0 px-2"},{default:n(()=>[r(C,{size:"40 "},{default:n(()=>[d("img",{src:s.headImg,alt:s.speakCname,class:"w-100"},null,8,te)]),_:2},1024),d("div",oe,[N(g(s.speakCname),1),s.subTitle?(l(),f("span",le,"（"+g(s.subTitle)+"）",1)):b("",!0)]),s.speakName===t(k)?(l(),h(ee,{key:0})):b("",!0)]),_:2},1032,["onClick"])]),_:2},1032,["id","class"]))),256))])]),r(D,{variant:"outlined",class:"mt-8"},{default:n(()=>[r(J,{class:"d-flex align-center"},{default:n(()=>[r(C,{size:"40 "},{default:n(()=>[d("img",{src:t(p).headImg,alt:t(p).speakCname,class:"w-100"},null,8,ne)]),_:1}),d("div",re,[N(g(t(p).speakCname),1),t(p).subTitle?(l(),f("span",ce,"（"+g(t(p).subTitle)+"）",1)):b("",!0)])]),_:1}),r($,{class:"mt-4"},{default:n(()=>[d("div",ie,[r(Q,{src:t(p).speakDemo,type:"simple"},null,8,["src"])])]),_:1})]),_:1})],64)):(l(),h(o,{key:1}))],2)}}},me=Y(de,[["__scopeId","data-v-9c41b3cf"]]);export{ee as I,me as S};
