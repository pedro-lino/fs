FROM nginx:alpine

RUN apk add --no-cache openssl
RUN mkdir -p /etc/nginx/ssl

# Generate self-signed certificate for development
RUN openssl req -x509 -newkey rsa:4096 -nodes -keyout /etc/nginx/ssl/key.pem -out /etc/nginx/ssl/cert.pem -days 365 -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost"

COPY nginx/nginx.conf /etc/nginx/nginx.conf

# Expose only the HTTPS port
EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]