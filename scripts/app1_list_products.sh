host=http://localhost:8080
username=admin
password=pass
topic="test"
curl -X GET ${host}/v1/products/ -u ${username}:${password}