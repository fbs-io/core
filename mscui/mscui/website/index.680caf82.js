import{_ as K,ac as B,E as O,ad as z,b as j,ae as q,d as U,af as M,N as G,O as H,ab as J,a2 as R,k as $,V as W,o as p,c as Q,f as i,w as l,W as X,e as f,g as _,t as v,ag as D,C,J as Y,l as m,m as Z,q as ee,a4 as te,ah as ie,a5 as se}from"./index.24f9f7c9.js";/* empty css                      *//* empty css                  *//* empty css                *//* empty css                *//* empty css                        */import le from"./dic.440d6c23.js";import ae from"./list.22331293.js";const ne={name:"dic",components:{dicDialog:le,listDialog:ae},data(){return{dialog:{dic:!1,info:!1},showDicloading:!0,dicList:[],dicFilterText:"",dicProps:{label:"name"},listApi:null,listApiParams:{},selection:[]}},watch:{dicFilterText(e){this.$refs.dic.filter(e)}},mounted(){this.getDic(),this.rowDrop()},methods:{async getDic(){var e=await this.$API.system.dic.tree.get();this.showDicloading=!1,this.dicList=e.data;var t=this.dicList[0];t&&(this.$nextTick(()=>{this.$refs.dic.setCurrentKey(t.id)}),this.listApiParams={code:t.code},this.listApi=this.$API.system.dic.list)},dicFilterNode(e,t){if(!e)return!0;var s=t.name+t.code;return s.indexOf(e)!==-1},addDic(){this.dialog.dic=!0,this.$nextTick(()=>{this.$refs.dicDialog.open()})},dicEdit(e){this.dialog.dic=!0,this.$nextTick(()=>{var t=this.$refs.dic.getNode(e.id),s=t.level==1?void 0:t.parent.data.id;e.parentId=s,this.$refs.dicDialog.open("edit").setData(e)})},dicClick(e){this.$refs.table.reload({code:e.code})},dicDel(e,t){this.$confirm(`\u786E\u5B9A\u5220\u9664 ${t.name} \u9879\u5417\uFF1F`,"\u63D0\u793A",{type:"warning"}).then(()=>{this.showDicloading=!0;var s=this.$refs.dic.getCurrentKey();if(this.$refs.dic.remove(t.id),s==t.id){var c=this.dicList[0];c?(this.$refs.dic.setCurrentKey(c.id),this.$refs.table.upData({code:c.code})):(this.listApi=null,this.$refs.table.tableData=[])}this.showDicloading=!1,this.$message.success("\u64CD\u4F5C\u6210\u529F")}).catch(()=>{})},rowDrop(){const e=this,t=this.$refs.table.$el.querySelector(".el-table__body-wrapper tbody");B.create(t,{handle:".move",animation:300,ghostClass:"ghost",onEnd({newIndex:s,oldIndex:c}){const a=e.$refs.table.tableData,o=a.splice(c,1)[0];a.splice(s,0,o),e.$message.success("\u6392\u5E8F\u6210\u529F")}})},addInfo(){this.dialog.list=!0,this.$nextTick(()=>{var e=this.$refs.dic.getCurrentKey();const t={dic:e};this.$refs.listDialog.open().setData(t)})},table_edit(e){this.dialog.list=!0,this.$nextTick(()=>{this.$refs.listDialog.open("edit").setData(e)})},async table_del(e,t){var s={id:e.id},c=await this.$API.demo.post.post(s);c.code==200?(this.$refs.table.tableData.splice(t,1),this.$message.success("\u5220\u9664\u6210\u529F")):this.$alert(c.message,"\u63D0\u793A",{type:"error"})},async batch_del(){this.$confirm(`\u786E\u5B9A\u5220\u9664\u9009\u4E2D\u7684 ${this.selection.length} \u9879\u5417\uFF1F`,"\u63D0\u793A",{type:"warning"}).then(()=>{const e=this.$loading();this.selection.forEach(t=>{this.$refs.table.tableData.forEach((s,c)=>{t.id===s.id&&this.$refs.table.tableData.splice(c,1)})}),e.close(),this.$message.success("\u64CD\u4F5C\u6210\u529F")}).catch(()=>{})},saveList(){this.$refs.listDialog.submit(async e=>{this.isListSaveing=!0;var t=await this.$API.demo.post.post(e);this.isListSaveing=!1,t.code==200?(this.listDialogVisible=!1,this.$message.success("\u64CD\u4F5C\u6210\u529F")):this.$alert(t.message,"\u63D0\u793A",{type:"error"})})},selectionChange(e){this.selection=e},changeSwitch(e,t){t.yx=t.yx=="1"?"0":"1",t.$switch_yx=!0,setTimeout(()=>{delete t.$switch_yx,t.yx=e,this.$message.success(`\u64CD\u4F5C\u6210\u529Fid:${t.id} val:${e}`)},500)},handleDicSuccess(e,t){if(t=="add")e.id=new Date().getTime(),this.dicList.length>0?this.$refs.table.upData({code:e.code}):(this.listApiParams={code:e.code},this.listApi=this.$API.dic.info),this.$refs.dic.append(e,e.parentId[0]),this.$refs.dic.setCurrentKey(e.id);else if(t=="edit"){var s=this.$refs.dic.getNode(e.id),c=s.level==1?void 0:s.parent.data.id;if(c!=e.parentId){var a=s.data;this.$refs.dic.remove(e.id),this.$refs.dic.append(a,e.parentId[0])}Object.assign(s.data,e)}},handleListSuccess(e,t){t=="add"?(e.id=new Date().getTime(),this.$refs.table.tableData.push(e)):t=="edit"&&this.$refs.table.tableData.filter(s=>s.id===e.id).forEach(s=>{Object.assign(s,e)})}}},oe={class:"custom-tree-node"},ce={class:"label"},de={class:"code"},re={class:"do"},he=m("\u5B57\u5178\u5206\u7C7B"),_e={class:"left-panel"},pe=m("\u7F16\u8F91"),ue=m("\u5220\u9664");function fe(e,t,s,c,a,o){const w=Z,g=O,r=ee,b=te,x=z,y=j,k=q,u=U,E=M,h=G,S=ie,T=H,A=J,I=R,L=se,P=$("dic-dialog"),V=$("list-dialog"),N=W;return p(),Q(Y,null,[i(u,null,{default:l(()=>[X((p(),f(E,{width:"300px"},{default:l(()=>[i(u,null,{default:l(()=>[i(g,null,{default:l(()=>[i(w,{placeholder:"\u8F93\u5165\u5173\u952E\u5B57\u8FDB\u884C\u8FC7\u6EE4",modelValue:a.dicFilterText,"onUpdate:modelValue":t[0]||(t[0]=n=>a.dicFilterText=n),clearable:""},null,8,["modelValue"])]),_:1}),i(y,{class:"nopadding"},{default:l(()=>[i(x,{ref:"dic",class:"menu","node-key":"id",data:a.dicList,props:a.dicProps,"highlight-current":!0,"expand-on-click-node":!1,"filter-node-method":o.dicFilterNode,onNodeClick:o.dicClick},{default:l(({node:n,data:d})=>[_("span",oe,[_("span",ce,v(n.label),1),_("span",de,v(d.code),1),_("span",re,[i(b,null,{default:l(()=>[i(r,{icon:"el-icon-edit",size:"small",onClick:D(F=>o.dicEdit(d),["stop"])},null,8,["onClick"]),i(r,{icon:"el-icon-delete",size:"small",onClick:D(F=>o.dicDel(n,d),["stop"])},null,8,["onClick"])]),_:2},1024)])])]),_:1},8,["data","props","filter-node-method","onNodeClick"])]),_:1}),i(k,{style:{height:"51px"}},{default:l(()=>[i(r,{type:"primary",size:"small",icon:"el-icon-plus",style:{width:"100%"},onClick:o.addDic},{default:l(()=>[he]),_:1},8,["onClick"])]),_:1})]),_:1})]),_:1})),[[N,a.showDicloading]]),i(u,{class:"is-vertical"},{default:l(()=>[i(g,null,{default:l(()=>[_("div",_e,[i(r,{type:"primary",icon:"el-icon-plus",onClick:o.addInfo},null,8,["onClick"]),i(r,{type:"danger",plain:"",icon:"el-icon-delete",disabled:a.selection.length==0,onClick:o.batch_del},null,8,["disabled","onClick"])])]),_:1}),i(y,{class:"nopadding"},{default:l(()=>[i(L,{ref:"table",apiObj:a.listApi,"row-key":"id",params:a.listApiParams,onSelectionChange:o.selectionChange,stripe:"",paginationLayout:"prev, pager, next"},{default:l(()=>[i(h,{type:"selection",width:"50"}),i(h,{label:"",width:"60"},{default:l(()=>[i(T,{class:"move",style:{cursor:"move"}},{default:l(()=>[i(S,{style:{width:"1em",height:"1em"}})]),_:1})]),_:1}),i(h,{label:"\u540D\u79F0",prop:"name",width:"150"}),i(h,{label:"\u952E\u503C",prop:"key",width:"150"}),i(h,{label:"\u662F\u5426\u6709\u6548",prop:"yx",width:"100"},{default:l(n=>[i(A,{modelValue:n.row.yx,"onUpdate:modelValue":d=>n.row.yx=d,onChange:d=>o.changeSwitch(d,n.row),loading:n.row.$switch_yx,"active-value":"1","inactive-value":"0"},null,8,["modelValue","onUpdate:modelValue","onChange","loading"])]),_:1}),i(h,{label:"\u64CD\u4F5C",fixed:"right",align:"right",width:"120"},{default:l(n=>[i(b,null,{default:l(()=>[i(r,{text:"",type:"primary",size:"small",onClick:d=>o.table_edit(n.row,n.$index)},{default:l(()=>[pe]),_:2},1032,["onClick"]),i(I,{title:"\u786E\u5B9A\u5220\u9664\u5417\uFF1F",onConfirm:d=>o.table_del(n.row,n.$index)},{reference:l(()=>[i(r,{text:"",type:"primary",size:"small"},{default:l(()=>[ue]),_:1})]),_:2},1032,["onConfirm"])]),_:2},1024)]),_:1})]),_:1},8,["apiObj","params","onSelectionChange"])]),_:1})]),_:1})]),_:1}),a.dialog.dic?(p(),f(P,{key:0,ref:"dicDialog",onSuccess:o.handleDicSuccess,onClosed:t[1]||(t[1]=n=>a.dialog.dic=!1)},null,8,["onSuccess"])):C("",!0),a.dialog.list?(p(),f(V,{key:1,ref:"listDialog",onSuccess:o.handleListSuccess,onClosed:t[2]||(t[2]=n=>a.dialog.list=!1)},null,8,["onSuccess"])):C("",!0)],64)}const we=K(ne,[["render",fe],["__scopeId","data-v-e431ef84"]]);export{we as default};