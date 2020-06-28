from unittest import mock

import unittest
import mimik

class MimikTest(unittest.TestCase):
    
    @mock.patch('builtins.open', new_callable=mock.mock_open(read_data='version=v1'))
    def test_mimik_version(self, mock_open):
        self.assertEqual(mimik.get_version(), 'v1')
        
    @mock.patch('requests.get')
    @mock.patch('mimik.get_version')
    def test_mimik_passthrough(self, mock_get_version, mock_get):
        mock_get.return_value.text = 'users (v1)'
        mock_get_version.return_value = 'v1'

        self.assertEqual(mimik.mimik('groups', {}, 'passthrough', 'http://users:5000/users', 'false'), 'groups (v1) -> users (v1)')

    @mock.patch('mimik.get_version')
    def test_mimik_edge(self, mock_get_version):
        mock_get_version.return_value = 'v1'

        self.assertEqual(mimik.mimik('roles', {}, 'edge', 'http://roles:5000/roles', 'false'), 'roles (v1)')

    def test_mimik_error(self):
        self.assertEqual(mimik.mimik('roles', {}, 'edge', 'http://roles:5000/roles', 'true'), ('Error', 503))


if __name__ == '__main__':
    unittest.main()