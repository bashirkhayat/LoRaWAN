import base64
import sys


def main(argv):
	mystr = base64.b64decode(str(argv))
	for c in mystr:
		print ord(c)


if __name__ == "__main__":
	main(sys.argv[1])	