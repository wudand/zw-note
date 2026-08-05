let list = [
    {
        id: '1',
        title: '前端笔记',
        description: '每天进步一点。日积月累。',
        author: '吴烦恼'
    },
    {
        id: '2',
        title: '后端技术简析',
        description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
        author: 'zww'
    },
    {
        id: '3',
        title: '后端技术简析',
        description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
        author: 'zww'
    },
    {
        id: '4',
        title: 'vue3笔记',
        description: 'vue3笔记',
        author: 'zww'
    },
    {
        id: '5',
        title: 'react笔记',
        description: 'react笔记',
        author: 'zww'
    },
    {
        id: '6',
        title: 'angular笔记',
        description: 'angular笔记',
        author: 'zww'
    },
    {
        id: '7',
        title: 'nodejs笔记',
        description: 'nodejs笔记',
        author: 'zww'
    },
    {
        id: '8',
        title: 'php笔记',
        description: 'php笔记',
        author: 'zww'
    },
    {
        id: '9',
        title: 'python笔记',
        description: 'python笔记',
        author: 'zww'
    },
    {
        id: '10',
        title: 'javascript笔记',
        description: 'javascript笔记',
        author: 'zww'
    }
]

const outline = {
    '1': [
        {
            id: '1',
            title: 'Vue2.0笔记',
            parentId: '0',
            children: [
                { id: '1-1', title: '搭建环境', parentId: '1' },
                { id: '1-2', title: '生命周期', parentId: '1' },
                { id: '1-3', title: '模版语法', parentId: '1' },
                { id: '1-4', title: '计算属性 VS 侦听器', parentId: '1' },
                { id: '1-5', title: '父子组件间的数据传递', parentId: '1' },
                { id: '1-6', title: '非父子组件之间的传值', parentId: '1' },
                { id: '1-7', title: '插槽(slot)', parentId: '1' }
            ]
        },
        {
            id: '2',
            title: '2. 常用命令',
            parentId: '0',
            children: [
                { id: '2-1', title: '2.1 镜像操作', parentId: '2' },
                { id: '2-2', title: '2.2 容器操作', parentId: '2' },
                { id: '2-3', title: '2.3 日志和调试', parentId: '2' },
            ]
        },
        {
            id: '3',
            title: '3. Dockerfile',
            parentId: '0',
            children: [
                { id: '3-1', title: '3.1 基本语法', parentId: '3' },
                { id: '3-2', title: '3.2 构建镜像', parentId: '3' },
            ]
        },
        {
            id: '4',
            title: '4. Docker Compose',
            parentId: '0',
            children: [
                { id: '4-1', title: '4.1 配置文件', parentId: '4' },
                { id: '4-2', title: '4.2 常用命令', parentId: '4' },
            ]
        },
        {
            id: '5',
            title: '5. 最佳实践',
            parentId: '0',
            children: [
                { id: '5-1', title: '5.1 最佳实践', parentId: '5' },
            ]
        }
    ],
    '2': [
        {
            id: '6',
            title: 'java 基础',
            parentId: '0',
            children: [
                { id: '6-1', title: '1.1 变量', parentId: '6' },
                { id: '6-2', title: '1.2 数据类型', parentId: '6' },
                { id: '6-3', title: '1.3 java 流程控制', parentId: '6' },
                { id: '6-4', title: '1.4 java 数组', parentId: '6' },
            ]
        },
        {
            id: '7',
            title: 'Docker 容器技术快速入门',
            parentId: '0',
            children: [
                { id: '7-1', title: '1.1 基础概念', parentId: '7' },
                { id: '7-2', title: '1.2 常用命令', parentId: '7' },
                { id: '7-3', title: '1.3 最佳实践', parentId: '7' },
            ]
        }
    ]
}

