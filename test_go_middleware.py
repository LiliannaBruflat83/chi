import os
import shutil
import subprocess
import unittest


class GoMiddlewareChainTest(unittest.TestCase):
    def test_go_middleware_short_circuits(self):
        go = shutil.which("go")
        if go is None:
            self.skipTest("Go toolchain not installed")

        repo_root = os.getcwd()
        result = subprocess.run(
            [go, "test", "./..."],
            cwd=repo_root,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            check=False,
        )
        output = result.stdout
        self.assertEqual(result.returncode, 0, msg=output)
        self.assertNotIn("superfluous response.WriteHeader", output, msg=output)


if __name__ == "__main__":
    unittest.main()
