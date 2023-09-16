/*
 * @Author: reel
 * @Date: 2023-06-04 15:27:47
 * @LastEditors: reel
 * @LastEditTime: 2023-06-07 22:02:58
 * @Description: 请填写简介
 */
import config from "@/config"
import http from "@/utils/request"

export default {
	token: {
		url: `${config.API_URL}/login`,
		name: "登录获取Cookie",
		post: async function(data={}){
			return await http.post(this.url, data);
		}
	}
}
