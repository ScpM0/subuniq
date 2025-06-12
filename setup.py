from setuptools import setup, find_packages

setup(
    name='subuniq',
    version='1.0.0',
    author='ScpM0',
    description='Remove duplicate subdomains CLI tool',
    packages=find_packages(),
    entry_points={
        'console_scripts': [
            'subuniq=subuniq.cli:main',
        ],
    },
    python_requires='>=3.6',
)
