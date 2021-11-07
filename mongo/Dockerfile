# syntax=docker/dockerfile:1

FROM mongo:5

RUN mkdir /var/mongodb
RUN chgrp -R 0 /var/mongodb \
	&& chmod -R g+rwX /var/mongodb \
	&& chown -LR mongodb:mongodb /var/mongodb

COPY mongod.conf /etc/mongod.conf

CMD [ "mongod", "--config", "/etc/mongod.conf" ]
