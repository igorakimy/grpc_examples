# Генерация файлов протокола буферов
gen-proto:
	# путь к proto-файлу, который нужно скомпилировать
	protoc product_info.proto \
	  # путь, по которому будет сохранен сгенерированный код
	  --go_out=. \
	  --go-grpc_out=require_unimplemented_servers=false:. \
	  --go_opt=paths=source_relative \
	  --go-grpc_opt=paths=source_relative \
	  # путь, по которому хранится исходный proto-файл и его зависимости
	  --proto_path=.

	# protoc product_info.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=.
	# protoc -I . --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt paths=source_relative product_info.proto
	# protoc -I proto --openapiv2_out pb --openapiv2_opt logtostderr=true proto/product_info.proto


# Генерация сертификатов безопасности
gen-root-ca-key:
	openssl genrsa -out ca.key 2048
gen-root-certificate:
	openssl req -new -x509 -days 365 -key ca.key -subj "/C=RU/ST=DPR/L=Donetsk/O=Akim, Inc./CN=Akim Root CA" -out ca.crt
gen-server-key:
	openssl req -newkey rsa:2048 -nodes -keyout server.key -config localhost.cnf -out server.csr
gen-server-certificate:
	openssl x509 -req -extfile localhost.cnf -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extensions v3_req

# Генерация mock-файлов
gen-mock:
	mockgen -source=.\server\ecommerce\product_info_grpc.pb.go -destination=.\mock_prodinfo\prodinfo_mock.go ProductInfoClient

# Билдинг докер образа и деплой его на Docker Hub
docker-build:
	docker image build -t igorakimov/grpc-productinfo-server:1.0.0 -f server/Dockerfile .
docker-push:
	docker image push igorakimov/grpc-productinfo-server:1.0.0

# Применение ресурсов kubernetes на основе yaml-файлов
k8s-apply-resources:
	kubectl apply -f server/kubernetes-server.yml &&
	kubectl apply -f client/kubernetes-client.yml

# Получение развертывания во всех пространствах имен в kubernetes
k8s-get-deployments:
	kubectl get deployments --all-namespaces

# Получение подов указанного развертывания в kubernetes
k8s-stop-and-delete-pods:
	kubectl delete -n default deployment grpc-productinfo-server

# Удаление джоба в kubernetes
k8s-stop-and-delete-job:
	kubectl delete job grpc-productinfo-client

	# [Environment]::SetEnvironmentVariable("KUBECONFIG", $HOME + "\.kube\config", [EnvironmentVariableTarget]::Machine)