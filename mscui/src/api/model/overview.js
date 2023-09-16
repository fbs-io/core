/*
 * @Author: reel
 * @Date: 2023-06-10 21:44:10
 * @LastEditors: reel
 * @LastEditTime: 2023-06-14 20:51:48
 * @Description: 请填写简介
 */
import config from "@/config"
import http from "@/utils/request"
export default {
	hostinfo: {
		url: `${config.API_URL}/hostinfo`,
		name: "获取服务器信息",
		get: async function(data={}){
			return await http.get(this.url, data);
		}
	},
	sysinfo: {
		url: `${config.API_URL}/sysinfo`,
		name: "获取服务器资源信息",
		get: async function(data={}){
			return await http.get(this.url, data);
		}
	},
	processinfo: {
		url: `${config.API_URL}/processinfo`,
		name: "获取进程相关信息",
		get: async function(data={}){
			return await http.get(this.url, data);
		}
	},
	srvstatus: {
		url: `${config.API_URL}/srvstatus`,
		name: "获取服务状态信息",
		get: async function(data={}){
			return await http.get(this.url, data);
		}
	},
	srvrestart: {
		url: `${config.API_URL}/srvrestart`,
		name: "重启服务",
		post: async function(data={}){
			return await http.post(this.url, data);
		}
	}
}
