import{_ as E,a as N,E as k,$ as q,aq as F,ar as U,b as P,Z as S,o as u,c as d,w as t,d as r,P as n,l as f,F as g,n as x,e as _,ax as C,f as D}from"./index.d39b14b4.js";const I={emits:["success","closed"],data(){return{mode:"add",titleMap:{add:"\u65B0\u589E\u7528\u6237",edit:"\u7F16\u8F91\u7528\u6237",show:"\u67E5\u770B"},visible:!1,isSaveing:!1,form:{id:"",userName:"",avatar:"",name:"",dept:"",group:[]},rules:{avatar:[{required:!0,message:"\u8BF7\u4E0A\u4F20\u5934\u50CF"}],userName:[{required:!0,message:"\u8BF7\u8F93\u5165\u767B\u5F55\u8D26\u53F7"}],name:[{required:!0,message:"\u8BF7\u8F93\u5165\u771F\u5B9E\u59D3\u540D"}],password:[{required:!0,message:"\u8BF7\u8F93\u5165\u767B\u5F55\u5BC6\u7801"},{validator:(l,e,i)=>{this.form.password2!==""&&this.$refs.dialogForm.validateField("password2"),i()}}],password2:[{required:!0,message:"\u8BF7\u518D\u6B21\u8F93\u5165\u5BC6\u7801"},{validator:(l,e,i)=>{e!==this.form.password?i(new Error("\u4E24\u6B21\u8F93\u5165\u5BC6\u7801\u4E0D\u4E00\u81F4!")):i()}}],dept:[{required:!0,message:"\u8BF7\u9009\u62E9\u6240\u5C5E\u90E8\u95E8"}],group:[{required:!0,message:"\u8BF7\u9009\u62E9\u6240\u5C5E\u89D2\u8272",trigger:"change"}]},groups:[],groupsProps:{value:"id",multiple:!0,checkStrictly:!0},depts:[],deptsProps:{value:"id",checkStrictly:!0}}},mounted(){this.getGroup(),this.getDept()},methods:{open(l="add"){return this.mode=l,this.visible=!0,this},async getGroup(){var l=await this.$API.system.role.list.get();this.groups=l.data.rows},async getDept(){var l=await this.$API.system.dept.list.get();this.depts=l.data},submit(){this.$refs.dialogForm.validate(async l=>{if(l){this.isSaveing=!0;var e=await this.$API.demo.post.post(this.form);this.isSaveing=!1,e.code==200?(this.$emit("success",this.form,this.mode),this.visible=!1,this.$message.success("\u64CD\u4F5C\u6210\u529F")):this.$alert(e.message,"\u63D0\u793A",{type:"error"})}else return!1})},setData(l){this.form.id=l.id,this.form.userName=l.userName,this.form.avatar=l.avatar,this.form.name=l.name,this.form.group=l.group,this.form.dept=l.dept}}},B=_("\u53D6 \u6D88"),A=_("\u4FDD \u5B58");function G(l,e,i,M,o,v){const w=C,a=N,m=k,b=q,c=F,V=U,h=P,p=D,y=S;return u(),d(y,{title:o.titleMap[o.mode],modelValue:o.visible,"onUpdate:modelValue":e[9]||(e[9]=s=>o.visible=s),width:500,"destroy-on-close":"",onClosed:e[10]||(e[10]=s=>l.$emit("closed"))},{footer:t(()=>[r(p,{onClick:e[7]||(e[7]=s=>o.visible=!1)},{default:t(()=>[B]),_:1}),o.mode!="show"?(u(),d(p,{key:0,type:"primary",loading:o.isSaveing,onClick:e[8]||(e[8]=s=>v.submit())},{default:t(()=>[A]),_:1},8,["loading"])):n("",!0)]),default:t(()=>[r(h,{model:o.form,rules:o.rules,disabled:o.mode=="show",ref:"dialogForm","label-width":"100px","label-position":"left"},{default:t(()=>[r(a,{label:"\u5934\u50CF",prop:"avatar"},{default:t(()=>[r(w,{modelValue:o.form.avatar,"onUpdate:modelValue":e[0]||(e[0]=s=>o.form.avatar=s),title:"\u4E0A\u4F20\u5934\u50CF"},null,8,["modelValue"])]),_:1}),r(a,{label:"\u767B\u5F55\u8D26\u53F7",prop:"userName"},{default:t(()=>[r(m,{modelValue:o.form.userName,"onUpdate:modelValue":e[1]||(e[1]=s=>o.form.userName=s),placeholder:"\u7528\u4E8E\u767B\u5F55\u7CFB\u7EDF",clearable:""},null,8,["modelValue"])]),_:1}),r(a,{label:"\u59D3\u540D",prop:"name"},{default:t(()=>[r(m,{modelValue:o.form.name,"onUpdate:modelValue":e[2]||(e[2]=s=>o.form.name=s),placeholder:"\u8BF7\u8F93\u5165\u5B8C\u6574\u7684\u771F\u5B9E\u59D3\u540D",clearable:""},null,8,["modelValue"])]),_:1}),o.mode=="add"?(u(),f(g,{key:0},[r(a,{label:"\u767B\u5F55\u5BC6\u7801",prop:"password"},{default:t(()=>[r(m,{type:"password",modelValue:o.form.password,"onUpdate:modelValue":e[3]||(e[3]=s=>o.form.password=s),clearable:"","show-password":""},null,8,["modelValue"])]),_:1}),r(a,{label:"\u786E\u8BA4\u5BC6\u7801",prop:"password2"},{default:t(()=>[r(m,{type:"password",modelValue:o.form.password2,"onUpdate:modelValue":e[4]||(e[4]=s=>o.form.password2=s),clearable:"","show-password":""},null,8,["modelValue"])]),_:1})],64)):n("",!0),r(a,{label:"\u6240\u5C5E\u90E8\u95E8",prop:"dept"},{default:t(()=>[r(b,{modelValue:o.form.dept,"onUpdate:modelValue":e[5]||(e[5]=s=>o.form.dept=s),options:o.depts,props:o.deptsProps,clearable:"",style:{width:"100%"}},null,8,["modelValue","options","props"])]),_:1}),r(a,{label:"\u6240\u5C5E\u89D2\u8272",prop:"group"},{default:t(()=>[r(V,{modelValue:o.form.group,"onUpdate:modelValue":e[6]||(e[6]=s=>o.form.group=s),multiple:"",filterable:"",style:{width:"100%"}},{default:t(()=>[(u(!0),f(g,null,x(o.groups,s=>(u(),d(c,{key:s.id,label:s.label,value:s.id},null,8,["label","value"]))),128))]),_:1},8,["modelValue"])]),_:1})]),_:1},8,["model","rules","disabled"])]),_:1},8,["title","modelValue"])}const O=E(I,[["render",G]]);export{O as default};