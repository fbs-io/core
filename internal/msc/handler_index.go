/*
 * @Author: reel
 * @Date: 2023-05-28 14:59:38
 * @LastEditors: reel
 * @LastEditTime: 2023-06-23 21:02:25
 * @Description: msc首页
 */
package msc

import (
	"github.com/gin-gonic/gin"
)

func (m *handler) index() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		index := indexIndex
		if !m.config.IsLoad {
			index = indexInstall
		}
		ctx.String(200, index)

	}
}

var indexIndex = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" href="/mscui/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>FBS Manager System Conter</title>
    <script src="/mscui/config.js"></script>
    <script type="module" crossorigin src="/mscui/assets/index.d39b14b4.js"></script>
    <link rel="stylesheet" href="/mscui/assets/index.785fc723.css">
  </head>
  <body>
    <noscript>
			<strong>We're sorry but <%= VUE_APP_TITLE %> doesn't work properly without JavaScript
				enabled. Please enable it to continue.</strong>
		</noscript>
    <script type="text/javascript">
			var dark = window.localStorage.getItem('APP_DARK');
			if(dark){
				document.documentElement.classList.add("dark")
			}
		</script>
    <div id="app" class="aminui">
      <div class="app-loading">
				<div class="app-loading__logo">
					<img src="/mscui/img/logo.png"/>
				</div>
				<div class="app-loading__loader"></div>
				<div class="app-loading__title">FBS Manager System Center</div>
			</div>
			<style>
				.app-loading {position: absolute;top:0px;left:0px;right:0px;bottom:0px;display: flex;justify-content: center;align-items: center;flex-direction: column;background: #fff;}
				.app-loading__logo {margin-bottom: 30px;}
				.app-loading__logo img {width: 90px;vertical-align: bottom;}
				.app-loading__loader {box-sizing: border-box;width: 35px;height: 35px;border: 5px solid transparent;border-top-color: #000;border-radius: 50%;animation: .5s loader linear infinite;position: relative;}
				.app-loading__loader:before {box-sizing: border-box;content: '';display: block;width: inherit;height: inherit;position: absolute;top: -5px;left: -5px;border: 5px solid #ccc;border-radius: 50%;opacity: .5;}
				.app-loading__title {font-size: 24px;color: #333;margin-top: 30px;}
				.dark .app-loading {background: #222225;}
				.dark .app-loading__loader {border-top-color: #fff;}
				.dark .app-loading__title {color: #d0d0d0;}
				@keyframes loader {
				    0% {
				        transform: rotate(0deg);
				    }
				    100% {
				        transform: rotate(360deg);
				    }
				}
			</style>
    </div>
    
  </body>
  <div id="versionCheck" style="display: none;position: absolute;z-index: 99;top:0;left:0;right:0;bottom:0;padding:40px;background:rgba(255,255,255,0.9);color: #333;">
		<h2 style="line-height: 1;margin: 0;font-size: 24px;">当前使用的浏览器内核版本过低 :(</h2>
		<p style="line-height: 1;margin: 0;font-size: 14px;margin-top: 20px;opacity: 0.8;">当前版本：<span id="versionCheck-type">--</span> <span id="versionCheck-version">--</span></p>
		<p style="line-height: 1;margin: 0;font-size: 14px;margin-top: 10px;opacity: 0.8;">最低版本要求：Chrome 71+、Firefox 65+、Safari 12+、Edge 97+。</p>
		<p style="line-height: 1;margin: 0;font-size: 14px;margin-top: 10px;opacity: 0.8;">请升级浏览器版本，或更换现代浏览器，如果你使用的是双核浏览器,请切换到极速/高速模式。</p>
	</div>
  <script type="text/javascript">
    function getBrowerInfo(){
      var userAgent = window.navigator.userAgent;
      var browerInfo = {
        type: 'unknown',
        version: 'unknown',
        userAgent: userAgent
      };
      if(document.documentMode){
        browerInfo.type = "IE"
        browerInfo.version = document.documentMode + ''
      }else if(indexOf(userAgent, "Firefox")){
        browerInfo.type = "Firefox"
        browerInfo.version = userAgent.match(/Firefox\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Opera")){
        browerInfo.type = "Opera"
        browerInfo.version = userAgent.match(/Opera\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Edg")){
        browerInfo.type = "Edg"
        browerInfo.version = userAgent.match(/Edg\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Chrome")){
        browerInfo.type = "Chrome"
        browerInfo.version = userAgent.match(/Chrome\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Safari")){
        browerInfo.type = "Safari"
        browerInfo.version = userAgent.match(/Safari\/([\d.]+)/)[1]
      }
      return browerInfo
    }
      function indexOf(userAgent, brower) {
          return userAgent.indexOf(brower) > -1
      }
    function isSatisfyBrower(){
      var minVer = {
        "Chrome": 71,
        "Firefox": 65,
        "Safari": 12,
        "Edg": 97,
        "IE": 999
      }
      var browerInfo = getBrowerInfo()
      var materVer = browerInfo.version.split('.')[0]
          return materVer >= minVer[browerInfo.type]
    }
    if(!isSatisfyBrower()){
      document.getElementById('versionCheck').style.display = 'block';
      document.getElementById('versionCheck-type').innerHTML = getBrowerInfo().type;
      document.getElementById('versionCheck-version').innerHTML = getBrowerInfo().version;
    }
    </script>
</html>
`
var indexInstall = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" href="/install/img/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>MSC</title>
    <script type="module" crossorigin src="/install/assets/index.428d4807.js"></script>
    <link rel="stylesheet" href="/install/assets/index.74a5aef5.css">
  </head>
  <body>
    <noscript>
			<strong>We're sorry but <%= VUE_APP_TITLE %> doesn't work properly without JavaScript
				enabled. Please enable it to continue.</strong>
		</noscript>
    <script type="text/javascript">
			var dark = window.localStorage.getItem('APP_DARK');
			if(dark){
				document.documentElement.classList.add("dark")
			}
		</script>
    <div id="app" class="aminui">
      <div class="app-loading">
				<div class="app-loading__logo">
					<img src="install/img/logo.png"/>
				</div>
				<div class="app-loading__loader"></div>
				<div class="app-loading__title">scuiViteUI</div>
			</div>
			<style>
				.app-loading {position: absolute;top:0px;left:0px;right:0px;bottom:0px;display: flex;justify-content: center;align-items: center;flex-direction: column;background: #fff;}
				.app-loading__logo {margin-bottom: 30px;}
				.app-loading__logo img {width: 90px;vertical-align: bottom;}
				.app-loading__loader {box-sizing: border-box;width: 35px;height: 35px;border: 5px solid transparent;border-top-color: #000;border-radius: 50%;animation: .5s loader linear infinite;position: relative;}
				.app-loading__loader:before {box-sizing: border-box;content: '';display: block;width: inherit;height: inherit;position: absolute;top: -5px;left: -5px;border: 5px solid #ccc;border-radius: 50%;opacity: .5;}
				.app-loading__title {font-size: 24px;color: #333;margin-top: 30px;}
				.dark .app-loading {background: #222225;}
				.dark .app-loading__loader {border-top-color: #fff;}
				.dark .app-loading__title {color: #d0d0d0;}
				@keyframes loader {
				    0% {
				        transform: rotate(0deg);
				    }
				    100% {
				        transform: rotate(360deg);
				    }
				}
			</style>
    </div>
    
  </body>
  <div id="versionCheck" style="display: none;position: absolute;z-index: 99;top:0;left:0;right:0;bottom:0;padding:40px;background:rgba(255,255,255,0.9);color: #333;">
		<h2 style="line-height: 1;margin: 0;font-size: 24px;">当前使用的浏览器内核版本过低 :(</h2>
		<p style="line-height: 1;margin: 0;font-size: 14px;margin-top: 20px;opacity: 0.8;">当前版本：<span id="versionCheck-type">--</span> <span id="versionCheck-version">--</span></p>
		<p style="line-height: 1;margin: 0;font-size: 14px;margin-top: 10px;opacity: 0.8;">最低版本要求：Chrome 71+、Firefox 65+、Safari 12+、Edge 97+。</p>
		<p style="line-height: 1;margin: 0;font-size: 14px;margin-top: 10px;opacity: 0.8;">请升级浏览器版本，或更换现代浏览器，如果你使用的是双核浏览器,请切换到极速/高速模式。</p>
	</div>
  <script type="text/javascript">
    function getBrowerInfo(){
      var userAgent = window.navigator.userAgent;
      var browerInfo = {
        type: 'unknown',
        version: 'unknown',
        userAgent: userAgent
      };
      if(document.documentMode){
        browerInfo.type = "IE"
        browerInfo.version = document.documentMode + ''
      }else if(indexOf(userAgent, "Firefox")){
        browerInfo.type = "Firefox"
        browerInfo.version = userAgent.match(/Firefox\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Opera")){
        browerInfo.type = "Opera"
        browerInfo.version = userAgent.match(/Opera\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Edg")){
        browerInfo.type = "Edg"
        browerInfo.version = userAgent.match(/Edg\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Chrome")){
        browerInfo.type = "Chrome"
        browerInfo.version = userAgent.match(/Chrome\/([\d.]+)/)[1]
      }else if(indexOf(userAgent, "Safari")){
        browerInfo.type = "Safari"
        browerInfo.version = userAgent.match(/Safari\/([\d.]+)/)[1]
      }
      return browerInfo
    }
      function indexOf(userAgent, brower) {
          return userAgent.indexOf(brower) > -1
      }
    function isSatisfyBrower(){
      var minVer = {
        "Chrome": 71,
        "Firefox": 65,
        "Safari": 12,
        "Edg": 97,
        "IE": 999
      }
      var browerInfo = getBrowerInfo()
      var materVer = browerInfo.version.split('.')[0]
          return materVer >= minVer[browerInfo.type]
    }
    if(!isSatisfyBrower()){
      document.getElementById('versionCheck').style.display = 'block';
      document.getElementById('versionCheck-type').innerHTML = getBrowerInfo().type;
      document.getElementById('versionCheck-version').innerHTML = getBrowerInfo().version;
    }
    </script>
</html>
`
