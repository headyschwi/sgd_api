import unittest
import requests
import os 

class TestProducts(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.base_url = os.getenv("API_URL", "http://localhost:7777/produtos")

    def test1_createProduct(self):
        url = f'{self.base_url}/'
        data = {"name": "Product", "price": 10.0, "stock": 10}

        response = requests.post(url, json=data)
        self.assertEqual(response.status_code, 201, f"Erro ao criar produto: {response.content}")

    def test2_getProducts(self):
        url = f'{self.base_url}'
        response = requests.get(url)

        self.assertEqual(response.status_code, 200, f"Erro ao buscar produtos: {response.content}")

    def test3_getProduct(self):
        url = f'{self.base_url}/1'
        response = requests.get(url)

        self.assertEqual(response.status_code, 200, f"Erro ao buscar produto: {response.content}")

    def test4_updateProduct(self):
        url = f'{self.base_url}/1'
        data = {"name": "Updated Product", "price": 20.0}

        response = requests.put(url, json=data)
        self.assertEqual(response.status_code, 200, f"Erro ao atualizar produto: {response.content}")

    def test5_deleteProduct(self):
        url = f'{self.base_url}/1'
        response = requests.delete(url)

        self.assertEqual(response.status_code, 200, f"Erro ao deletar produto: {response.content}")
    
if __name__ == '__main__':
    test_suite = unittest.TestSuite()
    test_suite.addTest(TestProducts('test1_createProduct'))
    test_suite.addTest(TestProducts('test2_getProducts'))
    test_suite.addTest(TestProducts('test3_getProduct'))
    test_suite.addTest(TestProducts('test4_updateProduct'))
    test_suite.addTest(TestProducts('test5_deleteProduct'))
    runner = unittest.TextTestRunner()
    runner.run(test_suite)

    

    