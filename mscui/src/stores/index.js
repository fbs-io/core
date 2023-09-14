/**
 * @description 自动import导入所有 vuex 模块
 */

import { createStore } from 'vuex';

// const files = require.context('./modules', false, /\.js$/);
const modules = {}
// files.keys().forEach((key) => {
// 	modules[key.replace(/(\.\/|\.js)/g, '')] = files(key).default
// })


// const metas = import.meta.globEager('./modules/*.js')

const metas = import.meta.glob('./modules/*.js', {import: 'default',eager: true})

for (let key in metas) {
	let k = key.replace('modules/', '')
	modules[k.replace(/(\.\/|\.js)/g, '')] = metas[key]
}

export default createStore({
	modules
});
