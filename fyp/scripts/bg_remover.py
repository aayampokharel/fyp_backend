from rembg import remove
import sys

try:
    # Read input image bytes
    input_data = sys.stdin.buffer.read()

    output_data = remove(input_data)
    sys.stdout.buffer.write(output_data)
    sys.stdout.buffer.flush()  
except Exception as e:
    print(f"Error: {e}", file=sys.stderr)
    sys.exit(1)


