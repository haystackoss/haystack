import click

control_content = """Package: nabaz-test
Version: {version}
Architecture: amd64
Maintainer: nabaz.io <hello@nabaz.io>
Homepage: https://nabaz.io/
Depends: git
Description: The nabaz test runner
 analyzes code changes and selects which tests to run.
"""

def write_to_file(version, path):
    with open(path, 'w') as f:
        f.write(control_content.format(version=version))

@click.option('--version', help='The version of the package', required=True)
@click.option('--output', help='The output path', required=True)
@click.command()
def main(version, output):
    write_to_file(version, output)

if __name__ == "__main__":
    main()