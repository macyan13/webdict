Webdict is a self-hosted web application to create and use your own dictionaries for different languages.

### Main features
* Tags and search by tags support.
* Multi-language support. You can create dictionaries for different languages.
* Multi-account support. As admin, you can create many users with their own dictionaries.
* Login via email.
* Automatic backup.
* Letsencrypt support with automatic renew.


### Tech details
* Current implementation relies on MongoDB as a database. It's possible to easily change DB providing different implementation for app repository interfaces.
* DB queries cache layer is application RAM.
* Docker compose installation supports automatic renew for letsencrypt cert by initial cert has to be acquired manually. It's possible to do it with the following command.
```
docker compose run --rm  certbot certonly --webroot --webroot-path /var/www/certbot/ -d example.org
```

**Note**: Need to comment entrypoint override for this container in docker compose file before the command run.

As Letsencrypt makes acknowledge via http for initial run you need to adjust nginx config to serve probe without SSl certs.

```
server {
listen 80;
listen [::]:80;

    server_name ${SERVER_NAME} www.${SERVER_NAME};
    server_tokens off;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        proxy_pass http://webdict:${APP_PORT};
    }
}
```


## Related projects
* [mongodb-backup](https://github.com/macyan13/mongodb-backup)
