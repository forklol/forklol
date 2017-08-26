FROM nginx:alpine
COPY ./dist /www
RUN rm -rf /etc/nginx/conf.d/*
COPY nginx.conf /etc/nginx/conf.d/forklol-web.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]