import os
from setuptools import setup, find_packages

here = os.path.abspath(os.path.dirname(__file__))
readme_path = os.path.join(here, "README.md")
long_description = open(readme_path, encoding="utf-8").read() if os.path.exists(readme_path) else ""

setup(
    name="studex-cli",
    version="1.0.1",
    description="A command line interface for Studex platform",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="itshivams",
    license="MIT",
    url="https://github.com/itshivams/Studex-CLI",
    packages=find_packages(),
    entry_points={
        "console_scripts": [
            "studex-cli=studex_cli.run:main",
        ]
    },
)
