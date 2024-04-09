# Redirect to homeserver's homepage

# Endpoints
```
/update
```
updates homeserver's ip

# Data
`ip` stores the ip address of the homeserver

# Installation
### Home server
`update.py` on homeserver. Run as cron job every 15 mins.
### VPS
Use systemd service file to run the flask app on nginx+letsencrypt
