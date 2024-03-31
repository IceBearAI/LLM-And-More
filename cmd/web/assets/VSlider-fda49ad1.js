import{c as s,H as z,s as le,r as G,D as he,i as ie,j as n,Q as ae,S as ke,a3 as ye,F as Se,L as pe}from"./utils-77c08ab5.js";import{p as H,q as ge,o as Ve,aH as se,u as J,aa as we,P as oe,m as ue,g as Q,ar as Ce,aI as _e,b as W,k as K,w as Te,aJ as xe,aK as Fe,y as Pe,v as re,aG as ze,ak as Le,A as Me,az as Re,ae as ne,a4 as Ee}from"./index-733db90d.js";const Z=Symbol.for("vuetify:v-slider");function Ne(e,t,a){const i=a==="vertical",o=t.getBoundingClientRect(),h="touches"in e?e.touches[0]:e;return i?h.clientY-(o.top+o.height/2):h.clientX-(o.left+o.width/2)}function Be(e,t){return"touches"in e&&e.touches.length?e.touches[0][t]:"changedTouches"in e&&e.changedTouches.length?e.changedTouches[0][t]:e[t]}const De=H({disabled:{type:Boolean,default:null},error:Boolean,readonly:{type:Boolean,default:null},max:{type:[Number,String],default:100},min:{type:[Number,String],default:0},step:{type:[Number,String],default:0},thumbColor:String,thumbLabel:{type:[Boolean,String],default:void 0,validator:e=>typeof e=="boolean"||e==="always"},thumbSize:{type:[Number,String],default:20},showTicks:{type:[Boolean,String],default:!1,validator:e=>typeof e=="boolean"||e==="always"},ticks:{type:[Array,Object]},tickSize:{type:[Number,String],default:2},color:String,trackColor:String,trackFillColor:String,trackSize:{type:[Number,String],default:4},direction:{type:String,default:"horizontal",validator:e=>["vertical","horizontal"].includes(e)},reverse:Boolean,...ge(),...Ve({elevation:2}),ripple:{type:Boolean,default:!0}},"Slider"),Ie=e=>{const t=s(()=>parseFloat(e.min)),a=s(()=>parseFloat(e.max)),i=s(()=>+e.step>0?parseFloat(e.step):0),o=s(()=>Math.max(se(i.value),se(t.value)));function h(k){if(k=parseFloat(k),i.value<=0)return k;const c=oe(k,t.value,a.value),g=t.value%i.value,C=Math.round((c-g)/i.value)*i.value+g;return parseFloat(Math.min(C,a.value).toFixed(o.value))}return{min:t,max:a,step:i,decimals:o,roundValue:h}},qe=e=>{let{props:t,steps:a,onSliderStart:i,onSliderMove:o,onSliderEnd:h,getActiveThumb:k}=e;const{isRtl:c}=J(),g=z(t,"reverse"),C=s(()=>t.direction==="vertical"),T=s(()=>C.value!==g.value),{min:f,max:V,step:x,decimals:E,roundValue:L}=a,I=s(()=>parseInt(t.thumbSize,10)),N=s(()=>parseInt(t.tickSize,10)),M=s(()=>parseInt(t.trackSize,10)),F=s(()=>(V.value-f.value)/x.value),B=z(t,"disabled"),_=s(()=>t.error||t.disabled?void 0:t.thumbColor??t.color),u=s(()=>t.error||t.disabled?void 0:t.trackColor??t.color),w=s(()=>t.error||t.disabled?void 0:t.trackFillColor??t.color),v=le(!1),m=le(0),y=G(),S=G();function r(l){var te;const d=t.direction==="vertical",de=d?"top":"left",ce=d?"height":"width",ve=d?"clientY":"clientX",{[de]:me,[ce]:be}=(te=y.value)==null?void 0:te.$el.getBoundingClientRect(),fe=Be(l,ve);let Y=Math.min(Math.max((fe-me-m.value)/be,0),1)||0;return(d?T.value:T.value!==c.value)&&(Y=1-Y),L(f.value+Y*(V.value-f.value))}const P=l=>{h({value:r(l)}),v.value=!1,m.value=0},O=l=>{S.value=k(l),S.value&&(S.value.focus(),v.value=!0,S.value.contains(l.target)?m.value=Ne(l,S.value,t.direction):(m.value=0,o({value:r(l)})),i({value:r(l)}))},R={passive:!0,capture:!0};function $(l){o({value:r(l)})}function j(l){l.stopPropagation(),l.preventDefault(),P(l),window.removeEventListener("mousemove",$,R),window.removeEventListener("mouseup",j)}function b(l){var d;P(l),window.removeEventListener("touchmove",$,R),(d=l.target)==null||d.removeEventListener("touchend",b)}function p(l){var d;O(l),window.addEventListener("touchmove",$,R),(d=l.target)==null||d.addEventListener("touchend",b,{passive:!1})}function D(l){l.preventDefault(),O(l),window.addEventListener("mousemove",$,R),window.addEventListener("mouseup",j,{passive:!1})}const q=l=>{const d=(l-f.value)/(V.value-f.value)*100;return oe(isNaN(d)?0:d,0,100)},A=z(t,"showTicks"),U=s(()=>A.value?t.ticks?Array.isArray(t.ticks)?t.ticks.map(l=>({value:l,position:q(l),label:l.toString()})):Object.keys(t.ticks).map(l=>({value:parseFloat(l),position:q(parseFloat(l)),label:t.ticks[l]})):F.value!==1/0?we(F.value+1).map(l=>{const d=f.value+l*x.value;return{value:d,position:q(d)}}):[]:[]),X=s(()=>U.value.some(l=>{let{label:d}=l;return!!d})),ee={activeThumbRef:S,color:z(t,"color"),decimals:E,disabled:B,direction:z(t,"direction"),elevation:z(t,"elevation"),hasLabels:X,isReversed:g,indexFromEnd:T,min:f,max:V,mousePressed:v,numTicks:F,onSliderMousedown:D,onSliderTouchstart:p,parsedTicks:U,parseMouseMove:r,position:q,readonly:z(t,"readonly"),rounded:z(t,"rounded"),roundValue:L,showTicks:A,startOffset:m,step:x,thumbSize:I,thumbColor:_,thumbLabel:z(t,"thumbLabel"),ticks:z(t,"ticks"),tickSize:N,trackColor:u,trackContainerRef:y,trackFillColor:w,trackSize:M,vertical:C};return he(Z,ee),ee},Ke=H({focused:Boolean,max:{type:Number,required:!0},min:{type:Number,required:!0},modelValue:{type:Number,required:!0},position:{type:Number,required:!0},ripple:{type:[Boolean,Object],default:!0},...ue()},"VSliderThumb"),Oe=Q()({name:"VSliderThumb",directives:{Ripple:Ce},props:Ke(),emits:{"update:modelValue":e=>!0},setup(e,t){let{slots:a,emit:i}=t;const o=ie(Z),{isRtl:h,rtlClasses:k}=J();if(!o)throw new Error("[Vuetify] v-slider-thumb must be used inside v-slider or v-range-slider");const{thumbColor:c,step:g,disabled:C,thumbSize:T,thumbLabel:f,direction:V,isReversed:x,vertical:E,readonly:L,elevation:I,mousePressed:N,decimals:M,indexFromEnd:F}=o,{textColorClasses:B,textColorStyles:_}=_e(c),{pageup:u,pagedown:w,end:v,home:m,left:y,right:S,down:r,up:P}=Fe,O=[u,w,v,m,y,S,r,P],R=s(()=>g.value?[1,2,3]:[1,5,10]);function $(b,p){if(!O.includes(b.key))return;b.preventDefault();const D=g.value||.1,q=(e.max-e.min)/D;if([y,S,r,P].includes(b.key)){const U=(E.value?[h.value?y:S,x.value?r:P]:F.value!==h.value?[y,P]:[S,P]).includes(b.key)?1:-1,X=b.shiftKey?2:b.ctrlKey?1:0;p=p+U*D*R.value[X]}else if(b.key===m)p=e.min;else if(b.key===v)p=e.max;else{const A=b.key===w?1:-1;p=p-A*D*(q>100?q/10:10)}return Math.max(e.min,Math.min(e.max,p))}function j(b){const p=$(b,e.modelValue);p!=null&&i("update:modelValue",p)}return W(()=>{const b=K(F.value?100-e.position:e.position,"%"),{elevationClasses:p}=Te(s(()=>C.value?void 0:I.value));return n("div",{class:["v-slider-thumb",{"v-slider-thumb--focused":e.focused,"v-slider-thumb--pressed":e.focused&&N.value},e.class,k.value],style:[{"--v-slider-thumb-position":b,"--v-slider-thumb-size":K(T.value)},e.style],role:"slider",tabindex:C.value?-1:0,"aria-valuemin":e.min,"aria-valuemax":e.max,"aria-valuenow":e.modelValue,"aria-readonly":!!L.value,"aria-orientation":V.value,onKeydown:L.value?void 0:j},[n("div",{class:["v-slider-thumb__surface",B.value,p.value],style:{..._.value}},null),ae(n("div",{class:["v-slider-thumb__ripple",B.value],style:_.value},null),[[ke("ripple"),e.ripple,null,{circle:!0,center:!0}]]),n(xe,{origin:"bottom center"},{default:()=>{var D;return[ae(n("div",{class:"v-slider-thumb__label-container"},[n("div",{class:["v-slider-thumb__label"]},[n("div",null,[((D=a["thumb-label"])==null?void 0:D.call(a,{modelValue:e.modelValue}))??e.modelValue.toFixed(g.value?M.value:1)])])]),[[ye,f.value&&e.focused||f.value==="always"]])]}})])}),{}}});const $e=H({start:{type:Number,required:!0},stop:{type:Number,required:!0},...ue()},"VSliderTrack"),Ae=Q()({name:"VSliderTrack",props:$e(),emits:{},setup(e,t){let{slots:a}=t;const i=ie(Z);if(!i)throw new Error("[Vuetify] v-slider-track must be inside v-slider or v-range-slider");const{color:o,parsedTicks:h,rounded:k,showTicks:c,tickSize:g,trackColor:C,trackFillColor:T,trackSize:f,vertical:V,min:x,max:E,indexFromEnd:L}=i,{roundedClasses:I}=Pe(k),{backgroundColorClasses:N,backgroundColorStyles:M}=re(T),{backgroundColorClasses:F,backgroundColorStyles:B}=re(C),_=s(()=>`inset-${V.value?"block":"inline"}-${L.value?"end":"start"}`),u=s(()=>V.value?"height":"width"),w=s(()=>({[_.value]:"0%",[u.value]:"100%"})),v=s(()=>e.stop-e.start),m=s(()=>({[_.value]:K(e.start,"%"),[u.value]:K(v.value,"%")})),y=s(()=>c.value?(V.value?h.value.slice().reverse():h.value).map((r,P)=>{var R;const O=r.value!==x.value&&r.value!==E.value?K(r.position,"%"):void 0;return n("div",{key:r.value,class:["v-slider-track__tick",{"v-slider-track__tick--filled":r.position>=e.start&&r.position<=e.stop,"v-slider-track__tick--first":r.value===x.value,"v-slider-track__tick--last":r.value===E.value}],style:{[_.value]:O}},[(r.label||a["tick-label"])&&n("div",{class:"v-slider-track__tick-label"},[((R=a["tick-label"])==null?void 0:R.call(a,{tick:r,index:P}))??r.label])])}):[]);return W(()=>n("div",{class:["v-slider-track",I.value,e.class],style:[{"--v-slider-track-size":K(f.value),"--v-slider-tick-size":K(g.value)},e.style]},[n("div",{class:["v-slider-track__background",F.value,{"v-slider-track__background--opacity":!!o.value||!T.value}],style:{...w.value,...B.value}},null),n("div",{class:["v-slider-track__fill",N.value],style:{...m.value,...M.value}},null),c.value&&n("div",{class:["v-slider-track__ticks",{"v-slider-track__ticks--always-show":c.value==="always"}]},[y.value])])),{}}}),je=H({...ze(),...De(),...Le(),modelValue:{type:[Number,String],default:0}},"VSlider"),Xe=Q()({name:"VSlider",props:je(),emits:{"update:focused":e=>!0,"update:modelValue":e=>!0,start:e=>!0,end:e=>!0},setup(e,t){let{slots:a,emit:i}=t;const o=G(),{rtlClasses:h}=J(),k=Ie(e),c=Me(e,"modelValue",void 0,u=>k.roundValue(u??k.min.value)),{min:g,max:C,mousePressed:T,roundValue:f,onSliderMousedown:V,onSliderTouchstart:x,trackContainerRef:E,position:L,hasLabels:I,readonly:N}=qe({props:e,steps:k,onSliderStart:()=>{i("start",c.value)},onSliderEnd:u=>{let{value:w}=u;const v=f(w);c.value=v,i("end",v)},onSliderMove:u=>{let{value:w}=u;return c.value=f(w)},getActiveThumb:()=>{var u;return(u=o.value)==null?void 0:u.$el}}),{isFocused:M,focus:F,blur:B}=Re(e),_=s(()=>L(c.value));return W(()=>{const u=ne.filterProps(e),w=!!(e.label||a.label||a.prepend);return n(ne,pe({class:["v-slider",{"v-slider--has-labels":!!a["tick-label"]||I.value,"v-slider--focused":M.value,"v-slider--pressed":T.value,"v-slider--disabled":e.disabled},h.value,e.class],style:e.style},u,{focused:M.value}),{...a,prepend:w?v=>{var m,y;return n(Se,null,[((m=a.label)==null?void 0:m.call(a,v))??(e.label?n(Ee,{id:v.id.value,class:"v-slider__label",text:e.label},null):void 0),(y=a.prepend)==null?void 0:y.call(a,v)])}:void 0,default:v=>{let{id:m,messagesId:y}=v;return n("div",{class:"v-slider__container",onMousedown:N.value?void 0:V,onTouchstartPassive:N.value?void 0:x},[n("input",{id:m.value,name:e.name||m.value,disabled:!!e.disabled,readonly:!!e.readonly,tabindex:"-1",value:c.value},null),n(Ae,{ref:E,start:0,stop:_.value},{"tick-label":a["tick-label"]}),n(Oe,{ref:o,"aria-describedby":y.value,focused:M.value,min:g.value,max:C.value,modelValue:c.value,"onUpdate:modelValue":S=>c.value=S,position:_.value,elevation:e.elevation,onFocus:F,onBlur:B,ripple:e.ripple},{"thumb-label":a["thumb-label"]})])}})}),{}}});export{Xe as V,qe as a,Ae as b,Oe as c,Ne as g,De as m,Ie as u};