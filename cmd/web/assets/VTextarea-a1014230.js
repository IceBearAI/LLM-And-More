import{c as y,r as v,s as p,y as ee,o as te,w as x,K as ae,P as R,j as l,L as S,F as b,Q as G,S as ne,aJ as le}from"./utils-1d42a8ce.js";import{p as oe,ak as ue,az as ie,g as re,be as se,A as ce,aA as de,b as fe,aC as ve,k as xe,an as me,ae as _,aD as ge,aE as he,aF as we,aG as Ve,P as ye}from"./index-19641142.js";const Fe=oe({autoGrow:Boolean,autofocus:Boolean,counter:[Boolean,Number,String],counterValue:Function,prefix:String,placeholder:String,persistentPlaceholder:Boolean,persistentCounter:Boolean,noResize:Boolean,rows:{type:[Number,String],default:5,validator:e=>!isNaN(parseFloat(e))},maxRows:{type:[Number,String],validator:e=>!isNaN(parseFloat(e))},suffix:String,modelModifiers:Object,...ue(),...ie()},"VTextarea"),ke=re()({name:"VTextarea",directives:{Intersect:se},inheritAttrs:!1,props:Fe(),emits:{"click:control":e=>!0,"mousedown:control":e=>!0,"update:focused":e=>!0,"update:modelValue":e=>!0},setup(e,D){let{attrs:F,emit:A,slots:i}=D;const o=ce(e,"modelValue"),{isFocused:f,focus:E,blur:U}=de(e),O=y(()=>typeof e.counterValue=="function"?e.counterValue(o.value):(o.value||"").toString().length),$=y(()=>{if(F.maxlength)return F.maxlength;if(!(!e.counter||typeof e.counter!="number"&&typeof e.counter!="string"))return e.counter});function j(t,n){var a,u;!e.autofocus||!t||(u=(a=n[0].target)==null?void 0:a.focus)==null||u.call(a)}const M=v(),m=v(),z=p(""),g=v(),J=y(()=>e.persistentPlaceholder||f.value||e.active);function P(){var t;g.value!==document.activeElement&&((t=g.value)==null||t.focus()),f.value||E()}function K(t){P(),A("click:control",t)}function L(t){A("mousedown:control",t)}function Q(t){t.stopPropagation(),P(),R(()=>{o.value="",Ve(e["onClick:clear"],t)})}function q(t){var a;const n=t.target;if(o.value=n.value,(a=e.modelModifiers)!=null&&a.trim){const u=[n.selectionStart,n.selectionEnd];R(()=>{n.selectionStart=u[0],n.selectionEnd=u[1]})}}const c=v(),h=v(+e.rows),C=y(()=>["plain","underlined"].includes(e.variant));ee(()=>{e.autoGrow||(h.value=+e.rows)});function d(){e.autoGrow&&R(()=>{if(!c.value||!m.value)return;const t=getComputedStyle(c.value),n=getComputedStyle(m.value.$el),a=parseFloat(t.getPropertyValue("--v-field-padding-top"))+parseFloat(t.getPropertyValue("--v-input-padding-top"))+parseFloat(t.getPropertyValue("--v-field-padding-bottom")),u=c.value.scrollHeight,w=parseFloat(t.lineHeight),k=Math.max(parseFloat(e.rows)*w+a,parseFloat(n.getPropertyValue("--v-input-control-height"))),I=parseFloat(e.maxRows)*w+a||1/0,s=ye(u??0,k,I);h.value=Math.floor((s-a)/w),z.value=xe(s)})}te(d),x(o,d),x(()=>e.rows,d),x(()=>e.maxRows,d),x(()=>e.density,d);let r;return x(c,t=>{t?(r=new ResizeObserver(d),r.observe(c.value)):r==null||r.disconnect()}),ae(()=>{r==null||r.disconnect()}),fe(()=>{const t=!!(i.counter||e.counter||e.counterValue),n=!!(t||i.details),[a,u]=me(F),{modelValue:w,...k}=_.filterProps(e),I=ge(e);return l(_,S({ref:M,modelValue:o.value,"onUpdate:modelValue":s=>o.value=s,class:["v-textarea v-text-field",{"v-textarea--prefixed":e.prefix,"v-textarea--suffixed":e.suffix,"v-text-field--prefixed":e.prefix,"v-text-field--suffixed":e.suffix,"v-textarea--auto-grow":e.autoGrow,"v-textarea--no-resize":e.noResize||e.autoGrow,"v-input--plain-underlined":C.value},e.class],style:e.style},a,k,{centerAffix:h.value===1&&!C.value,focused:f.value}),{...i,default:s=>{let{id:V,isDisabled:B,isDirty:H,isReadonly:W,isValid:X}=s;return l(he,S({ref:m,style:{"--v-textarea-control-height":z.value},onClick:K,onMousedown:L,"onClick:clear":Q,"onClick:prependInner":e["onClick:prependInner"],"onClick:appendInner":e["onClick:appendInner"]},I,{id:V.value,active:J.value||H.value,centerAffix:h.value===1&&!C.value,dirty:H.value||e.dirty,disabled:B.value,focused:f.value,error:X.value===!1}),{...i,default:Y=>{let{props:{class:N,...T}}=Y;return l(b,null,[e.prefix&&l("span",{class:"v-text-field__prefix"},[e.prefix]),G(l("textarea",S({ref:g,class:N,value:o.value,onInput:q,autofocus:e.autofocus,readonly:W.value,disabled:B.value,placeholder:e.placeholder,rows:e.rows,name:e.name,onFocus:P,onBlur:U},T,u),null),[[ne("intersect"),{handler:j},null,{once:!0}]]),e.autoGrow&&G(l("textarea",{class:[N,"v-textarea__sizer"],id:`${T.id}-sizer`,"onUpdate:modelValue":Z=>o.value=Z,ref:c,readonly:!0,"aria-hidden":"true"},null),[[le,o.value]]),e.suffix&&l("span",{class:"v-text-field__suffix"},[e.suffix])])}})},details:n?s=>{var V;return l(b,null,[(V=i.details)==null?void 0:V.call(i,s),t&&l(b,null,[l("span",null,null),l(we,{active:e.persistentCounter||f.value,value:O.value,max:$.value},i.counter)])])}:void 0})}),ve({},M,m,g)}});export{ke as V};
