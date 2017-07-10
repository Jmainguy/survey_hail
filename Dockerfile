# Maintainer Jonathan Mainguy <jon@soh.re>
FROM centos:centos7
ADD survey_hail /usr/sbin/survey_hail
VOLUME /opt/survey_hail/ /etc/survey_hail/
CMD ["/usr/sbin/survey_hail"]
