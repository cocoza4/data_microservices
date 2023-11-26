host=http://localhost:8080
username=admin
password=pass
curl -X POST ${host}/v1/products/ -u ${username}:${password} -d '{"name": "product_curl1", "description": "desc_curl1"}' -H 'Content-Type: application/json'
curl -X POST ${host}/v1/products/ -u ${username}:${password} -d '{"name": "product_curl2", "description": "desc_curl2"}' -H 'Content-Type: application/json'
curl -X POST ${host}/v1/products/ -u ${username}:${password} -d '{"name": "product_curl3", "description": "desc_curl3"}' -H 'Content-Type: application/json'
curl -X POST ${host}/v1/products/ -u ${username}:${password} -d '{"name": "product_curl4", "description": "desc_curl4"}' -H 'Content-Type: application/json'