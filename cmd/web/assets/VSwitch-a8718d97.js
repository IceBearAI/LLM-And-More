import{p as T,al as $,aj as p,g as M,A as m,bc as q,aA as E,ab as G,b as H,ao as J,ae as g,ak as V,b3 as K,af as O,bd as Q,ad as W}from"./index-b84f97ad.js";import{r as X,c as h,j as a,L as b,F as Y}from"./utils-fc8ebe1f.js";const Z=T({indeterminate:Boolean,inset:Boolean,flat:Boolean,loading:{type:[Boolean,String],default:!1},...$(),...p()},"VSwitch"),te=M()({name:"VSwitch",inheritAttrs:!1,props:Z(),emits:{"update:focused":e=>!0,"update:modelValue":e=>!0,"update:indeterminate":e=>!0},setup(e,k){let{attrs:C,slots:o}=k;const n=m(e,"indeterminate"),s=m(e,"modelValue"),{loaderClasses:w}=q(e),{isFocused:y,focus:S,blur:P}=E(e),f=X(),A=h(()=>typeof e.loading=="string"&&e.loading!==""?e.loading:e.color),F=G(),_=h(()=>e.id||`switch-${F}`);function x(){n.value&&(n.value=!1)}function B(i){var u,d;i.stopPropagation(),i.preventDefault(),(d=(u=f.value)==null?void 0:u.input)==null||d.click()}return H(()=>{const[i,u]=J(C),d=g.filterProps(e),I=V.filterProps(e);return a(g,b({class:["v-switch",{"v-switch--inset":e.inset},{"v-switch--indeterminate":n.value},w.value,e.class]},i,d,{modelValue:s.value,"onUpdate:modelValue":r=>s.value=r,id:_.value,focused:y.value,style:e.style}),{...o,default:r=>{let{id:L,messagesId:U,isDisabled:j,isReadonly:z,isValid:D}=r;return a(V,b({ref:f},I,{modelValue:s.value,"onUpdate:modelValue":[l=>s.value=l,x],id:L.value,"aria-describedby":U.value,type:"checkbox","aria-checked":n.value?"mixed":void 0,disabled:j.value,readonly:z.value,onFocus:S,onBlur:P},u),{...o,default:l=>{let{backgroundColorClasses:c,backgroundColorStyles:t}=l;return a("div",{class:["v-switch__track",...c.value],style:t.value,onClick:B},null)},input:l=>{let{inputNode:c,icon:t,backgroundColorClasses:N,backgroundColorStyles:R}=l;return a(Y,null,[c,a("div",{class:["v-switch__thumb",{"v-switch__thumb--filled":t||e.loading},e.inset?void 0:N.value],style:e.inset?void 0:R.value},[a(K,null,{default:()=>[e.loading?a(Q,{name:"v-switch",active:!0,color:D.value===!1?void 0:A.value},{default:v=>o.loader?o.loader(v):a(W,{active:v.isActive,color:v.color,indeterminate:!0,size:"16",width:"2"},null)}):t&&a(O,{key:t,icon:t,size:"x-small"},null)]})])])}})}})}),{}}});export{te as V};