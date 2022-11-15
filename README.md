# Noman

Nomad management tool, based on DockerHub webhooks.

Job definitions pull the docker tag from consul and restart the job when that value changes. This service received a webhook from DockerHub when a new tag has been built and updates the key which causes the job to restart using the new container version.

Note that it does just enough for my purposes and no more, and much debugging stuff is still present.
