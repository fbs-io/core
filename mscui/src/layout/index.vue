<template>
		<header class="adminui-header">
			<div class="adminui-header-left">
				<div class="logo-bar">
					<img class="logo" src="/img/logo.png">
					<span>{{ $CONFIG.APP_NAME }}</span>
				</div>
				<!-- <div v-if="!ismobile" style="width: 5rem;">
					<el-menu mode="horizontal" :default-active="active" router background-color="#222b45" text-color="#fff" active-text-color="var(--el-color-primary)">
						<NavMenu :navMenus="menu"></NavMenu>
					</el-menu>
				</div> -->
				<div v-if="!ismobile" class="adminui-header-menu">
					<el-menu mode="horizontal" :default-active="active" router background-color="#222b45" text-color="#fff" active-text-color="var(--el-color-primary)">
						<NavMenu :navMenus="menu"></NavMenu>
					</el-menu>
				</div>
			</div>
			<div class="adminui-header-right">

				<Side-m v-if="ismobile"></Side-m>
				<userbar></userbar>
			</div>
		</header>
		<section class="aminui-wrapper">
			<div class="aminui-body el-container">
				<Tags v-if="!ismobile && layoutTags" ></Tags>
				<div class="adminui-main" id="adminui-main">
					<router-view v-slot="{ Component }">
					    <keep-alive :include="this.$store.state.keepAlive.keepLiveRoute">
					        <component :is="Component" :key="$route.fullPath" v-if="$store.state.keepAlive.routeShow"/>
					    </keep-alive>
					</router-view>
					<iframe-view></iframe-view>
				</div>
			</div>
		</section>
	<!-- </template> -->

	<!-- 默认布局 -->


	<div class="main-maximize-exit" @click="exitMaximize"><el-icon><el-icon-close /></el-icon></div>
</template>

<script>
	import SideM from './components/sideM.vue';
	import Topbar from './components/topbar.vue';
	import Tags from './components/tags.vue';
	import NavMenu from './components/NavMenu.vue';
	import userbar from './components/userbar.vue';
	import setting from './components/setting.vue';
	import iframeView from './components/iframeView.vue';

	export default {
		name: 'index',
		components: {
			SideM,
			Topbar,
			Tags,
			NavMenu,
			userbar,
			setting,
			iframeView
		},
		data() {
			return {
				settingDialog: false,
				menu: [],
				nextMenu: [],
				pmenu: {},
				active: ''
			}
		},
		computed:{
			ismobile(){
				return this.$store.state.global.ismobile
			},
			layout(){
				return this.$store.state.global.layout
			},
			layoutTags(){
				return this.$store.state.global.layoutTags
			},
			menuIsCollapse(){
				return this.$store.state.global.menuIsCollapse
			}
		},
		created() {
			this.onLayoutResize();
			window.addEventListener('resize', this.onLayoutResize);
			var menu = this.$router.sc_getMenu();
			this.menu = this.filterUrl(menu);
			this.showThis()
		},
		watch: {
			$route() {
				this.showThis()
			},
			layout: {
				handler(val){
					document.body.setAttribute('data-layout', val)
				},
				immediate: true,
			}
		},
		methods: {
			openSetting(){
				this.settingDialog = true;
			},
			onLayoutResize(){
				this.$store.commit("SET_ismobile", document.body.clientWidth < 992)
			},
			//路由监听高亮
			showThis(){
				this.pmenu = this.$route.meta.breadcrumb ? this.$route.meta.breadcrumb[0] : {}
				this.nextMenu = this.filterUrl(this.pmenu.children);
				this.$nextTick(()=>{
					this.active = this.$route.meta.active || this.$route.fullPath;
				})
			},
			//点击显示
			showMenu(route) {
				this.pmenu = route;
				this.nextMenu = this.filterUrl(route.children);
				if((!route.children || route.children.length == 0) && route.component){
					this.$router.push({path: route.path})
				}
			},
			//转换外部链接的路由
			filterUrl(map){
				var newMap = []
				map && map.forEach(item => {
					item.meta = item.meta?item.meta:{};
					//处理隐藏
					if(item.meta.hidden || item.meta.type=="button"){
						return false
					}
					//处理http
					if(item.meta.type=='iframe'){
						item.path = `/i/${item.name}`;
					}
					//递归循环
					if(item.children&&item.children.length > 0){
						item.children = this.filterUrl(item.children)
					}
					newMap.push(item)
				})
				return newMap;
			},
			//退出最大化
			exitMaximize(){
				document.getElementById('app').classList.remove('main-maximize')
			}
		}
	}
</script>
