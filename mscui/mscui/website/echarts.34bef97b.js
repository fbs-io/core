import{_ as l,U as p,V as d,W as h,o as _,e as u,w as f,f as m}from"./index.24f9f7c9.js";/* empty css                */import{_ as r}from"./index.89b34818.js";const v={title:"\u5B9E\u65F6\u6536\u5165",icon:"el-icon-data-line",description:"Echarts\u7EC4\u4EF6\u6F14\u793A",components:{scEcharts:r},data(){return{loading:!0,option:{}}},created(){var n=this;setTimeout(function(){n.loading=!1},500);var a={tooltip:{trigger:"axis"},xAxis:{boundaryGap:!1,type:"category",data:function(){for(var e=new Date,t=[],o=30;o--;)t.unshift(e.toLocaleTimeString().replace(/^\D*/,"")),e=new Date(e-2e3);return t}()},yAxis:[{type:"value",name:"\u4EF7\u683C",splitLine:{show:!1}}],series:[{name:"\u6536\u5165",type:"line",symbol:"none",lineStyle:{width:1,color:"#409EFF"},areaStyle:{opacity:.1,color:"#79bbff"},data:function(){for(var e=[],t=30;t--;)e.push(Math.round(Math.random()*0));return e}()}]};this.option=a},mounted(){}};function g(n,a,e,t,o,w){const i=r,s=p,c=d;return h((_(),u(s,{shadow:"hover",header:"\u5B9E\u65F6\u6536\u5165"},{default:f(()=>[m(i,{ref:"c1",height:"300px",option:o.option},null,8,["option"])]),_:1})),[[c,o.loading]])}const b=l(v,[["render",g]]);export{b as default};