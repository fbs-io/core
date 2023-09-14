/*
 * @Author: reel
 * @Date: 2023-09-09 10:35:05
 * @LastEditors: reel
 * @LastEditTime: 2023-09-09 10:35:41
 * @Description: 请填写简介
 */
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
	install: {
		url: `${config.API_URL}/install`,
		name: "初次安装初始化",
		post: async function(data={}){
			return await http.post(this.url, data);
		}
	}
}
