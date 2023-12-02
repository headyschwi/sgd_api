import unittest
import requests
import os

class TestClients(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.base_url = os.getenv("API_URL", "http://localhost:7777/clientes")

    def test1_createClient(self):
        url = f'{self.base_url}'
        data = {"name": "teste", "email": "teste@gmail.com"}

        response = requests.post(url, json=data)
        self.assertEqual(response.status_code, 201, f"Erro ao criar cliente: {response.content}")

    def test2_getClients(self):
        url = f'{self.base_url}'
        response = requests.get(url)

        self.assertEqual(response.status_code, 200, f"Erro ao buscar clientes: {response.content}")

    def test3_getClient(self):
        url = f'{self.base_url}/1'
        response = requests.get(url)

        self.assertEqual(response.status_code, 200, f"Erro ao buscar cliente: {response.content}")

    def test4_updateClient(self):
        url = f'{self.base_url}/1'
        data = {"name": "teste", "email": "testeupdate@gmail.com"}

        response = requests.put(url, json=data)
        self.assertEqual(response.status_code, 200, f"Erro ao atualizar cliente: {response.content}")

    def test5_deleteClient(self):
        url = f'{self.base_url}/1'
        response = requests.delete(url)

        self.assertEqual(response.status_code, 200, f"Erro ao deletar cliente: {response.content}")

if __name__ == '__main__':
    test_suite = unittest.TestSuite()
    test_suite.addTest(TestClients('test1_createClient'))
    test_suite.addTest(TestClients('test2_getClients'))
    test_suite.addTest(TestClients('test3_getClient'))
    test_suite.addTest(TestClients('test4_updateClient'))
    test_suite.addTest(TestClients('test5_deleteClient'))
    runner = unittest.TextTestRunner()
    runner.run(test_suite)
