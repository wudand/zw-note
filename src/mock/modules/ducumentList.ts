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
    }
}