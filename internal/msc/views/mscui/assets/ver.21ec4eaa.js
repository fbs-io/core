import{_,C as l,o as p,c as g,w as o,m as t,t as n,d as r,e as a,f as h}from"./index.d39b14b4.js";/* empty css                */const m="/mscui/img/ver.svg",u={title:"\u7248\u672C\u4FE1\u606F",icon:"el-icon-monitor",description:"\u5F53\u524D\u9879\u76EE\u7248\u672C\u4FE1\u606F",data(){return{ver:"loading..."}},mounted(){this.getVer()},methods:{async getVer(){const e=await this.$API.demo.ver.get();this.ver=e.data},golog(){window.open("https://gitee.com/qbiancheng/scui-vite-ui/releases")},gogit(){window.open("https://gitee.com/qbiancheng/scui-vite-ui")}}},v={style:{height:"210px","text-align":"center"}},y=t("img",{src:m,style:{height:"140px"}},null,-1),f={style:{"margin-top":"15px"}},x={style:{"margin-top":"5px"}},C={style:{"margin-top":"20px"}},w=a("\u66F4\u65B0\u65E5\u5FD7"),V=a("gitee");function k(e,E,I,N,c,i){const s=h,d=l;return p(),g(d,{shadow:"hover",header:"\u7248\u672C\u4FE1\u606F"},{default:o(()=>[t("div",v,[y,t("h2",f,"SCUI VITE ADMIN "+n(e.$CONFIG.CORE_VER),1),t("p",x,"\u6700\u65B0\u7248\u672C "+n(c.ver),1)]),t("div",C,[r(s,{type:"primary",plain:"",round:"",onClick:i.golog},{default:o(()=>[w]),_:1},8,["onClick"]),r(s,{type:"primary",plain:"",round:"",onClick:i.gogit},{default:o(()=>[V]),_:1},8,["onClick"])])]),_:1})}const $=_(u,[["render",k]]);export{$ as default};