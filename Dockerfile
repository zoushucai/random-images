# 使用 Node.js 的 LTS 版本作为基础镜像
FROM node:18-alpine

# 设置工作目录
WORKDIR /usr/src/app

# 复制 package.json 和 package-lock.json 到工作目录
COPY package*.json ./

# 安装依赖
RUN npm install --production
RUN npm install sharp@latest
# 不再复制 images 文件夹
COPY ./index.js ./index.js
# 暴露容器端口
EXPOSE 2113

# 启动应用程序
CMD ["node", "index.js"]

# docker build -t my-rdimg-app .
