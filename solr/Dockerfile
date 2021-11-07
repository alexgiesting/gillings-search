# syntax=docker/dockerfile:1

FROM solr:8

USER root

LABEL io.k8s.description="Run SOLR search in OpenShift" \
	io.k8s.display-name="SOLR 8" \
	io.openshift.expose-services="8983:http" \
	io.openshift.tags="builder,solr,solr8"

# Give the Solr directory to root group (not root user)
# https://docs.openshift.org/latest/creating_images/guidelines.html#openshift-origin-specific-guidelines
RUN chgrp -R 0 /opt/solr \
	&& chmod -R g+rwX /opt/solr \
	&& chown -LR solr:root /opt/solr

RUN chgrp -R 0 /opt/docker-solr \
	&& chmod -R g+rwX /opt/docker-solr \
	&& chown -LR solr:root /opt/docker-solr

# In order to drop the root user, we have to make some directories writable
# to the root group as OpenShift default security model is to run the container
# under random UID.
RUN usermod -a -G 0 solr

USER 8983

# Limit the memory used by Solr's JVM (by default, Solr expects 512 MB of RAM
# for the JVM, but our Openshift container doesn't provide that much memory).
ENV SOLR_HEAP="64m"
ENV SOLR_OPTS="-XX:-UseGCOverheadLimit"

# Create document collections
# https://github.com/docker-solr/docker-solr/blob/master/scripts/solr-precreate
RUN [ "init-var-solr" ]
RUN [ "run-initdb" ]
RUN [ "precreate-core", "citations" ]
COPY --chown=solr:solr ./managed-schema /var/solr/data/citations/conf
# RUN [ "precreate-core", "faculty" ]
