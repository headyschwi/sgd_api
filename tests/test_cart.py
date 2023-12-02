import unittest
import requests
import os

class TestCart(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.base_url = os.getenv("API_URL", "http://localhost:7777/carrinho")
        
        user_response = requests.post("http://localhost:7777/clientes", json={"name": "Cart Test", "email": "cart@test"})
        
        if user_response.status_code == 201:
            cls.user_id = user_response.json().get("data")["ID"]
        else:
            raise ValueError("Falha na criação do usuário: " + user_response.text)

        product_response = requests.post("http://localhost:7777/produtos", json={"name": "Cart Test", "price": 10.0, "stock": 10})
        
        if product_response.status_code == 201:
            cls.product_id = product_response.json().get("data")["ID"]
        else:
            raise ValueError("Falha na criação do produto: " + product_response.text)

    @classmethod
    def tearDownClass(cls):
        requests.delete(f"http://localhost:7777/clientes/{cls.user_id}")
        requests.delete(f"http://localhost:7777/produtos/{cls.product_id}")
        
    def test1_addToCart(self):
        url = f"{self.base_url}/adicionar"
        data = {"client_id":self.user_id, "product_id": self.product_id, "amount": 1}

        response = requests.post(url, json=data)
        self.assertEqual(response.status_code, 200, f"Erro ao adicionar carrinho: {response.content}")

    def test2_updateCartItem(self):
        url = f'{self.base_url}/atualizar'
        data = {"client_id":self.user_id, "product_id": self.product_id, "amount": 2}

        response = requests.put(url, json=data)
        self.assertEqual(response.status_code, 200, f"Erro ao atualizar item de carrinho: {response.content}")

    def test3_removeFromCart(self):
        url = f'{self.base_url}/remover'
        data = {"client_id":self.user_id, "product_id": self.product_id}
        
        response = requests.delete(url, json=data)

        self.assertEqual(response.status_code, 200, f"Erro ao remover item de carrinho: {response.content}")


if __name__ == '__main__':
    test_suite = unittest.TestSuite()
    test_suite.addTest(TestCart('test1_addToCart'))
    test_suite.addTest(TestCart('test2_updateCartItem'))
    test_suite.addTest(TestCart('test3_removeFromCart'))
    runner = unittest.TextTestRunner()
    runner.run(test_suite)

