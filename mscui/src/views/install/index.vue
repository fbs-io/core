<!--
 * @Author: reel
 * @Date: 2023-05-30 20:47:44
 * @LastEditors: reel
 * @LastEditTime: 2023-09-12 05:27:27
 * @Description: 请填写简介
-->

<template>
	<common-page title="项目初始化">
		<el-steps :active="stepActive" simple finish-status="success">
		    <el-step title="服务配置" />
		    <el-step title="用户配置" />
			<el-step title="数据库配置" />
			<el-step title="缓存配置" />
			<el-step title="数据配置" />
		    <el-step title="完成注册" />
		</el-steps>
		<el-form v-if="stepActive==0" ref="stepForm_0" :model="form" :rules="rules" :label-width="120">
			<el-form-item label="端口号" prop="port">
				<el-input v-model="form.port" placeholder="请输入服务器端口号"></el-input>
				<div class="el-form-item-msg">该端口作为应用端口，非管理端端口号。</div>
			</el-form-item>
		</el-form>
		<el-form v-if="stepActive==1" ref="stepForm_1" :model="form" :rules="rules" :label-width="120">
			<el-form-item label="登录账号" prop="user">
				<el-input v-model="form.user" placeholder="请输入登录账号"></el-input>
				<div class="el-form-item-msg">登录账号将作为登录时的唯一凭证</div>
			</el-form-item>
			<el-form-item label="登录密码" prop="password">
				<el-input v-model="form.password" type="password" show-password placeholder="请输入登录密码"></el-input>
				<sc-password-strength v-model="form.password"></sc-password-strength>
				<div class="el-form-item-msg">默认密码root123，请修改密码并妥善保管</div>
			</el-form-item>
			<el-form-item label="确认密码" prop="password2">
				<el-input v-model="form.password2" type="password" show-password placeholder="请再一次输入登录密码"></el-input>
			</el-form-item>
			<el-form-item label="" prop="agree">
				<el-checkbox v-model="form.agree" label="">已阅读并同意</el-checkbox><span class="link" @click="showAgree=true">《平台服务协议》</span>
			</el-form-item>
		</el-form>
		<el-form v-if="stepActive==2" ref="stepForm_2" :model="form" :rules="rules" :label-width="120">
            <el-form-item label="数据库名称" prop="db_name">
                <el-input v-model="form.db_name" placeholder="请输入数据库名称"></el-input>
            </el-form-item>
			<el-form-item label="数据库地址" prop="db_host" v-if="isShowDB">
				<el-input v-model="form.db_host" placeholder="请输入数据库地址"></el-input>
			</el-form-item>
			<el-form-item label="数据库端口" prop="db_port" v-if="isShowDB">
				<el-input v-model="form.db_port" placeholder="请输入数据库端口"></el-input>
			</el-form-item>
			<el-form-item label="数据库用户" prop="db_user" v-if="isShowDB">
				<el-input v-model="form.db_user" placeholder="请输入数据库用户"></el-input>
			</el-form-item>
			<el-form-item label="数据库密码" prop="db_pwd" v-if="isShowDB">
				<el-input v-model="form.db_pwd" placeholder="请输入数据库密码"></el-input>
			</el-form-item>
			<el-form-item label="数据库类型" prop="db_type">
				<el-radio-group v-model="form.db_type">
					<el-radio-button label="sqlite" @click="isShowDB=false">Sqlite</el-radio-button>
					<el-radio-button label="postgres" @click="isShowDB=true">Postgres</el-radio-button>
					<!-- <el-radio-button label="mysql" @click="isShowDB=true">MySQL</el-radio-button> -->
				</el-radio-group>
			</el-form-item>

		</el-form>
		<el-form v-if="stepActive==3" ref="stepForm_3" :model="form" :rules="rules" :label-width="120">
            <el-form-item label="缓存名称" prop="cache_name">
                <el-input v-model="form.cache_name" placeholder="请输入数据库名称"></el-input>
            </el-form-item>
			<el-form-item label="缓存地址" prop="cache_host" v-if="isShowCache">
				<el-input v-model="form.cache_host" placeholder="请输入数据库地址"></el-input>
			</el-form-item>
			<el-form-item label="缓存端口" prop="cache_port" v-if="isShowCache">
				<el-input v-model="form.cache_port" placeholder="请输入数据库端口"></el-input>
			</el-form-item>
			<el-form-item label="缓存用户" prop="cache_user" v-if="isShowCache">
				<el-input v-model="form.cache_user" placeholder="请输入数据库用户"></el-input>
			</el-form-item>
			<el-form-item label="缓存密码" prop="cache_pwd" v-if="isShowCache">
				<el-input v-model="form.cache_pwd" placeholder="请输入数据库密码"></el-input>
			</el-form-item>
			<el-form-item label="缓存类型" prop="cache_type">
				<el-radio-group v-model="form.cache_type">
					<el-radio-button label="local" @click="isShowCache=false">Local</el-radio-button>
					<!-- <el-radio-button label="redis" @click="isShowCache=true">Redis</el-radio-button> -->
				</el-radio-group>
			</el-form-item>
		</el-form>
		<el-form v-if="stepActive==4" ref="stepForm_4" :model="form" :rules="rules" :label-width="120">
			<el-form-item label="数据存储配置" prop="data_path">
				<el-input v-model="form.data_path" placeholder="请输入数据存储位置"></el-input>
				<div class="el-form-item-msg">该地址作为日志，缓存，数据库等数据存储路径，默认在程序启动目录下。如需更改，请手动填写数据存储路径</div>
			</el-form-item>
		</el-form>
		<div v-if="stepActive==5">
			<el-result icon="success" title="配置成功" sub-title="可以使用账户密码登录后台管理系统">
				<template #extra>
					<el-button type="primary" @click="goMscLogin">后台登录</el-button>
					<el-button type="primary" @click="goAppLogin">应用登录</el-button>
				</template>
			</el-result>
		</div>
		<el-form style="text-align: center;">
			<el-button v-if="stepActive>0 && stepActive<5" @click="pre">上一步</el-button>
			<el-button v-if="stepActive<4" type="primary" @click="next">下一步</el-button>
			<el-button v-if="stepActive==4" type="primary" :loading="isloading"  @click="save">提交</el-button>
		</el-form>
		
		<!-- 软件服务协议 -->
		<el-dialog v-model="showAgree" title="软件服务协议" :width="800" destroy-on-close>
          <agreePage/>
			<template #footer>
				<el-button @click="showAgree=false">取消</el-button>
				<el-button type="primary" @click="showAgree=false;form.agree=true;">我已阅读并同意</el-button>
			</template>
		</el-dialog>

		<!-- 缓存和数据库配置时提醒 -->
        <el-dialog v-model="showWarning" title="警告"  :width="800" destroy-on-close>
            <p style="font-size: 16px; color: brown ;">请注意, 强烈不建议您在生产环境使用{{db_cache}}。</p>
			<template #footer>
				<el-button @click="showWarning=false">取消</el-button>
				<el-button type="primary" @click="showWarning=false;stepActive += 1">确定</el-button>
			</template>
		</el-dialog>
	</common-page>
