grafana:
	@echo " > starting up grafana..."
	docker run -d -p 3000:3000 --name=grafana -v grafana-storage:/var/lib/grafana grafana/grafana-enterprise
tarantool:
	@echo " > starting up tarantool..."
	docker run --name mytarantool -d -p 3301:3301 tarantool/tarantool:2.10.0-beta2
prometheus:
	@echo " > starting up prometheus"
	docker run -p 9090:9090 -v $(shell pwd)/prometheus/prometheus.yml prom/prometheus