FROM nginx:alpine
# Copia el contenido de tu carpeta frontend al directorio de Nginx
COPY . /usr/share/nginx/html
# Copia tu configuración específica de nginx si existe
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
