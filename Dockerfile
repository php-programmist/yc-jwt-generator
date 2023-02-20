FROM alpine
COPY build/jwt_generator_linux /usr/bin/jwt_generator
RUN chmod +x /usr/bin/jwt_generator
EXPOSE 80
CMD ["jwt_generator"]