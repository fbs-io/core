import config from "./config"
import api from './api'
import tool from './utils/tool'
import http from "./utils/request"
import { permission, rolePermission } from './utils/permission'

import scTable from './components/scTable'
import scTableColumn from './components/scTable/column.js'
// import scFilterBar from './components/scFilterBar'
import scUpload from './components/scUpload'
import scUploadMultiple from './components/scUpload/multiple'
import scUploadFile from './components/scUpload/file'
import scFormTable from './components/scFormTable'
import scTableSelect from './components/scTableSelect'
import scPageHeader from './components/scPageHeader'
import scSelect from './components/scSelect'
import scDialog from './components/scDialog'
import scForm from './components/scForm'
import scTitle from './components/scTitle'
import scWaterMark from './components/scWaterMark'
import scQrCode from './components/scQrCode'

import scStatusIndicator from './components/scMini/scStatusIndicator'
import scTrend from './components/scMini/scTrend'

import auth from './directives/auth'
import role from './directives/role'
import time from './directives/time'
import copy from './directives/copy'
import errorHandler from './utils/errorHandler'

import { ElMessage, ElMessageBox, ElLoading, ElAlert } from 'element-plus'
import "element-plus/es/components/message/style/css"
import "element-plus/es/components/message-box/style/css"
import "element-plus/es/components/loading/style/css"
import "element-plus/es/components/alert/style/css"

import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import * as elIcons from '@element-plus/icons-vue'
import * as scIcons from './assets/icons'

export default {
	install(app) {
		//挂载全局对象
		app.config.globalProperties.$CONFIG = config;
		app.config.globalProperties.$TOOL = tool;
		app.config.globalProperties.$HTTP = http;
		app.config.globalProperties.$API = api;
		app.config.globalProperties.$AUTH = permission;
		app.config.globalProperties.$ROLE = rolePermission;
		app.config.globalProperties.$confirm = ElMessageBox.confirm;
		app.config.globalProperties.$message = ElMessage;
		app.config.globalProperties.$loading = ElLoading.service;
		app.config.globalProperties.$alert = ElMessageBox.alert;

		//注册全局组件
		app.component('scTable', scTable);
		app.component('scTableColumn', scTableColumn);
		// app.component('scFilterBar', scFilterBar);
		app.component('scUpload', scUpload);
		app.component('scUploadMultiple', scUploadMultiple);
		app.component('scUploadFile', scUploadFile);
		app.component('scFormTable', scFormTable);
		app.component('scTableSelect', scTableSelect);
		app.component('scPageHeader', scPageHeader);
		app.component('scSelect', scSelect);
		app.component('scDialog', scDialog);
		app.component('scForm', scForm);
		app.component('scTitle', scTitle);
		app.component('scWaterMark', scWaterMark);
		app.component('scQrCode', scQrCode);
		app.component('scStatusIndicator', scStatusIndicator);
		app.component('scTrend', scTrend);

		//注册全局指令
		app.directive('auth', auth)
		app.directive('role', role)
		app.directive('time', time)
		app.directive('copy', copy)

		//统一注册el-icon图标
		for(let icon in elIcons){
			app.component(`ElIcon${icon}`, elIcons[icon])
		}
		//统一注册sc-icon图标
		for(let icon in scIcons){
			app.component(`ScIcon${icon}`, scIcons[icon])
		}

		// el的图标组件
		for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
			app.component(key, component)
		}

		//关闭async-validator全局控制台警告
		window.ASYNC_VALIDATOR_NO_WARNING = 1

		//全局代码错误捕捉
		app.config.errorHandler = errorHandler
	}
}