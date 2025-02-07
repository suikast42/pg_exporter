#!/bin/bash

 if [ x"${REGISTRY}" == "x" ]; then
     echo "Build version $BUILD_VERSION push to dockerhub"
     make build
     docker  build  .  -t suikast42/pg_exporter:$BUILD_VERSION
     echo docker push suikast42/pg_exporter:$BUILD_VERSION
     docker push suikast42/pg_exporter:$BUILD_VERSION
 else
     echo "Build version $BUILD_VERSION push to $REGISTRY"
     make build
     docker  build  .  -t  $REGISTRY/suikast42/pg_exporter:$BUILD_VERSION
     echo docker push $REGISTRY/suikast42/pg_exporter:$BUILD_VERSION
     docker push $REGISTRY/suikast42/pg_exporter:$BUILD_VERSION
 fi
