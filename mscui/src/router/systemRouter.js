/*
 * @Author: reel
 * @Date: 2023-06-04 15:27:47
 * @LastEditors: reel
 * @LastEditTime: 2023-09-09 09:14:10
 * @Description: 请填写简介
 */
import config from "@/config"

//系统路由
const routes = [
	{
		name: "layout",
		path: "/",
		component: () => import('@/layout'),
		redirect: config.DASHBOARD_URL || "/overview",
		children: []
	},
	{
		path: "/login",
		component: () => import('@/views/login'),
		meta: {
			title: "登录"
		}
	},
	{
		path: "/install",
		// hidden: true,
		component: () => import('@/views/install'),
	}

]

export default routes;