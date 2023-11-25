import requests

HOST = "http://localhost:8080/v1"
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

    # list all products
    headers = {"Authorization": f"Bearer {token}"}
    product_uri = f"{HOST}/products"
    resp = requests.get(product_uri, headers=headers)
    print("all products", resp.json())

    # create index
    index_uri = f"{HOST}/products/indexes?field=name"
    resp = requests.post(index_uri, headers=headers)
    print("create index", resp.json())

    # create product
    product_uri = f"{HOST}/products"
    data = {
        "name": "new product",
        "description": "new desc"
    }
    resp = requests.post(product_uri, headers=headers, json=data)
    print("create product", resp.json())

    # latest product
    latest_uri = f"{HOST}/products/latest"
    resp = requests.get(latest_uri, headers=headers)
    print("latest product", resp.json())


if __name__ == "__main__":
    main()