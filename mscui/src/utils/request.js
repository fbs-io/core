import axios from 'axios';
import { ElNotification, ElMessageBox } from 'element-plus';
import sysConfig from "@/config";
import tool from '@/utils/tool';
import router from '@/router';

axios.defaults.baseURL = import.meta.env.VITE_APP_API_BASEURL

axios.defaults.timeout = sysConfig.TIMEOUT

axios.defaults.withCredentials = true
// HTTP request 拦截器
axios.interceptors.request.use(
	(config) => {
		let cookie = tool.cookie.get("TOKEN");
		if(cookie){
			config.headers[sysConfig.TOKEN_NAME] = cookie
		}
		if(!sysConfig.REQUEST_CACHE && config.method == 'get'){
			config.params = config.params || {};
			config.params['_'] = new Date().getTime();
		}
		Object.assign(config.headers, sysConfig.HEADERS)
		return config;
	},
	(error) => {
		return Promise.reject(error);
	}
);
var isReloadLogin = false
// HTTP response 拦截器
axios.interceptors.response.use(
	(response) => {
		// console.log(response)
		var token = response.headers.get("X-Csrf-Token")
		if(response.data.errno == 0 && token){
			tool.cookie.set("TOKEN", token, {
				expires: 30*60
			})
		}
		if (response.data.errno ==0 ){
			if( response.data.notify){
				ElNotification.success({
					title: '请求成功',
					message: "Success"
				})
			}
		}else if (response.data.errno ==40001) {
			if (!isReloadLogin){
				isReloadLogin = true
				ElMessageBox.confirm('当前用户已被登出或无权限访问当前资源，请尝试重新登录后再操作。', '无权限访问', {
					 type: 'error',
					 closeOnClickModal: false,
					 center: true,
					 confirmButtonText: '重新登录'
				}).then(() => {
						router.replace({path: '/login'});
						isReloadLogin = false
					}).catch(() => {})
			}
		}else{
				ElNotification.error({
					title: '请求发生错误',
					message: response.data.api + response.data.message
				})
			

		}
		// console.log(response)
		return response;
	},
	(error) => {
		if (error.response) {
			if (error.response.status == 404) {
				ElNotification.error({
					title: '请求错误',
					message: "Status:404，正在请求不存在的服务器记录！"
				});
			} else if (error.response.status == 500) {
				ElNotification.error({
					title: '请求错误',
					message: error.response.data.message || "Status:500，服务器发生错误！"
				});
			} else if (error.response.status == 401) {
				ElMessageBox.confirm('当前用户已被登出或无权限访问当前资源，请尝试重新登录后再操作。', '无权限访问', {
					type: 'error',
					closeOnClickModal: false,
					center: true,
					confirmButtonText: '重新登录'
				}).then(() => {
					router.replace({path: '/login'});
				}).catch(() => {})
			} else {
				ElNotification.error({
					title: '请求错误',
					message: error.message || `Status:${error.response.status}，未知错误！`
				});
			}
		} else {
			ElNotification.error({
				title: '请求错误',
				message: "请求服务器无响应！"
			});
		}

		return Promise.reject(error.response);
	}
);
// // HTTP request 拦截器
// axios.interceptors.request.use(
// 	(config) => {
// 		let cookie = tool.cookie.get("TOKEN");
// 		if(cookie){
// 			config.headers[sysConfig.TOKEN_NAME] = cookie
// 		}
// 		if(!sysConfig.REQUEST_CACHE && config.method == 'get'){
// 			config.params = config.params || {};
// 			config.params['_'] = new Date().getTime();
// 		}
// 		Object.assign(config.headers, sysConfig.HEADERS)
// 		return config;
// 	},
// 	(error) => {
// 		return Promise.reject(error);
// 	}
// );

// // HTTP response 拦截器
// axios.interceptors.response.use(
// 	(response) => {
// 		return response;
// 	},
// 	(error) => {
// 		if (error.response) {
// 			if (error.response.status == 404) {
// 				ElNotification.error({
// 					title: '请求错误',
// 					message: "Status:404，正在请求不存在的服务器记录！"
// 				});
// 			} else if (error.response.status == 500) {
// 				ElNotification.error({
// 					title: '请求错误',
// 					message: error.response.data.message || "Status:500，服务器发生错误！"
// 				});
// 			} else if (error.response.status == 401) {
// 				ElMessageBox.confirm('当前用户已被登出或无权限访问当前资源，请尝试重新登录后再操作。', '无权限访问', {
// 					type: 'error',
// 					closeOnClickModal: false,
// 					center: true,
// 					confirmButtonText: '重新登录'
// 				}).then(() => {
// 					router.replace({path: '/login'});
// 				}).catch(() => {})
// 			} else {
// 				ElNotification.error({
// 					title: '请求错误',
// 					message: error.message || `Status:${error.response.status}，未知错误！`
// 				});
// 			}
// 		} else {
// 			ElNotification.error({
// 				title: '请求错误',
// 				message: "请求服务器无响应！"
// 			});
// 		}

// 		return Promise.reject(error.response);
// 	}
// );

var http = {

	/** get 请求
	 * @param  {接口地址} url
	 * @param  {请求参数} params
	 * @param  {参数} config
	 */
	get: function(url, params={}, config={}) {
		return new Promise((resolve, reject) => {
			axios({
				method: 'get',
				url: url,
				params: params,
				...config
			}).then((response) => {
				resolve(response.data);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** post 请求
	 * @param  {接口地址} url
	 * @param  {请求参数} data
	 * @param  {参数} config
	 */
	post: function(url, data={}, config={}) {
		return new Promise((resolve, reject) => {
			axios({
				method: 'post',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response.data);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** put 请求
	 * @param  {接口地址} url
	 * @param  {请求参数} data
	 * @param  {参数} config
	 */
	put: function(url, data={}, config={}) {
		return new Promise((resolve, reject) => {
			axios({
				method: 'put',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response.data);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** patch 请求
	 * @param  {接口地址} url
	 * @param  {请求参数} data
	 * @param  {参数} config
	 */
	patch: function(url, data={}, config={}) {
		return new Promise((resolve, reject) => {
			axios({
				method: 'patch',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response.data);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** delete 请求
	 * @param  {接口地址} url
	 * @param  {请求参数} data
	 * @param  {参数} config
	 */
	delete: function(url, data={}, config={}) {
		return new Promise((resolve, reject) => {
			axios({
				method: 'delete',
				url: url,
				data: data,
				...config
			}).then((response) => {
				resolve(response.data);
			}).catch((error) => {
				reject(error);
			})
		})
	},

	/** jsonp 请求
	 * @param  {接口地址} url
	 * @param  {JSONP回调函数名称} name
	 */
	jsonp: function(url, name='jsonp'){
		return new Promise((resolve) => {
			var script = document.createElement('script')
			var _id = `jsonp${Math.ceil(Math.random() * 1000000)}`
			script.id = _id
			script.type = 'text/javascript'
			script.src = url
			window[name] =(response) => {
				resolve(response)
				document.getElementsByTagName('head')[0].removeChild(script)
				try {
					delete window[name];
				}catch(e){
					window[name] = undefined;
				}
			}
			document.getElementsByTagName('head')[0].appendChild(script)
		})
	}
}

export default http;
