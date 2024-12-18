FROM node:lts-alpine AS build

WORKDIR /app

COPY package*.json ./
COPY src ./src
COPY public ./public
COPY index.html ./
COPY jsconfig.json ./
COPY vite.config.js ./

RUN npm install
RUN npm run build

FROM nginx:stable-alpine

COPY --from=build /app/dist /usr/share/nginx/html

COPY ./config/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]