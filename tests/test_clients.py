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
        self.assertEqual(response.status_code, 201)
    
    def teste_getClients(self):
        url = 'http://localhost:'+port+'/clients'
        response = requests.get(url)
        self.assertEqual(response.status_code, 200)
