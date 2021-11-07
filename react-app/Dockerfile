# based on https://mherman.org/blog/dockerizing-a-react-app/

# build environment
FROM node:16-alpine as build
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH
COPY package.json ./
COPY package-lock.json ./
RUN npm ci
RUN npm install react-scripts@4 -g
COPY public public
COPY src src
RUN npm run build

# production environment
FROM nginx:1
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build /usr/src/app/build /www
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
