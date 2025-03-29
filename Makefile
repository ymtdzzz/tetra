.PHONY: setup-db download-test-data clean-test-data

setup-db: download-test-data
	@echo "Importing test data..."
	@cd docker/init/mysql && \
	tar xzf world-db.tar.gz && \
	docker compose -f ../../../docker-compose.yml exec -T mysql mysql -uroot -proot < world-db/world.sql
	@make clean-test-data

download-test-data:
	@mkdir -p docker/init/mysql
	@echo "Downloading test data..."
	@cd docker/init/mysql && \
	curl -O https://downloads.mysql.com/docs/world-db.tar.gz

clean-test-data:
	@echo "Cleaning up test data..."
	@rm -rf docker/init/mysql
