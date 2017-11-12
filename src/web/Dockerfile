FROM node:8-alpine

ARG version=v0.0.0

USER root
RUN apk --update upgrade && \
  apk add nginx sed && \
  rm -rf /var/cache/apk/* && \
  mkdir /build

COPY . /build

RUN rm -rf /build/dist

#RUN npm install -g node-sass
#RUN npm rebuild node-sass --force 

WORKDIR /build
RUN npm install
RUN sed -i -E "s/VERSION: '\"(v(.*))?\"'/VERSION: '\"$version\"'/g" ./config/dev.env.js
RUN sed -i -E "s/VERSION: '\"(v(.*))?\"'/VERSION: '\"$version\"'/g" ./config/prod.env.js
RUN VERSION=$version npm run build

RUN mkdir /www && \
  cp -r /build/dist/* /www && \
  chown -R nginx:nginx /www && \
  chown -R nginx:nginx /var/lib/nginx && \
  rm -rf /build

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]
#CMD ["nginx", "-g", "daemon off;"]
