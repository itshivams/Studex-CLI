from setuptools import setup, find_packages

setup(
    name="studex-cli",
    version="1.0.1",
    description="A command line interface for Studex platform",
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
