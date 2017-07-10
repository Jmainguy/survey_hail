# This expects docker is already installed and running
docker:
	@go build
	@docker build -t="survey_hail" .
	@if [ ! -d /opt/survey_hail/etc ]; then mkdir -p /opt/survey_hail/etc; fi
	@cp config.yaml /opt/survey_hail/etc/
	@cp run.sh /opt/survey_hail/
	@echo "cd to /opt/survey_hail/, edit etc/config.yaml, run run.sh"
