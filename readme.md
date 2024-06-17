
# 随机图片 api

- 采用的是 nodejs 来实现的，使用了 express 框架

- 图片来源于互联网

## 使用方法

- 首先,准备好 nodejs 环境, 以及必要的依赖. `npm install` 安装依赖

- 准备好图片资源, 把图片放到 `images` 目录下, 如果要对图片进行分类, 可以在 `images` 目录下新建文件夹, 把图片放到对应的文件夹下. 例如 `images/cat` 文件夹下的图片, 就是猫的图片. 通过 运行 python 脚本 `python3 images2json.py` 来生成图片的索引文件 `images_info.json`, 用于后续的读取,方便获取图片

- 运行 `node index.js` 启动服务

- 访问 `http://localhost:3000/random` 即可随机获取一张图片

## 参数说明

- 访问 `http://localhost:3000/random?sub=xxx&width=xxx&type=xxx&contains=xxx&index=xxx` 即可获取指定的图片

    - `sub`: 文件夹名称, 例如 `http://localhost:3000/random?sub=cat` 即可获取 `images/cat` 文件夹下的图片

    - `width`: 图片宽度, 例如 `http://localhost:3000/random?width=200` 即可获取宽度为 200px 的图片, 默认为 1920px

    - `type`: 图片格式, 例如 `http://localhost:3000/random?type=images` 即可获取原始图片的格式, 如果未指定, 则默认为 webp 格式
    
    - `contains`: 包含关键字, 例如 `http://localhost:3000/random?contains=cat` 即可获取包含 `cat` 关键字的图片, 如果未指定, 则默认为随机获取

    - `index`: 索引, 例如 `http://localhost:3000/random?index=1` 即可获取索引为 1 的图片(按顺序), 如果未指定, 则默认为随机获取

    - `device`: 设备类型, 例如 `http://localhost:3000/random?device=mobile` 即可获取移动端图片, 如果未指定, 则默认为 pc 端图片

    - 上述选项是可以组合使用的, 不过有优先级之分, 注意区分即可(这里不在细分), 不建议组合太多选项, 以免出现问题.

- 最好,使用域名+反向代理的方式来访问, 不然可能会有跨域问题


## docker 部署

- 本项目已经提供了 `Dockerfile` 文件 和 `docker-compose.yml` 文件, 可以直接使用 docker 来部署

- 使用 docker 部署, 可以直接使用 `docker-compose up -d` 来启动服务
