FROM node AS build
WORKDIR /app

ENV NODE_ENV=production

COPY package.json package-lock.json ./

RUN npm ci

COPY . .

RUN npm run build

FROM nginxinc/nginx-unprivileged:1.16-alpine-perl as base
WORKDIR /usr/share/nginx/html

COPY --from=build /app/build .
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY scripts/entrypoint.sh entrypoint.sh

EXPOSE 3000

USER nginx
ENTRYPOINT [ "sh", "entrypoint.sh" ]