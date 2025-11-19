// 首先引入Mock
import Mock from 'mockjs';
 
// 设置拦截ajax请求的相应时间
Mock.setup({
  timeout: '200-600'
});
 
interface MockConfig {
  [key: string]: any;
}

let configArray: MockConfig[] = [];
// const modules = import.meta.glob('./modules/*/*.ts', { eager: true });
// 同时匹配两种路径模式
const moduleFiles = {
  ...import.meta.glob('./modules/*.ts', {eager: true}),
  ...import.meta.glob('./modules/*/*.ts', {eager: true})
};

Object.keys(moduleFiles).forEach(key => {
  configArray = configArray.concat(moduleFiles[key].default);
})

// export function initMock() {
  console.log('mock初始化完成');
  // 注册所有的mock服务
  configArray.forEach((item) => {
    for (let [path, target] of Object.entries(item)) {
      const protocol = path.split('|');
      Mock.mock(new RegExp('^' + protocol[1]), protocol[0], target);
    }
  });
  console.log('mock初始化完成', configArray);
// }
