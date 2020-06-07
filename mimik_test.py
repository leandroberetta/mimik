from unittest import mock

import unittest
import mimik

class MimikTest(unittest.TestCase):

    def test_mimik_version(self):
        version = mimik.get_version()

        self.assertEqual(version, 'v1')
        
    @mock.patch('requests.get')
    def test_mimik_passthrough(self, mock_get):
        mock_get.return_value.content = 'groups (v2) -> users (v1)'

        self.assertEqual(mimik.mimik('users', {'mimik-history': 'groups (v2)'}, 'passthrough', 'http://users:5000/users'), 'groups (v2) -> users (v1)')

    def test_mimik_edge(self):
        self.assertEqual(mimik.mimik('roles', {'mimik-history': 'groups (v2) -> users (v1)'}, 'edge', 'http://roles:5000/roles'), 'groups (v2) -> users (v1) -> roles (v1)')


if __name__ == '__main__':
    unittest.main()