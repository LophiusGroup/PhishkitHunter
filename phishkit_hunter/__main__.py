import fire
from phishkit_hunter import Hunter


def main():
    fire.Fire({
        'inspect': Hunter().inspect
    })

if __name__ == "__main__":
    main()
