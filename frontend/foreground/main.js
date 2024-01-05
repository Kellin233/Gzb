import Vue from 'vue'
import App from './App'

Vue.config.productionTip = false

import uView from 'uview-ui'
Vue.use(uView)

App.mpType = 'app'


const msg = (title, duration=1500, mask=false, icon='none')=>{
	//统一提示方便全局修改
	if(Boolean(title) === false){
		return;
	}
	uni.showToast({
		title,
		duration,
		mask,
		icon
	});
}

Vue.prototype.$api = {msg};



const app = new Vue({
    ...App
})

// 引入请求封装，将app参数传递到配置中
require('./common/request.js')(app)

// http拦截器，将此部分放在new Vue()和app.$mount()之间，才能App.vue中正常使用
// import httpInterceptor from '@/common/http.interceptor.js'
// Vue.use(httpInterceptor, app)

// http接口API集中管理引入部分
// import httpApi from '@/common/http.api.js'
// Vue.use(httpApi, app)

app.$mount()