const contentList = {
'1': `321321321`,
'1-1': `
1.安装node.js     node -v
2.安装淘宝镜像     cnpm -v
npm install -g cnpm --registry=https://registry.npm.taobao.org
3.安装webpack     webpack -v
npm install webpack -g
4.安装vue-cli脚手架    vue -V
npm install vue-cli -g
5.创建项目
\`vue create project-one\`

\`\`\`bash
cd project-one
npm run serve
\`\`\`

安装sass
\`\`\`bash
npm install sass-loader@7.3.1 --save-dev
npm install node-sass@4.14.1 --save-dev //这个版本的现在不支持了
\`\`\`
sass-loader@10.0.1
node-sass@6.0.1
`,
'1-2': `
    # 生命周期
\`\`\`
<div id="app"></div>


var vm = new Vue({
    el: '#app',
    template: "<div>{{test}}</div>",
    data: {
        test: 'hello'
    },
    beforeCreate: function (){
        //组件实例刚刚被创建，属性和方法都没有
        console.log('beforeCreate')
    },
    created: function(){
        //属性和方法已经绑定成功，但DOM还没有生成（可以进行一些资源的请求，不能操作dom）
        console.log("created")
    },
    beforeMount: function(){
        //模板编译之前
        console.log(this.$el);
        console.log("beforeMount")
    },
    mounted: function(){
        //模板已经编译完成，类似之前的ready，就是模板引擎数据有没有展示（可以进行一些数据请求操作dom）
        console.log(this.$el);
        console.log("mounted")
    },
    beforeUpdate: function(){
        //组件（数据）更新之前，用于监测实例内的变化
    },
    updated: function(){
        //组件更新完毕
    },
    beforeDestroy: function(){
        //组件销毁之前;  执行this.$destory();后触发
        console.log("beforeDestroy")
    },
    destroyed: function(){
        //组件销毁之后（可以进行一些优化操作，清空定时器，解除绑定事件）
        console.log("destroyed")
    }
})
\`\`\`

问题：
1. 何时需要使用beforeDestory
* 解绑自定义事件event.$off
* 清除定时器
* 解绑自定义的DOM事件，如window scroll等
`,
'1-3': `
# 模板语法
\`\`\`
//数据为普通文本
<div id="app1">
    <div v-text="name"></div>
    <div v-html="name"></div>
    <div>{{name}}</div>
</div>

var vm1 = new Vue({
    el: "#app1",
    data: {
        name: 'wudan'
    }
})
\`\`\`

\`\`\`
//数据为HTML代码
<div id="app1">
    <div v-text="name"></div>
    <div v-html="name"></div>
    <div>{{name}}</div>
</div>

var vm1 = new Vue({
    el: "#app1",
    data: {
        name: '<h1>wudan</h1>'
    }
})
\`\`\`
`,
'1-4': `
# 计算属性
\`\`\`
computed: {
    reversedMessage: function () { // \`this\` 指向 vm 实例  
        return this.message.split('').reverse().join('') 
    }
},
\`\`\`
# 侦听器
> 监听对象的属性
\`\`\`
data () {
    return {
        formInline: {
            bloodPressure: ''
        }
    }
}
\`\`\`
\`\`\`
watch: {
   'formInline.bloodPressure':function(val) {
       //针对监听属性进行的操作
   }
},
\`\`\`
`,
'1-5': `
# 父组件向子组件传值
## 父组件传过来的值在prop中是单向数据流，不能在子组件中改变
\`\`\`
//实现从父组件传count过来，将接收的count传入子组件的data中，点击累加修改count值
<div id="app">
    // :count -- 接收的是数字类型  count --  接收的是字符串类型
	<counter :count='0'></counter>
	<counter :count='1'></counter>
</div>

var counter = {
	template: '<div @click="handleClick">{{number}}</div>',
	props: ['count'],
	data: function(){
		return {
			number: this.count
		}
	},
	methods: {
		handleClick(){
			this.number ++
		}
	}
}
var app = new Vue({
	el: '#app',
	components: {
		counter: counter
	}
})
\`\`\`
# 子组件向父组件传值
\`\`\`
<div id="app">
	<counter @change="handleChange"></counter>

</div>

var counter = {
	template: '<div @click="handleClick">{{number}}</div>',
	data: function(){
		return {
			number: 6
		}
	},
	methods: {
		handleClick(){
			this.number ++;
			this.$emit('change',1);// 传递给父组件1
		}
	}
}
var app = new Vue({
	el: '#app',
	components: {
		counter: counter
	},
	methods: {
		handleChange(value){
			//从子组件获取到的值
			console.log(value)
		}
	}
})
\`\`\`
`,
'1-6': `
# 非父子之间的传值
\`\`\`
<div id="app">
	<child1 content="hello"></child1>
	<child2 content="你好"></child2>
</div>

//创建一个新的实例方法绑定到Vue原型上，让每一个vue实例都能用他的方法
Vue.prototype.bus = new Vue()
Vue.component('child1', {
	props: ['content'],
	template: '<div @click="handleClick">{{content}}</div>',
	methods: {
		handleClick(){
			this.bus.$emit('change', this.content)
		}
	},
	mounted(){
		this.bus.$on('change',function(val) {
			console.log("child1拿到数据："+ val)
		})
	}
})
Vue.component('child2', {
	props: ['content'],
	template: '<div @click="handleClick">{{content}}</div>',
	methods: {
		handleClick(){
			this.bus.$emit('change', this.content)
		}
	},
	mounted(){
		this.bus.$on('change',function(val) {
			console.log("child2拿到数据："+ val)
		})
	}
})

var app = new Vue({
	el: '#app',
	
})
\`\`\`
`,
'1-7': `
# 一般插槽
\`\`\`
<div id="app">
	<child>
		<div>hahhaha</div>
	</child>
</div>

//Vue的插槽(slot)
Vue.component('child', {
	template:  \`<div>
					<p>hello</p>
					<slot></slot>
				</div>\`
})
\`\`\`

# 具名插槽
\`\`\`
<div id="app">
	<child>
        <div slot="header">header</div>
        <div slot="footer">footer</div>
	</child>
</div>

//Vue的插槽(slot)
Vue.component('child', {
	template:  \`<div>
					<p>hello</p>
					<slot name="header"></slot>
					<slot name="footer"></slot>
				</div>\`
})
\`\`\`
# 作用域插槽
### 如果需要父组件的结构样式 使用子组件中的数据
\`\`\`
<div id="app">
	<child>
		<template v-slot="slotProps">
			<li>{{slotProps.item}}</li>
		</template>
	</child>
</div>

Vue.component('child', {
	data: function() {
		return {
			list: [1,2,3,4]
		}
	},
	template:  \`<div>
					<ul>
					<slot v-for="item in list" :item=item></slot>
					</ul>
				</div>\`
})

\`\`\`
`,
'6': `
# java 基础
`,
'6-1': `
# 变量
\`\`\`
// 变量
int a = 10;
String b = "hello";
double c = 3.14;
boolean d = true;
char e = 'a';
\`\`\`
`,
'6-2': `
# 数据类型
\`\`\`
// 数据类型
byte, short, int, long, float, double, boolean, char
String, StringBuilder, StringBuffer
List, Set, Map
ArrayList, LinkedList, HashMap, TreeMap, HashSet, TreeSet
Thread, Runnable, Callable, Future, FutureTask
ThreadLocal, AtomicInteger, AtomicLong, AtomicBoolean, AtomicReference
ThreadLocalRandom, Random, SecureRandom
ThreadLocalRandom, Random, SecureRandom
\`\`\`
`,
'6-3': `
# 流程控制
\`\`\`
// 流程控制
if (condition) {
    // code block
} else {
    // code block
}
\`\`\`
`,
'6-4': `
# 数组
\`\`\`
// 数组
int[] a = new int[10];
int[] a = {1, 2, 3, 4, 5};
\`\`\`
`,
'7': `
# Docker 容器技术快速入门
`,
'7-1': `
    > 本文档介绍 Docker 容器技术的基本使用方法

    # Docker 容器技术快速入门

    Docker 是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中。

    ## 1. 基础概念

    Docker 的核心概念包括：

    - **镜像（Image）**：一个只读的模板，包含了运行应用所需的代码、运行时、库等
    - **容器（Container）**：镜像的运行实例
    - **仓库（Registry）**：存储镜像的地方

    ### 1.1 镜像管理

    镜像是 Docker 的核心组件之一，以下是常用的镜像操作命令。

    ### 1.2 容器生命周期

    容器从创建到销毁的完整生命周期管理。

    ## 2. 常用命令

    以下是 Docker 的常用命令及其说明。

    ### 2.1 镜像操作

    #### 拉取镜像

    从 Docker Hub 拉取镜像到本地：

    \`\`\`bash
    docker pull [镜像名称]

    # 示例：拉取 nginx 镜像
    docker pull nginx

    # 拉取指定版本
    docker pull nginx:1.21
    \`\`\`

    #### 查看本地镜像

    \`\`\`bash
    docker images

    # 输出示例：
    # REPOSITORY   TAG       IMAGE ID       CREATED        SIZE
    # nginx        latest    605c77e624dd   2 weeks ago    141MB
    \`\`\`

    #### 删除镜像

    \`\`\`bash
    docker rmi [镜像名称或ID]

    # 示例
    docker rmi nginx
    docker rmi 605c77e624dd
    \`\`\`

    ### 2.2 容器操作

    #### 运行容器

    \`\`\`bash
    docker run [选项] [镜像名称] [命令]

    # 常用选项：
    # -d：后台运行容器
    # -p：端口映射，格式：主机端口:容器端口
    # --name：指定容器名称
    # -v：挂载数据卷
    # -e：设置环境变量

    # 示例：运行 nginx 容器
    docker run -d -p 80:80 --name nginx_demo nginx
    \`\`\`

    #### 查看容器

    \`\`\`bash
    # 查看运行中的容器
    docker ps

    # 查看所有容器（包括已停止）
    docker ps -a

    # 查看容器详细信息
    docker inspect [容器名称或ID]
    \`\`\`

    #### 启动和停止容器

    \`\`\`bash
    # 启动已停止的容器
    docker start [容器名称或ID]

    # 停止运行中的容器
    docker stop [容器名称或ID]

    # 重启容器
    docker restart [容器名称或ID]
    \`\`\`

    #### 删除容器

    \`\`\`bash
    # 删除已停止的容器
    docker rm [容器名称或ID]

    # 强制删除运行中的容器
    docker rm -f [容器名称或ID]

    # 删除所有已停止的容器
    docker container prune
    \`\`\`

    ### 2.3 日志和调试

    #### 查看容器日志

    \`\`\`bash
    # 查看容器日志
    docker logs [容器名称或ID]

    # 实时查看日志
    docker logs -f [容器名称或ID]

    # 查看最近的日志
    docker logs --tail 100 [容器名称或ID]
    \`\`\`

    #### 进入容器

    \`\`\`bash
    # 以交互式终端进入容器
    docker exec -it [容器名称或ID] /bin/bash

    # 或使用 sh（如果容器中没有 bash）
    docker exec -it [容器名称或ID] /bin/sh
    \`\`\`

    ## 3. Dockerfile

    Dockerfile 是用来构建镜像的文本文件，包含了一系列构建指令。

    ### 3.1 基本语法

    \`\`\`dockerfile
    # 基础镜像
    FROM node:18-alpine

    # 设置工作目录
    WORKDIR /app

    # 复制文件
    COPY package*.json ./

    # 运行命令
    RUN npm install

    # 复制源代码
    COPY . .

    # 暴露端口
    EXPOSE 3000

    # 启动命令
    CMD ["npm", "start"]
    \`\`\`

    ### 3.2 构建镜像

    \`\`\`bash
    # 构建镜像
    docker build -t [镜像名称]:[标签] [Dockerfile所在目录]

    # 示例
    docker build -t myapp:latest .
    \`\`\`

    ## 4. Docker Compose

    Docker Compose 用于定义和运行多容器 Docker 应用。

    ### 4.1 配置文件

    \`\`\`yaml
    version: '3.8'

    services:
    web:
        build: .
        ports:
        - "3000:3000"
        depends_on:
        - db
        environment:
        - NODE_ENV=production

    db:
        image: postgres:14
        volumes:
        - postgres_data:/var/lib/postgresql/data
        environment:
        - POSTGRES_PASSWORD=password

    volumes:
    postgres_data:
    \`\`\`

    ### 4.2 常用命令

    \`\`\`bash
    # 启动服务
    docker-compose up -d

    # 停止服务
    docker-compose down

    # 查看日志
    docker-compose logs -f

    # 重新构建
    docker-compose build
    \`\`\`

    ## 5. 最佳实践

    1. **使用官方镜像**：优先使用官方维护的镜像作为基础镜像
    2. **精简镜像大小**：使用 Alpine 版本，清理不必要的文件
    3. **利用缓存**：合理安排 Dockerfile 指令顺序
    4. **使用 .dockerignore**：排除不需要的文件
    5. **安全考虑**：不在镜像中存储敏感信息

    ---

    *文档持续更新中...*
`
}

export default {
    'get|/api/documentList': () => {
        return {
            code: 200,
            message: 'success',
            data: {
                list: list
            }
        }
    },
    'post|/api/document': (options: any) => {
        const id = `doc-${Date.now()}`
        list.push({
            ...JSON.parse(options.body),
            id
        })
        // 生成唯一ID
        return {
            code: 200,
            message: 'success',
            data: {
                ...JSON.parse(options.body),
                id
            }
        }
    },
    'put|/api/document': (options: any) => {
        console.log(options.body)
        list = list.map(item => item.id == JSON.parse(options.body).id ? { ...item, ...JSON.parse(options.body) } : item)
        return {
            code: 200,
            message: 'success',
            data: list
        }
    },
    'get|/api/document/outline/': (options: any) => {
        const url = options.url;
        const documentId = url.split('/').pop();
        return {
            code: 200,
            message: 'success',
            data: outline[documentId as keyof typeof outline]
        }
    },
    'get|/api/document/content/': (options: any) => {
        const url = options.url;
        const outlineId = url.split('/').pop();
        return {
            code: 200,
            message: 'success',
            data: contentList[outlineId as keyof typeof contentList]
        }
    }
}