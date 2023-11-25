import requests

HOST = "http://localhost:8080/v1"
TOPIC = "newtopic"
USERNAME = "admin"
PASSWORD = "pass"

def main():

    # login
    login_uri = f"{HOST}/login"
    credentials = {
        "username": USERNAME,
        "password": PASSWORD
    }
    resp = requests.post(login_uri, json=credentials)
    token = resp.json()["token"]
    print("Token", token)

    # create topic
    headers = {"Authorization": f"Bearer {token}"}
    topic_uri = f"{HOST}/kafka/topics/{TOPIC}"
    resp = requests.post(topic_uri, headers=headers)
    print("create topic", resp.json())

    # publish event
    publish_uri = f"{HOST}/kafka/topics/{TOPIC}/publish"
    for i in range(10):
        data = {
            "name": "test" + str(i),
            "description": "test" + str(i)
        }
        resp = requests.post(publish_uri, json=data, headers=headers)
        print(resp.json())

if __name__ == "__main__":
    main()