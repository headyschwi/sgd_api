import unittest
import requests
import os

port = os.getenv('API_PORT', '333')

class TestClients(unittest.TestCase):
    def teste_createClient(self):
        url = 'http://localhost:'+port+'/clients'
        data = {
            "name": "teste",
            "email": "teste@gmail.com"
        }

        response = requests.post(url, json=data)
        if response.status_code != 201:
            raise Exception('Erro ao criar cliente')

    def teste_getClients(self):
        url = 'http://localhost:'+port+'/clients'
        response = requests.get(url)
        
        if response.status_code != 200:
            raise Exception('Erro ao buscar clientes')
