import{_,o as i,c,g as d,u as h}from"./index.8c95325c.js";const f={props:{modelValue:{type:String,default:""}},data(){return{level:0}},watch:{modelValue(){this.strength(this.modelValue)}},mounted(){this.strength(this.modelValue)},methods:{strength(s){var e=0,a=s.length>=6,l=/\d/.test(s),t=/[a-z]/.test(s),r=/[A-Z]/.test(s),n=!/(\w)\1{2}/.test(s),o=/[`~!@#$%^&*()_+<>?:"{},./;'[\]]/.test(s);if(s.length<=0)return e=0,this.level=e,!1;if(!a)return e=1,this.level=e,!1;l&&(e+=1),t&&(e+=1),r&&(e+=1),n&&(e+=1),o&&(e+=1),this.level=e}}},p={class:"sc-password-strength"};function u(s,e,a,l,t,r){return i(),c("div",p,[d("div",{class:h(["sc-password-strength-bar",`sc-password-strength-level-${t.level}`])},null,2)])}const g=_(f,[["render",u],["__scopeId","data-v-691fbf79"]]);export{g as default};
