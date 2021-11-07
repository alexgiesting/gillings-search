# based on https://mherman.org/blog/dockerizing-a-react-app/

# build environment
FROM node:16-alpine as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json ./
COPY package-lock.json ./
RUN npm ci
RUN npm install react-scripts@4 -g
COPY public public
COPY src src
RUN npm run build

# production environment
FROM registry.access.redhat.com/ubi8/nginx-118
COPY ./nginx/nginx.conf "${NGINX_CONF_PATH}"
COPY ./nginx/nginx-default-cfg/*.conf "${NGINX_DEFAULT_CONF_PATH}"
COPY ./nginx/nginx-cfg/*.conf "${NGINX_CONFIGURATION_PATH}"
# COPY --from=build /app/build .
COPY ./nginx/*.html .
CMD ["nginx", "-g", "daemon off;"]
