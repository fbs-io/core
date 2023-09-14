import { markRaw } from 'vue';
const resultComps = {}

// const files = import.meta.globEager('./*.vue')

const files = import.meta.glob('./*.vue', {import: 'default',eager: true})

Object.keys(files).forEach(fileName => {
	let comp = files[fileName]
	resultComps[fileName.replace(/^\.\/(.*)\.\w+$/, '$1')] = comp
})

export default markRaw(resultComps)