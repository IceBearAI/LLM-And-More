import{aF as k,Z as x,a0 as y,P as v,V as m,$ as w,a1 as $}from"./index-1ef416e3.js";import{x as B,z as S,aa as g,k as o,l as u,m as t,j as s,n as l,X as i,Y as z,Z as D,_ as f,F}from"./utils-5f930323.js";const N={__name:"Confirm",emits:["close","submit"],setup(T,{expose:V,emit:p}){const e=B({style:{width:"auto"},formData:{},visible:!1}),{style:b,formData:P}=S(e),r=g(),n=p,c=()=>{e.visible=!1},C=()=>{c(),n("close")},_=()=>{n("submit")};return V({show({width:a="auto"}={}){e.style.width=a,e.visible=!0},hide(){c()}}),(a,d)=>(o(),u($,{modelValue:e.visible,"onUpdate:modelValue":d[0]||(d[0]=h=>e.visible=h),"max-width":"800",width:l(b).width},{default:t(()=>[s(w,{class:"px-2 pt-3"},{default:t(()=>[l(r).title?(o(),u(k,{key:0,class:"text-subtitle-1 color-font"},{default:t(()=>[i(a.$slots,"title")]),_:3})):z("",!0),s(x,{class:"text-body-1 color-font-light"},{default:t(()=>[i(a.$slots,"text")]),_:3}),s(y,null,{default:t(()=>[s(v),l(r).buttons?i(a.$slots,"buttons",{key:0}):(o(),D(F,{key:1},[s(m,{size:"small",color:"secondary",variant:"outlined",onClick:C},{default:t(()=>[f("取消")]),_:1}),s(m,{size:"small",color:"primary",variant:"flat",onClick:_},{default:t(()=>[f("确定")]),_:1})],64))]),_:3})]),_:3})]),_:3},8,["modelValue","width"]))}},A=N;export{A as C};
