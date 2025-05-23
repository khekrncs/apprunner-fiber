from locust import HttpUser, task, between
import random
import json

class MyApiUser(HttpUser):
    wait_time = between(1, 5)  

    def on_start(self):
        self.user_id = None 

    @task
    def create_and_get_user(self):
        user_data = {
            "name": f"User{random.randint(1000, 9999)}",
            "email": f"user{random.randint(1000,9999)}@example.com"
        }
        headers = {"Content-Type": "application/json", "x-api-key": "ecs-secret-key"}
        with self.client.post("/dev/api/v1/users", data=json.dumps(user_data), headers=headers, catch_response=True) as response:
            if response.status_code == 200:
                try:
                    self.user_id = response.json()["id"]
                    response.success()
                except Exception as e:
                    response.failure(f"JSON parse error or missing id: {e}")
                    return
            else:
                response.failure(f"Failed to create user: {response.status_code}")
                return
        if self.user_id:
            self.client.get(f"/dev/api/v1/users/{self.user_id}", headers=headers)

    @task(1)
    def healthCheck(self):
        self.client.get("/dev/api/v1/healthCheck")
