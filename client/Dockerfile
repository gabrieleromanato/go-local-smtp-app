FROM node:lts-alpine

WORKDIR /app

COPY package*.json ./
COPY src ./src
COPY public ./public
COPY index.html ./
COPY jsconfig.json ./
COPY vite.config.js ./
COPY postcss.config.js ./
COPY tailwind.config.js ./
COPY .env ./

RUN npm install http-server -g

RUN npm install
RUN npm run build

EXPOSE 3000

CMD ["http-server", "dist", "-p", "3000"]