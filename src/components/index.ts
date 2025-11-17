// 导入声明
// 导入声明
import type { App } from "vue";

// 第三方组件
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';

import {Pane, Splitpanes} from 'splitpanes';
import 'splitpanes/dist/splitpanes.css';

export default {
    install(app: App) {
        // 导入布局插件
        app.component('Splitpanes', Splitpanes);
        app.component('Pane', Pane);
        app.use(ElementPlus) // ELEMENT 组件
    }
}
