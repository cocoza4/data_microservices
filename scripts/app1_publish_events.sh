host=http://localhost:8080
username=admin
password=pass
topic="test"
curl -X POST "${host}/v1/kafka/topics/${topic}/publish" -u ${username}:${password} -d '{"name": "product_curl1", "description": "desc_curl1"}' -H 'Content-Type: application/json'