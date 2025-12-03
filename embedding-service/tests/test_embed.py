from fastapi.testclient import TestClient
from app.main import app
from app.model import get_model
import numpy as np

client = TestClient(app)

class MockModel:
    def encode(self, text):
        return np.array([0.1, 0.2, 0.3])


def override_get_model():
    return MockModel()


def test_embed():
    app.dependency_overrides[get_model] = override_get_model

    response = client.post("/embed", json={"text": "string"})

    assert response.status_code == 200
    assert response.json() == {"embedding": [0.1, 0.2, 0.3]}

    app.dependency_overrides.clear()
