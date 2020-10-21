# Basics 

To use the config as is , add the following lines to your hosts file, the run docker-compose up -d

127.0.0.1	scon.local
127.0.0.1	front.scon.local
127.0.0.1	monitor.scon.local
127.0.0.1	prom.scon.local

# Exposed services

The following services will be exposed on your local machine once all the services are running: 

- scon.local : This is the backend api service
- front.scon.local : This is the frontend (gui) service
- monitor.scon.local : This is the grafana monitoring service
- prom.scon.local : This is the prometheus service and does not need to be exposed but is available for debugging

Traefik is currently in place to handle traffic to the different services, this is not a requirement if you wish to hit the services directly.

# Alternative setup

Alternatively edit the .env file and set up the relavent domains to point to the front end and backend (the monitoring is optional).

Ensure that API_URL is set to the base address of your back end domain (e.g http://scon.local)