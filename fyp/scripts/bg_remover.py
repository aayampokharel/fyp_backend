

from rembg import remove, new_session
import sys
import io
from PIL import Image


try:
    session = new_session(
        model_name="u2net",  # Quality model
        providers=['CUDAExecutionProvider']
    )
except Exception as e:
    sys.exit(1)

def remove_bg_balanced(input_data):
    input_image = Image.open(io.BytesIO(input_data))
    
    max_size = 768  # Balanced size
    original_size = input_image.size
    
    if max(original_size) > max_size:
        ratio = max_size / max(original_size)
        new_size = tuple(int(dim * ratio) for dim in original_size)
        input_image = input_image.resize(new_size, Image.Resampling.LANCZOS)
    
    img_byte_arr = io.BytesIO()
    input_image.save(img_byte_arr, format='PNG')
    img_byte_arr = img_byte_arr.getvalue()
    
    output_data = remove(
        img_byte_arr,
        session=session,
        alpha_matting=False,      
        post_process_mask=True,   # for quality 
        only_mask=False
    )
    
    return output_data

try:
    input_data = sys.stdin.buffer.read()
    output_data = remove_bg_balanced(input_data)
    sys.stdout.buffer.write(output_data)
    sys.stdout.buffer.flush()
except Exception as e:
    sys.exit(1)