</template>

<script>
	import scPasswordStrength from '@/components/scPasswordStrength/index.vue';
	import commonPage from './components/commonPage.vue'
	import agreePage from './components/agree.vue'
	import axios from 'axios';

	export default {
		components: {
			commonPage,
			agreePage,
			scPasswordStrength
		},
		data() {
			return {
				stepActive: 0,
				showAgree: false,
				isloading: false,
                isShowDB:false,
                isShowCache:false,
                showWarning: false,
                db_cache: "Sqlite数据库",
				form: {
					port:"80",
					user: "root",
					password: "root123",
					password2: "root123",
					agree: false,
                    // 数据库相关
					db_name: "",
					db_host: "",
					db_port: "",
					db_user: "",
					db_pwd: "",
					db_type: "sqlite",
                    // 缓存相关
					cache_name: "",
					cache_host: "",
					cache_port: "",
					cache_user: "",
					cache_pwd: "",
					cache_type: "local",
					// 数据存放位置
					data_path: "data"
				},
				rules: {
                    user: [
						{ required: true, message: '请输入账号名'}
					],
					password: [
						{ required: true, message: '请输入密码'}
					],
					password2: [
						{ required: true, message: '请再次输入密码'},
						{validator: (rule, value, callback) => {
							if (value !== this.form.password) {
								callback(new Error('两次输入密码不一致'));
							}else{
								callback();
							}
						}}
					],
					agree: [
						{validator: (rule, value, callback) => {
							if (!value) {
								callback(new Error('请阅读并同意协议'));
							}else{
								callback();
							}
						}}
					],
                    db_name: [
                        { required: true, message: '请输入数据库名称'}
                    ],
                    db_host: [
						{ required: true, message: '请输入数据库地址'}
					],
					db_port: [
						{ required: true, message: '请输入数据库端口'}
					],
					db_user: [
						{ required: true, message: '请输入数据库用户'}
					],
					db_pwd: [
						{ required: true, message: '请输入数据库密码'}
					],
					db_type: [
						{ required: true, message: '请选择数据库类型'}
					],
                    cache_name: [
                        { required: true, message: '请输入数据库名称'}
                    ],
                    cache_host: [
						{ required: true, message: '请输入数据库地址'}
					],
					cache_port: [
						{ required: true, message: '请输入数据库端口'}
					],
					cache_user: [
						{ required: true, message: '请输入数据库用户'}
					],
					cache_pwd: [
						{ required: true, message: '请输入数据库密码'}
					],
					cache_type: [
						{ required: true, message: '请选择缓存类型'}
					],
                },

			}
		},
		mounted() {
			// console.log(location)
		},
		methods: {
			pre(){
				this.stepActive -=1
			},
			next(){
				const formName = `stepForm_${this.stepActive}`
				this.$refs[formName].validate((valid) => {
					if (valid) {
                        
                        if ((this.form.db_type=="sqlite" && this.stepActive==2 )|| (this.form.cache_type =="local"&& this.stepActive==3)){
                            this.showWarning = true
                            this.checkDBCache()
                            return true
                        }
						this.stepActive += 1
					}else{
						return false
					}
				})

			},
			async save(){
				this.isloading = true
				const formName = `stepForm_${this.stepActive}`
				this.$refs[formName].validate((valid) => {
					if (!valid) {
						return false
					}
			})
			
			var res = await this.$API.install.install.post(this.form)
			this.isloading = false
			if (res.errno == 0){
				this.stepActive += 1
			}else{
				return false
			}
			// this.isloading = false

		
			},
            checkDBCache(){
                if (this.stepActive==2){
                    this.db_cache = "Sqlite数据库"
                }else if (this.stepActive==3){
                    this.db_cache = "本地缓存"
                }
            },
			goMscLogin(){
				this.$CONFIG.APP_INIT = false
				this.$router.replace({
					path: '/',
				})
			},
			goAppLogin(){
				location.replace(location.protocol+"//"+location.hostname+":"+this.form.port)
			}
		}
	}
</script>

<style scoped>


</style>

