# Redirect to homeserver's homepage

# Endpoints
```
/update
```
updates homeserver's ip

# Data
`secret` stores the 'secret' to keep `/update` secure
`ip` stores the ip address of the homeserver

# Installation
### Home server
`update.sh` on homeserver. Run as cron job every 15 mins.
### VPS
Use systemd service file to run the flask app on nginx+letsencrypt

# Testing Data
for testing flask app with `data.json`:
```
curl -H "Content-Type: application/json" --data @data.json http://127.0.0.1:5000/update
```
