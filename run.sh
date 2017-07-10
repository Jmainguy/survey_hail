#!/bin/bash
docker kill survey_hail
docker rm survey_hail
docker run -d --name survey_hail --restart always -v /opt/survey_hail/etc/:/etc/survey_hail/ -v /opt/survey_hail:/opt/survey_hail survey_hail
