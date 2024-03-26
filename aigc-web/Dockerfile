FROM node:18.16.1-alpine3.18 AS node

COPY . /app/web
WORKDIR /app/web
# install packages
RUN npm install pnpm -g --registry https://registry.npmmirror.com/

ARG NODE_OPTIONS=--max_old_space_size=4096
RUN pnpm config set registry https://registry.npmmirror.com/
RUN pnpm install

# build
RUN pnpm build
