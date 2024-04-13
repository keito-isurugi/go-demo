# SQLをコンテナに流す
exec-schema:
	cat ./DDL/*.up.sql > ./DDL/schema.sql
	docker cp DDL/schema.sql go-demo-db:/ && docker exec -it go-demo-db psql -U postgres -d go_demo -f /schema.sql
	rm ./DDL/schema.sql
exec-dummy:
	docker cp DDL/insert_dummy_data.sql go-demo-db:/ && docker exec -it go-demo-db psql -U postgres -d go_demo -f /insert_dummy_data.sql

# テーブルをリフレッシュ
refresh-schema:
	@make exec-schema
	@make exec-dummy