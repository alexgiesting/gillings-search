# based on https://mherman.org/blog/dockerizing-a-react-app/
# and https://catalog.redhat.com/software/containers/ubi8/nginx-118/5f521a6f9dd2d5ca7158e5dc

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
# Add application sources to a directory that the assemble script expects them
# and set permissions so that the container runs without root access
USER 0
COPY --from=build /app/build /tmp/src/
RUN chown -R 1001:0 /tmp/src
USER 1001
# Let the assemble script install the dependencies
RUN /usr/libexec/s2i/assemble
# Run script uses standard ways to run the application
CMD /usr/libexec/s2i/run
