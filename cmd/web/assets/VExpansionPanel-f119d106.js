import{p as m,m as f,R as _,a as g,r as A,g as P,S as R,s as D,G as E,b as y,au as S,av as G,O as L,a6 as b,aw as $,v as k,af as j,o as z,T as O,q,U as F,w as H,y as N}from"./index-b84f97ad.js";import{c as d,H as p,j as l,i as I,Q as C,a3 as Q,S as U,D as J}from"./utils-fc8ebe1f.js";const r=Symbol.for("vuetify:v-expansion-panel"),K=["default","accordion","inset","popout"],M=m({color:String,static:Boolean,variant:{type:String,default:"default",validator:e=>K.includes(e)},readonly:Boolean,...f(),..._(),...g(),...A()},"VExpansionPanels"),ne=P()({name:"VExpansionPanels",props:M(),emits:{"update:modelValue":e=>!0},setup(e,o){let{slots:n}=o;R(e,r);const{themeClasses:a}=D(e),i=d(()=>e.variant&&`v-expansion-panels--variant-${e.variant}`);return E({VExpansionPanel:{color:p(e,"color")},VExpansionPanelTitle:{readonly:p(e,"readonly"),static:p(e,"static")}}),y(()=>l(e.tag,{class:["v-expansion-panels",a.value,i.value,e.class],style:e.style},n)),{}}}),W=m({...f(),...S()},"VExpansionPanelText"),X=P()({name:"VExpansionPanelText",props:W(),setup(e,o){let{slots:n}=o;const a=I(r);if(!a)throw new Error("[Vuetify] v-expansion-panel-text needs to be placed inside v-expansion-panel");const{hasContent:i,onAfterLeave:u}=G(e,a.isSelected);return y(()=>l(L,{onAfterLeave:u},{default:()=>{var c;return[C(l("div",{class:["v-expansion-panel-text",e.class],style:e.style},[n.default&&i.value&&l("div",{class:"v-expansion-panel-text__wrapper"},[(c=n.default)==null?void 0:c.call(n)])]),[[Q,a.isSelected.value]])]}})),{}}}),T=m({color:String,expandIcon:{type:b,default:"$expand"},collapseIcon:{type:b,default:"$collapse"},hideActions:Boolean,static:Boolean,ripple:{type:[Boolean,Object],default:!1},readonly:Boolean,...f()},"VExpansionPanelTitle"),Y=P()({name:"VExpansionPanelTitle",directives:{Ripple:$},props:T(),setup(e,o){let{slots:n}=o;const a=I(r);if(!a)throw new Error("[Vuetify] v-expansion-panel-title needs to be placed inside v-expansion-panel");const{backgroundColorClasses:i,backgroundColorStyles:u}=k(e,"color"),c=d(()=>({collapseIcon:e.collapseIcon,disabled:a.disabled.value,expanded:a.isSelected.value,expandIcon:e.expandIcon,readonly:e.readonly}));return y(()=>{var v;return C(l("button",{class:["v-expansion-panel-title",{"v-expansion-panel-title--active":a.isSelected.value,"v-expansion-panel-title--static":e.static},i.value,e.class],style:[u.value,e.style],type:"button",tabindex:a.disabled.value?-1:void 0,disabled:a.disabled.value,"aria-expanded":a.isSelected.value,onClick:e.readonly?void 0:a.toggle},[l("span",{class:"v-expansion-panel-title__overlay"},null),(v=n.default)==null?void 0:v.call(n,c.value),!e.hideActions&&l("span",{class:"v-expansion-panel-title__icon"},[n.actions?n.actions(c.value):l(j,{icon:a.isSelected.value?e.collapseIcon:e.expandIcon},null)])]),[[U("ripple"),e.ripple]])}),{}}}),Z=m({title:String,text:String,bgColor:String,...f(),...z(),...O(),...S(),...q(),...g(),...T()},"VExpansionPanel"),le=P()({name:"VExpansionPanel",props:Z(),emits:{"group:selected":e=>!0},setup(e,o){let{slots:n}=o;const a=F(e,r),{backgroundColorClasses:i,backgroundColorStyles:u}=k(e,"bgColor"),{elevationClasses:c}=H(e),{roundedClasses:v}=N(e),h=d(()=>(a==null?void 0:a.disabled.value)||e.disabled),V=d(()=>a.group.items.value.reduce((s,t,x)=>(a.group.selected.value.includes(t.id)&&s.push(x),s),[])),w=d(()=>{const s=a.group.items.value.findIndex(t=>t.id===a.id);return!a.isSelected.value&&V.value.some(t=>t-s===1)}),B=d(()=>{const s=a.group.items.value.findIndex(t=>t.id===a.id);return!a.isSelected.value&&V.value.some(t=>t-s===-1)});return J(r,a),E({VExpansionPanelText:{eager:p(e,"eager")}}),y(()=>{const s=!!(n.text||e.text),t=!!(n.title||e.title);return l(e.tag,{class:["v-expansion-panel",{"v-expansion-panel--active":a.isSelected.value,"v-expansion-panel--before-active":w.value,"v-expansion-panel--after-active":B.value,"v-expansion-panel--disabled":h.value},v.value,i.value,e.class],style:[u.value,e.style]},{default:()=>{var x;return[l("div",{class:["v-expansion-panel__shadow",...c.value]},null),t&&l(Y,{key:"title",collapseIcon:e.collapseIcon,color:e.color,expandIcon:e.expandIcon,hideActions:e.hideActions,ripple:e.ripple},{default:()=>[n.title?n.title():e.title]}),s&&l(X,{key:"text"},{default:()=>[n.text?n.text():e.text]}),(x=n.default)==null?void 0:x.call(n)]}})}),{}}});export{ne as V,le as a,Y as b,X as c};
