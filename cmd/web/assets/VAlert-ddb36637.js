import{b9 as L,p as $,a6 as D,m as z,a7 as R,aQ as w,o as F,aR as j,aS as E,q as N,a as O,r as p,ba as q,g as H,A as M,s as Q,bb as U,b3 as G,aT as J,w as K,aU as W,aV as X,y as Y,aI as Z,a9 as ee,bc as ae,af as te,N as d,V as le}from"./index-733db90d.js";import{c as o,H as ne,j as t,L as se}from"./utils-77c08ab5.js";const oe=L("v-alert-title"),re=["success","info","warning","error"],ie=$({border:{type:[Boolean,String],validator:e=>typeof e=="boolean"||["top","end","bottom","start"].includes(e)},borderColor:String,closable:Boolean,closeIcon:{type:D,default:"$close"},closeLabel:{type:String,default:"$vuetify.close"},icon:{type:[Boolean,String,Function,Object],default:null},modelValue:{type:Boolean,default:!0},prominent:Boolean,title:String,text:String,type:{type:String,validator:e=>re.includes(e)},...z(),...R(),...w(),...F(),...j(),...E(),...N(),...O(),...p(),...q({variant:"flat"})},"VAlert"),de=H()({name:"VAlert",props:ie(),emits:{"click:close":e=>!0,"update:modelValue":e=>!0},setup(e,v){let{emit:m,slots:a}=v;const r=M(e,"modelValue"),n=o(()=>{if(e.icon!==!1)return e.type?e.icon??`$${e.type}`:e.icon}),y=o(()=>({color:e.color??e.type,variant:e.variant})),{themeClasses:b}=Q(e),{colorClasses:f,colorStyles:k,variantClasses:V}=U(y),{densityClasses:P}=G(e),{dimensionStyles:C}=J(e),{elevationClasses:g}=K(e),{locationStyles:S}=W(e),{positionClasses:x}=X(e),{roundedClasses:_}=Y(e),{textColorClasses:A,textColorStyles:T}=Z(ne(e,"borderColor")),{t:B}=ee(),i=o(()=>({"aria-label":B(e.closeLabel),onClick(s){r.value=!1,m("click:close",s)}}));return()=>{const s=!!(a.prepend||n.value),I=!!(a.title||e.title),h=!!(a.close||e.closable);return r.value&&t(e.tag,{class:["v-alert",e.border&&{"v-alert--border":!!e.border,[`v-alert--border-${e.border===!0?"start":e.border}`]:!0},{"v-alert--prominent":e.prominent},b.value,f.value,P.value,g.value,x.value,_.value,V.value,e.class],style:[k.value,C.value,S.value,e.style],role:"alert"},{default:()=>{var c,u;return[ae(!1,"v-alert"),e.border&&t("div",{key:"border",class:["v-alert__border",A.value],style:T.value},null),s&&t("div",{key:"prepend",class:"v-alert__prepend"},[a.prepend?t(d,{key:"prepend-defaults",disabled:!n.value,defaults:{VIcon:{density:e.density,icon:n.value,size:e.prominent?44:28}}},a.prepend):t(te,{key:"prepend-icon",density:e.density,icon:n.value,size:e.prominent?44:28},null)]),t("div",{class:"v-alert__content"},[I&&t(oe,{key:"title"},{default:()=>{var l;return[((l=a.title)==null?void 0:l.call(a))??e.title]}}),((c=a.text)==null?void 0:c.call(a))??e.text,(u=a.default)==null?void 0:u.call(a)]),a.append&&t("div",{key:"append",class:"v-alert__append"},[a.append()]),h&&t("div",{key:"close",class:"v-alert__close"},[a.close?t(d,{key:"close-defaults",defaults:{VBtn:{icon:e.closeIcon,size:"x-small",variant:"text"}}},{default:()=>{var l;return[(l=a.close)==null?void 0:l.call(a,{props:i.value})]}}):t(le,se({key:"close-btn",icon:e.closeIcon,size:"x-small",variant:"text"},i.value),null)])]}})}}});export{de as V};