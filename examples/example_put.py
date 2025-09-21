import requests

url = "https://darkroom.ast.cam.ac.uk:9443/mybucket/report.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&..."
with open("local_report.pdf", "rb") as f:
    response = requests.put(url, data=f)

print(response.status_code)  # 200 means success
