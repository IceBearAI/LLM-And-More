import{_ as we}from"./NavBack.vue_vue_type_script_setup_true_lang-b6a7b049.js";import{ac as ce,i as Te,c as Ce,o as me,U as D,k as p,l as k,m as e,j as a,_ as i,$ as f,n as x,Y as U,W as s,aB as pe,al as ee,am as ae,d as H,au as te,x as F,r as I,aC as Z,av as E,K as De,Z as A,a5 as Me,F as Y,z as _e,a0 as se,aD as ne,D as re,a8 as de}from"./utils-5a8267ef.js";import{_ as ie}from"./UiParentCard.vue_vue_type_script_setup_true_lang-3c4d2264.js";import{_ as O,ae as G,e as $e,V as J,a5 as z,ai as Pe}from"./index-bdfe8eaa.js";import{V as M}from"./VCol-9b7c01c8.js";import{V as W}from"./VRow-522ff18f.js";import{C as Se}from"./CustomUpload-6dde583c.js";import{_ as Ue}from"./Explain.vue_vue_type_style_index_0_lang-5ba40fa4.js";import{V as Be}from"./VTextarea-e29a39b5.js";import{V as Le,a as Ne,b as Re,c as Ge}from"./VExpansionPanel-e3d4823f.js";import{V as Ee,a as ue}from"./VRadioGroup-d88f4c3c.js";import{V as Fe}from"./VForm-d4c8b802.js";import{i as Ae}from"./index-8cd50de6.js";import{V as X}from"./VAlert-6bf79ec1.js";import{_ as je}from"./ConfirmByInput.vue_vue_type_style_index_0_lang-0f40bd3c.js";import{_ as fe}from"./ConfirmByClick.vue_vue_type_style_index_0_lang-e4b43f61.js";import{_ as ze}from"./DialogLog.vue_vue_type_style_index_0_lang-05451b0f.js";import{I as qe}from"./IconTerminal2-210c96bf.js";import{V as Ye,a as Je,b as We,c as He}from"./VWindowItem-5afdd853.js";import"./IconInfoCircle-2525280f.js";import"./Confirm-f9ddfe55.js";import"./TextLog-eed03216.js";const j=g=>(ee("data-v-11edd06b"),g=g(),ae(),g),Oe=j(()=>s("label",null,"供应",-1)),Qe=j(()=>s("label",null,"类型",-1)),Ke=j(()=>s("label",null,"最长上下文",-1)),Ze=j(()=>s("label",null,"微调",-1)),Xe=j(()=>s("label",null,"参数量",-1)),ea={style:{color:"#539bff"}},aa=j(()=>s("label",null,"创建时间",-1)),ta=j(()=>s("label",null,"备注",-1)),la={__name:"ModelListDetailBaseInfo",setup(g){const{loadDictTree:S,getLabels:b}=ce(),o=Te("provideModelListDetail"),c=Ce(()=>o.rawData),n=async()=>{await S(["model_type"])};return me(()=>{n()}),(t,$)=>{const w=D("router-link"),C=D("el-tooltip");return p(),k(W,{class:"my-form waterfall"},{default:e(()=>[a(M,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[a(G,{"hide-details":""},{prepend:e(()=>[Oe]),default:e(()=>[i(" "+f(c.value.providerName),1)]),_:1})]),_:1}),a(M,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[a(G,{"hide-details":""},{prepend:e(()=>[Qe]),default:e(()=>[i(" "+f(x(b)([["model_type",c.value.modelType]])),1)]),_:1})]),_:1}),a(M,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[a(G,{"hide-details":""},{prepend:e(()=>[Ke]),default:e(()=>[i(" "+f(c.value.maxTokens),1)]),_:1})]),_:1}),c.value.isFineTuning?(p(),k(M,{key:0,xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[a(G,{"hide-details":""},{prepend:e(()=>[Ze]),default:e(()=>[a(C,{content:"微调详情",placement:"top"},{default:e(()=>[a(w,{to:"/model/fine-tuning/detail?jobId="+c.value.jobId,class:"link"},{default:e(()=>[i(f(c.value.isFineTuning?"是":""),1)]),_:1},8,["to"])]),_:1})]),_:1})]),_:1})):U("",!0),a(M,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[a(G,{"hide-details":""},{prepend:e(()=>[Xe]),default:e(()=>[s("span",ea,f(c.value.parameters)+"B",1)]),_:1})]),_:1}),a(M,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[a(G,{"hide-details":""},{prepend:e(()=>[aa]),default:e(()=>[i(" "+f(x(pe).dateFormat(c.value.createdAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1})]),_:1}),a(M,{xs:"12",sm:"6",md:"4",lg:"3"},{default:e(()=>[a(G,{"hide-details":""},{prepend:e(()=>[ta]),default:e(()=>[i(" "+f(c.value.remark),1)]),_:1})]),_:1})]),_:1})}}},oa=O(la,[["__scopeId","data-v-11edd06b"]]),L=g=>(ee("data-v-7c045d1f"),g=g(),ae(),g),sa={class:"mx-auto mt-3",style:{width:"540px"}},na=L(()=>s("label",{class:"required"},"评测指标",-1)),ra=L(()=>s("label",{class:"required"},"待评测数据集",-1)),da=L(()=>s("label",{class:"required"},"最大输出序列长度",-1)),ia=L(()=>s("label",{class:"required"},"单卡Batch大小",-1)),ua=L(()=>s("label",null,"备注",-1)),ca=L(()=>s("label",null,"调度标签",-1)),ma=L(()=>s("label",null,"k8s集群",-1)),pa=L(()=>s("label",null,"推理类型",-1)),_a=L(()=>s("label",{class:"required"},"CPU数量",-1)),fa=L(()=>s("label",{class:"required"},"GPU数量",-1)),va={class:"required"},ya=H({__name:"CreatePerformanceEvalPane",emits:["submit"],setup(g,{expose:S,emit:b}){const o=te(),c=F({operateType:"add"}),n=I(""),t=F({modelId:Number(o.query.jobId),fileId:"",evalTargetType:null,maxLength:512,batchSize:32,remark:"",label:null,k8sCluster:null,inferredType:"",cpu:0,gpu:0,maxGpuMemory:""}),$=b,w=I(),C=I(),v=F({fileId:[l=>!!l||"请选择微调文件"],evalTargetType:[l=>!!l||"请选择评测指标"],maxLength:[l=>!!l||"请输入最大输出序列长度"],batchSize:[l=>!!l||"请输入单卡Batch大小"],cpu:[l=>Z.validNumberInput(l,0,100,"请输入使用CPU数量",!0)],gpu:[l=>Z.validNumberInput(l,0,20,"请输入使用GPU数量",!0)],maxGpuMemory:[l=>Z.validNumberInput(l,1,80,"请输入GPU内存",!0)]}),T=({res:l})=>{l&&(t.fileId=l.fileId,n.value=l.filename)},u=()=>{n.value="",t.fileId=""},_=()=>{t.cpu!==1&&(t.cpu=1),t.gpu!==1&&(t.gpu=1)},y=async({valid:l,errors:r,showLoading:B})=>{if(l){const N={...t};t.inferredType==="cpu"?N.gpu=0:N.cpu=0,N.gpu<2&&delete N.maxGpuMemory;const[m,Q]=await E.post({showLoading:B,showSuccess:!0,url:"/evaluate/create",data:N});Q&&(w.value.hide(),$("submit"))}};return S({show({title:l,operateType:r}){w.value.show({title:l,refForm:C}),c.operateType=r,r=="add"&&(u(),t.evalTargetType=null,t.maxLength=512,t.batchSize=32,t.remark="",t.label=null,t.k8sCluster=null,t.inferredType="",t.cpu=0,t.gpu=0,t.maxGpuMemory="")}}),(l,r)=>{const B=D("Select"),N=D("Pane");return p(),k(N,{ref_key:"refPane",ref:w,onSubmit:y},{default:e(()=>[s("div",sa,[a(Fe,{ref_key:"refForm",ref:C,class:"my-form"},{default:e(()=>[a(B,{placeholder:"请选择评测指标",rules:v.evalTargetType,mapDictionary:{code:"model_evaluate_target_type"},modelValue:t.evalTargetType,"onUpdate:modelValue":r[0]||(r[0]=m=>t.evalTargetType=m)},{prepend:e(()=>[na]),_:1},8,["rules","modelValue"]),t.evalTargetType!=="five"?(p(),k(G,{key:0,rules:v.fileId,modelValue:t.fileId,"onUpdate:modelValue":r[1]||(r[1]=m=>t.fileId=m),"hide-details":"auto"},{prepend:e(()=>[ra]),default:e(()=>[n.value?(p(),k($e,{key:0,closable:"",color:"info","onClick:close":u},{default:e(()=>[i(f(n.value),1)]),_:1})):(p(),k(Se,{key:1,"file-type":[".txt",".json",".jsonl"],isSuffixValid:"",onAfterUpload:T},{trigger:e(()=>[a(J,{color:"info",variant:"outlined"},{default:e(()=>[i("上传文件")]),_:1})]),_:1},8,["file-type"]))]),_:1},8,["rules","modelValue"])):U("",!0),a(z,{type:"number",placeholder:"请输入最大输出序列长度","hide-details":"auto",rules:v.maxLength,modelValue:t.maxLength,"onUpdate:modelValue":r[2]||(r[2]=m=>t.maxLength=m),modelModifiers:{number:!0}},{prepend:e(()=>[da]),_:1},8,["rules","modelValue"]),a(z,{type:"number",placeholder:"请输入单卡Batch大小","hide-details":"auto",rules:v.batchSize,modelValue:t.batchSize,"onUpdate:modelValue":r[3]||(r[3]=m=>t.batchSize=m),modelModifiers:{number:!0}},{prepend:e(()=>[ia]),_:1},8,["rules","modelValue"]),a(Be,{modelValue:t.remark,"onUpdate:modelValue":r[4]||(r[4]=m=>t.remark=m),modelModifiers:{trim:!0},placeholder:"请输入你微调模型的备注",clearable:""},{prepend:e(()=>[ua]),_:1},8,["modelValue"]),a(G,null,{default:e(()=>[a(Le,null,{default:e(()=>[a(Ne,{elevation:"10"},{default:e(()=>[a(Re,{class:"text-h6"},{default:e(()=>[i("高级")]),_:1}),a(Ge,{class:"mt-4",eager:""},{default:e(()=>[a(B,{placeholder:"请选择调度标签","hide-details":!1,mapDictionary:{code:"model_deploy_label"},modelValue:t.label,"onUpdate:modelValue":r[5]||(r[5]=m=>t.label=m)},{prepend:e(()=>[ca]),_:1},8,["modelValue"]),a(B,{placeholder:"请选择k8s集群","hide-details":!1,mapDictionary:{code:"k8s_cluster"},modelValue:t.k8sCluster,"onUpdate:modelValue":r[6]||(r[6]=m=>t.k8sCluster=m)},{prepend:e(()=>[ma]),_:1},8,["modelValue"]),a(Ee,{modelValue:t.inferredType,"onUpdate:modelValue":[r[7]||(r[7]=m=>t.inferredType=m),_],inline:"",color:"primary"},{prepend:e(()=>[pa]),default:e(()=>[a(ue,{label:"CPU",value:"cpu"}),a(ue,{label:"GPU",value:"gpu"})]),_:1},8,["modelValue"]),t.inferredType==="cpu"?(p(),k(z,{key:0,type:"number",placeholder:"请输入使用CPU数量",rules:v.cpu,modelValue:t.cpu,"onUpdate:modelValue":r[8]||(r[8]=m=>t.cpu=m),modelModifiers:{number:!0}},{prepend:e(()=>[_a]),_:1},8,["rules","modelValue"])):U("",!0),t.inferredType==="gpu"?(p(),k(z,{key:1,type:"number",placeholder:"请输入使用GPU数量",rules:v.gpu,modelValue:t.gpu,"onUpdate:modelValue":r[9]||(r[9]=m=>t.gpu=m),modelModifiers:{number:!0}},{prepend:e(()=>[fa]),_:1},8,["rules","modelValue"])):U("",!0),t.gpu>1?(p(),k(z,{key:2,type:"number",placeholder:"请输入GPU内存",rules:v.maxGpuMemory,modelValue:t.maxGpuMemory,"onUpdate:modelValue":r[10]||(r[10]=m=>t.maxGpuMemory=m),modelModifiers:{number:!0},max:"80","hide-details":"auto"},{prepend:e(()=>[s("label",va,[i("GPU内存 "),a(Ue,null,{default:e(()=>[i("指定每个 GPU 用于存储模型权重的最大内存。这允许它为激活分配更多内存，因此您可以使用更长的上下文长度或更大的批量大小")]),_:1})])]),"append-inner":e(()=>[i(" GiB ")]),_:1},8,["rules","modelValue"])):U("",!0)]),_:1})]),_:1})]),_:1})]),_:1})]),_:1},512)])]),_:1},512)}}});const ha=O(ya,[["__scopeId","data-v-7c045d1f"]]);const ba={__name:"GradeRadarChart",setup(g,{expose:S}){const b=I();var o=null;const c=({title:n,radar:t,seriesData:$})=>{o=Ae(b.value);const w=$.map(C=>C.name);o.setOption({title:{text:n},legend:{icon:"circle",data:w},radar:t,series:[{type:"radar",data:$}]})};return De(()=>{o&&(o.dispose(),o=null)}),Pe(b,n=>{o==null||o.resize()}),S({initChart:c}),(n,t)=>(p(),A("div",{ref_key:"refChart",ref:b,class:"chart-item w-100 h-100"},null,512))}},ga=O(ba,[["__scopeId","data-v-8ceb45d1"]]),xa=g=>(ee("data-v-6813598c"),g=g(),ae(),g),ka={style:{height:"500px"}},Va=xa(()=>s("h5",{class:"text-h6 text-capitalize"},"建议进一步操作",-1)),Ia={class:"text-center mt-1"},wa=H({__name:"ViewEvalDetailPane",setup(g,{expose:S}){const b=I(),o=I(),c=F({currentModelId:null,compare1ModelId:null,compare2ModelId:null}),n=I({}),t=[{key:"current",value:"currentModelId",disabled:!0},{key:"compare1",value:"compare1ModelId",disabled:!1},{key:"compare2",value:"compare2ModelId",disabled:!1}],$=F({id:""}),w=()=>{const v={indicator:[{name:"中文能力",max:10,axisLabel:{show:!0}},{name:"推理能力",max:10},{name:"指令遵从能力",max:10},{name:"创新能力",max:10},{name:"阅读理解",max:10}]},T=[];n.value&&t.forEach(u=>{n.value[u.key].value&&T.push({name:n.value[u.key].name,value:n.value[u.key].value})}),o.value.initChart({title:"",radar:v,seriesData:T})},C=async(v=!1)=>{if(v)return;const[T,u]=await E.post({url:"/evaluate/fivegraph",showLoading:b.value.el,data:{currentModelId:c.currentModelId,compare1ModelId:c.compare1ModelId,compare2ModelId:c.compare2ModelId,currentModelEvaluateId:$.id}});u&&(n.value=u,w())};return S({show({title:v,modelId:T,id:u}){b.value.show({width:900,title:v,hasSubmitBtn:!1}),c.currentModelId=parseInt(T),c.compare1ModelId=null,c.compare2ModelId=null,$.id=u,C()}}),(v,T)=>{const u=D("Select"),_=D("Pane");return p(),k(_,{ref_key:"refPane",ref:b},{default:e(()=>[a(W,null,{default:e(()=>{var y;return[a(M,{cols:(y=n.value.current)!=null&&y.isFineTuning?7:12},{default:e(()=>[s("div",ka,[a(ga,{ref_key:"gradeRadarChartRef",ref:o},null,512)])]),_:1},8,["cols"]),n.value.current&&n.value.current.isFineTuning?(p(),k(M,{key:0,cols:"5"},{default:e(()=>[a(X,{type:n.value.current.riskOver?"error":"success",variant:"tonal"},{default:e(()=>[i("过拟合风险")]),_:1},8,["type"]),a(X,{type:n.value.current.riskUnder?"error":"success",variant:"tonal",class:"mt-4"},{default:e(()=>[i("欠拟合风险")]),_:1},8,["type"]),a(X,{type:"info",variant:"tonal",class:"mt-4"},{default:e(()=>[Va,s("div",null,f(n.value.current.remind),1)]),_:1})]),_:1})):U("",!0),a(M,{cols:"12"},{default:e(()=>[a(W,null,{default:e(()=>[(p(),A(Y,null,Me(t,l=>a(M,{cols:"4"},{default:e(()=>{var r;return[a(u,{placeholder:"请选择模型",modelValue:c[l.value],"onUpdate:modelValue":B=>c[l.value]=B,disabled:l.disabled,clearable:!1,mapAPI:{url:"/channels/models",data:{pageSize:-1,providerName:"LocalAI",modelType:"text-generation",evalTag:"five"},labelField:"modelName",valueField:"id"},onChange:B=>C(l.disabled)},null,8,["modelValue","onUpdate:modelValue","disabled","onChange"]),s("div",Ia,f((r=n.value[l.key])==null?void 0:r.score),1)]}),_:2},1024)),64))]),_:1})]),_:1})]}),_:1})]),_:1},512)}}});const Ta=O(wa,[["__scopeId","data-v-6813598c"]]),Ca={class:""},Da={class:"d-flex"},Ma={style:{width:"250px"},class:"mr-4"},$a={style:{width:"250px"}},Pa=s("span",null,"请先部署模型 ",-1),Sa={key:1,class:"d-flex align-center justify-center"},Ua=["onClick"],Ba=["onClick"],La=s("span",{class:"text-primary"},"取消",-1),Na=s("br",null,null,-1),Ra={class:"text-primary font-weight-black"},Ga=s("br",null,null,-1),Ea=H({__name:"TabModelEstimate",props:{showArrange:String,modelTitle:String,providerName:String},setup(g){const S=te(),{modelName:b,jobId:o}=S.query,{loadDictTree:c,getLabels:n}=ce(),t=I(),$=I(),w=I(),C=I(),v=I(),T=I(),u=F({style:{},formData:{modelName:"",evalTargetType:""},showTooltip:!1,timer:null,tableInfos:{list:[],total:0},currentJobId:""}),_=F({status:null,evalTargetType:null}),y=F({uuid:null}),l=I();_e(u);const r=()=>{$.value.show({title:"性能评估",operateType:"add"})},B=h=>{w.value.show({title:"详情",id:h.id,modelId:o})},N=h=>{u.currentJobId=h.uuid,v.value.show({width:"450px",confirmText:u.currentJobId})},m=h=>{y.uuid=h.uuid,T.value.show({width:"450px"})},Q=async(h={})=>{const[V,P]=await E.delete({...h,showSuccess:!0,url:`/evaluate/delete/${y.uuid}`});P&&(T.value.hide(),oe())},ve=async(h={})=>{const[V,P]=await E.put({...h,showSuccess:!0,url:`/evaluate/cancel/${u.currentJobId}`});P&&(v.value.hide(),oe())},le=async(h={})=>{await c(["model_evaluate_data_type"]);const[V,P]=await E.get({url:"/evaluate/list",showLoading:t.value.el,data:{modelId:o,..._,...h}});P?(u.tableInfos.total=P.total,u.tableInfos.list=P.list||[]):(u.tableInfos.list=[],u.tableInfos.total=0)},K=()=>{t.value.query({page:1})},oe=()=>{t.value.query()},ye=h=>{let V=[];return h.status=="waiting"||h.status=="running"?V.push({text:"取消",color:"",click(){N(h)}}):V.push({text:"删除",color:"error",click(){m(h)}}),V},he=async h=>{u.currentJobId=h.uuid,l.value.show()},be=async()=>{let[h,V]=await E.get({url:`/api/evaluate/${b}/eval-log/${u.currentJobId}`});if(V){const P=Object.keys(V).length===0?"":V;l.value.setContent(P)}};return me(()=>{le()}),(h,V)=>{const P=D("Select"),q=D("el-tooltip"),ge=D("ButtonsInForm"),R=D("el-table-column"),xe=D("router-link"),ke=D("ButtonsInTable"),Ve=D("TableWithPager");return p(),A(Y,null,[s("div",Ca,[a(W,null,{default:e(()=>[a(M,{cols:"12",class:"d-flex justify-space-between align-center",style:{height:"76px"}},{default:e(()=>[s("div",Da,[s("div",Ma,[a(P,{onChange:K,label:"请选择评估状态",mapDictionary:{code:"model_eval_status"},modelValue:_.status,"onUpdate:modelValue":V[0]||(V[0]=d=>_.status=d)},null,8,["modelValue"])]),s("div",$a,[a(P,{onChange:K,label:"请选择指标",mapDictionary:{code:"model_evaluate_target_type"},modelValue:_.evalTargetType,"onUpdate:modelValue":V[1]||(V[1]=d=>_.evalTargetType=d)},null,8,["modelValue"])])]),a(q,{ref:"tooltipRef",visible:u.showTooltip,"popper-options":{modifiers:[{name:"computeStyles",options:{adaptive:!1,enabled:!1}}]},"auto-close":1,"virtual-ref":C.value,"virtual-triggering":"","popper-class":"singleton-tooltip"},{content:e(()=>[Pa]),_:1},8,["visible","virtual-ref"]),a(ge,null,{default:e(()=>[g.providerName==="LocalAI"?(p(),k(J,{key:0,color:"primary",onClick:r},{default:e(()=>[i("创建评估")]),_:1})):U("",!0)]),_:1})]),_:1}),a(M,{cols:"12"},{default:e(()=>[a(Ve,{onQuery:le,ref_key:"refTableWithPager",ref:t,infos:u.tableInfos},{default:e(()=>[a(R,{label:"指标","min-width":"100px"},{default:e(({row:d})=>[i(f(x(n)([["model_evaluate_target_type",d.evalTargetType]])),1)]),_:1}),a(R,{label:"评估状态","min-width":"120px"},{default:e(({row:d})=>[d.status==="failed"?(p(),k(q,{key:0,content:d.statusMsg,placement:"top","raw-content":""},{default:e(()=>[s("span",{class:se(`text-${x(ne).statusMap[d.status].color}`)},f(x(n)([["model_eval_status",d.status]])),3)]),_:2},1032,["content"])):(p(),A("div",Sa,[s("span",{class:se(`text-${x(ne).statusMap[d.status].color}`)},f(x(n)([["model_eval_status",d.status]])),3),["running","success"].includes(d.status)?(p(),A("div",{key:0,class:"link ml-1",onClick:Ie=>he(d)},"(日志)",8,Ua)):U("",!0)]))]),_:1}),a(R,{label:"webshell","min-width":"90px"},{default:e(({row:d})=>[d.status==="running"?(p(),k(q,{key:0,content:"进入终端",placement:"top"},{default:e(()=>[a(xe,{class:"link",to:{path:"/model/terminal",query:{resourceType:"eval-job",serviceName:d.jobName}},target:"_blank"},{default:e(()=>[a(x(qe),{class:"align-top",size:20})]),_:2},1032,["to"])]),_:2},1024)):U("",!0)]),_:1}),a(R,{label:"数据量","min-width":"100px"},{default:e(({row:d})=>[i(f(d.dataSize),1)]),_:1}),a(R,{label:"评估数据集","min-width":"120px"},{default:e(({row:d})=>[i(f(x(n)([["model_evaluate_data_type",d.dataType]])),1)]),_:1}),a(R,{label:"平均分","min-width":"100px"},{default:e(({row:d})=>[d.evalTargetType==="five"&&d.score>0?(p(),k(q,{key:0,content:"点击可查看五维图指标详情",placement:"top"},{default:e(()=>[s("span",{class:"link",onClick:Ie=>B(d)},f(d.score),9,Ba)]),_:2},1024)):(p(),A(Y,{key:1},[i(f(d.score),1)],64))]),_:1}),a(R,{label:"备注","min-width":"200px","show-overflow-tooltip":""},{default:e(({row:d})=>[i(f(d.remark),1)]),_:1}),a(R,{label:"创建时间","min-width":"160px"},{default:e(({row:d})=>[i(f(x(pe).dateFormat(d.createdAt,"YYYY-MM-DD HH:mm:ss")),1)]),_:1}),a(R,{label:"操作",width:"80px",fixed:"right"},{default:e(({row:d})=>[a(ke,{buttons:ye(d),onlyOne:""},null,8,["buttons"])]),_:1})]),_:1},8,["infos"])]),_:1})]),_:1})]),a(ze,{ref_key:"refDialogLog",ref:l,interval:20,onRefresh:be},null,512),a(je,{ref_key:"refConfirmAbort",ref:v,onSubmit:ve},{text:e(()=>[i(" 此操作将会"),La,i("正在进行的模型评估"),Na,i(" 任务ID："),s("span",Ra,f(u.currentJobId),1),Ga,i(" 你还要继续吗？ ")]),_:1},512),a(fe,{ref_key:"refConfirmDelete",ref:T,onSubmit:Q},{text:e(()=>[i("确定要删除该评估数据？")]),_:1},512),a(ha,{ref_key:"createPerformanceEvalPaneRef",ref:$,onSubmit:K},null,512),a(Ta,{ref_key:"viewEvalDetailPaneRef",ref:w},null,512)],64)}}});const Fa=["innerHTML"],it=H({__name:"modelListDetail",setup(g){const S=te(),b=I(),o=F({tabIndex:"2",style:{},formData:{},rawData:{modelName:"",deployStatus:"",providerName:""},confirmByClickInfo:{html:"",action:"",row:{}}}),{style:c,rawData:n,formData:t}=_e(o);re("provideModelListDetail",o),re("provideModelListDetailed",o),(async()=>{let{jobId:_}=S.query,[y,l]=await E.get({showLoading:!0,url:`/api/models/${_}`});l&&(o.rawData=l)})();const w=(_={})=>{let{action:y,row:l}=o.confirmByClickInfo;y=="deploy"?C(l,_):y=="undeploy"&&v(l,_)},C=async(_,y)=>{let[l,r]=await E.post({...y,showSuccess:!0,url:`/api/models/${_.id}/deploy`,data:{...o.formData,id:_.id}});r&&b.value.hide()},v=async(_,y)=>{let[l,r]=await E.post({...y,showSuccess:!0,url:`/api/models/${_.id}/undeploy`,data:{...o.formData,id:_.id}});r&&b.value.hide()},T=()=>{de.assign(o.confirmByClickInfo,{html:"确认部署模型 <span class='text-primary mx-1  font-weight-black'>"+o.rawData.modelName+"</span>吗？",action:"deploy",row:o.rawData}),b.value.show({width:"350px"})},u=()=>{de.assign(o.confirmByClickInfo,{html:"确认卸载模型 <span class='text-primary mx-1  font-weight-black'>"+o.rawData.modelName+"</span>吗？",action:"undeploy",row:o.rawData}),b.value.show({width:"350px"})};return(_,y)=>(p(),A(Y,null,[a(we,{backUrl:"/model/model-list/list"},{default:e(()=>[i("模型详情")]),_:1}),a(ie,{class:"mt-4"},{header:e(()=>[i(" 模型："+f(x(n).modelName),1)]),action:e(()=>[x(n).deployStatus=="success"||x(n).deployStatus=="running"?(p(),k(J,{key:0,flat:"",color:"secondary",onClick:u},{default:e(()=>[i("卸载")]),_:1})):U("",!0),x(n).deployStatus=="failed"?(p(),k(J,{key:1,flat:"",color:"secondary",onClick:T},{default:e(()=>[i("部署")]),_:1})):U("",!0)]),default:e(()=>[a(oa)]),_:1}),a(ie,{class:"mt-5"},{header:e(()=>[a(Ye,{modelValue:o.tabIndex,"onUpdate:modelValue":y[0]||(y[0]=l=>o.tabIndex=l),"align-tabs":"start",color:"primary"},{default:e(()=>[a(Je,{value:2},{default:e(()=>[i("当前模型下的评估列表")]),_:1})]),_:1},8,["modelValue"])]),default:e(()=>[a(We,{modelValue:o.tabIndex,"onUpdate:modelValue":y[1]||(y[1]=l=>o.tabIndex=l)},{default:e(()=>[a(He,{value:2},{default:e(()=>[a(Ea,{showArrange:x(n).deployStatus,modelTitle:x(n).modelName,providerName:x(n).providerName},null,8,["showArrange","modelTitle","providerName"])]),_:1})]),_:1},8,["modelValue"])]),_:1}),a(fe,{ref_key:"refConfirmByClick",ref:b,onSubmit:w},{text:e(()=>[s("div",{innerHTML:o.confirmByClickInfo.html},null,8,Fa)]),_:1},512)],64))}});export{it as default};
